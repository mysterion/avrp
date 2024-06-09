package dist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func latestRelease() (Release, error) {
	repoOwner := "mysterion"
	repoName := "aframe-vr-player"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Release{}, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Release{}, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return Release{}, err
	}

	return release, nil
}

func allReleases() ([]Release, error) {
	repoOwner := "mysterion"
	repoName := "aframe-vr-player"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", repoOwner, repoName)

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
