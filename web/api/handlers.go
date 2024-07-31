package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mysterion/avrp/internal/thumbnails"
)

var servDir string

func Init(sd string) {
	servDir = sd
}

type ListData struct {
	Files   []string `json:"files"`
	Folders []string `json:"folders"`
}

const distPath = "/"

const listPath = "/list/"

func listHandler(w http.ResponseWriter, r *http.Request) {
	listDir, err := url.PathUnescape(filepath.Join(servDir, filepath.FromSlash(r.URL.Path[len(listPath):])))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	files, folders, err := listFilesAndFolders(listDir)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	listData := ListData{
		Files:   files,
		Folders: folders,
	}
	jsonData, err := json.Marshal(listData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func listFilesAndFolders(dirPath string) ([]string, []string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}
	files := make([]string, 0)
	folders := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		} else {
			files = append(files, entry.Name())
		}
	}
	return files, folders, nil
}

const filePath = "/file/"

const thumbPath = "/thumb/"

func thumbHandler(w http.ResponseWriter, r *http.Request) {

	if !thumbnails.Available {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Thumbnails not available for now. Check logs for more info.")
		return
	}

	id := r.URL.Query().Get("id")
	file, err := url.PathUnescape(filepath.Join(servDir, filepath.FromSlash(r.URL.Path[len(thumbPath):])))

	if id == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = os.Stat(file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not Found")
		return
	}

	if !thumbnails.Generated(file) {
		thumbnails.Generate(file)
	}

	p, err := thumbnails.Get(id, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	http.ServeFile(w, r, p)
}
