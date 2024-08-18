package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

// Updates if latest version > current version
func Update() {

	v := Ver()
	log.Printf("Current version: %v\n", v)

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
}
