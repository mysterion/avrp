package dist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

// TODO:
// when http query to github is done. LAST_UPDATE_CHECK file should be updated
// this function will ONLY READ the LAST_UPDATE_CHECK file and return true or false

func CheckUFile() (bool, error) {

	fd, err := os.Open(utils.UpdateFile)
	if err != nil {
		return false, err
	}
	defer fd.Close()

	lastUpdate := make([]byte, 128)
	n, err := fd.Read(lastUpdate)
	if n == 0 {
		return true, nil
	} else if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, string(lastUpdate[:n]))
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
