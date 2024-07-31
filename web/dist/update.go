package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mysterion/avrp/internal/utils"
)

type Asset struct {
	ID                 int       `json:"id"`
	Size               int       `json:"size"`
	CreatedAt          time.Time `json:"created_at"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
}

type Release struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Tag    string  `json:"tag_name"`
	URL    string  `json:"html_url"`
	Assets []Asset `json:"assets"`
}

func (r *Release) Version() int {
	vs := strings.ReplaceAll(r.Tag, ".", "")
	v, err := strconv.Atoi(vs)
	if err != nil {
		return 0
	}
	return v
}

// returns `true` if update-check was more than 7 days ago
// or running this app for the first time
func checkUFile() (bool, error) {

	lastUpdate, err := os.ReadFile(utils.UpdateFile)
	if errors.Is(err, fs.ErrNotExist) {
		return true, nil
	} else if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, string(lastUpdate))
	if err != nil {
		return false, err
	}

	if t.Before(time.Now().AddDate(0, 0, -7)) {
		return true, nil
	}
	return false, nil
}

func UpdateUFile() error {
	fd, err := os.OpenFile(utils.UpdateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	t := time.Now().Format(time.RFC3339)

	_, err = fd.WriteString(t)
	if err != nil {
		return err
	}

	return nil
}

func LatestRelease() (Release, error) {
	repoOwner := "mysterion"
	repoName := "aframe-vr-player"
	var r Release

	url := os.Getenv("URL_LATEST_RELEASE")
	if url == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return r, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func AllReleases() ([]Release, error) {
	repoOwner := "mysterion"
	repoName := "aframe-vr-player"

	url := os.Getenv("URL_ALL_RELEASES")
	if url == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", repoOwner, repoName)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]Release, 0), err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return make([]Release, 0), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return make([]Release, 0), err
	}

	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return make([]Release, 0), err
	}

	return releases, nil
}

// Fetches the latest dist if no dist is found
// or
// Updates the existing
func TryUpdate() {
	tryUpdate, err := checkUFile()

	v := Ver()
	log.Printf("Current version: %v\n", v)

	valid := Valid()
	if !valid {
		tryUpdate = true
	} else if valid && err != nil {
		log.Println("ERR: while checking for update", err)
		log.Println("Skipping Update check")
		return
	}

	if !tryUpdate {
		return
	}

	if v == math.MaxInt {
		return
	}

	log.Println("Checking for updates")

	r, err := LatestRelease()
	if err != nil {
		log.Println("ERR: Failed to fetch latest Release", err)
		log.Println("Skipping Update check")
		return
	}

	log.Printf("Latest release: %v, v:%v\n", r.Tag, r.Version())

	if v >= r.Version() {
		log.Println("Already on the latest version")
		if v == r.Version() {
			utils.Panic(UpdateUFile())
		}
		return
	}

	log.Println("Downloading the latest version")

	err = DownloadRelease(r)
	if err != nil {
		log.Println("ERR: Failed to fetch latest Release", err)
		log.Println("Skipping Update check")
		if !Valid() {
			if err := Delete(); errors.Is(err, fs.ErrNotExist) {
				panic(err)
			}
		}
		return
	}

	utils.Panic(UpdateUFile())
}
