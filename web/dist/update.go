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
)

type Commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Tag struct {
	Name       string `json:"name"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
	NodeId     string `json:"node_id"`
}

func (t *Tag) Version() int {
	vs := strings.ReplaceAll(t.Name, ".", "")
	v, err := strconv.Atoi(vs)
	if err != nil {
		return 0
	}
	return v
}

func LatestTag() (Tag, error) {
	tags, err := AllTags()
	if err != nil {
		return Tag{}, err
	}
	return tags[0], nil
}

func AllTags() ([]Tag, error) {

	url := os.Getenv("URL_ALL_RELEASES")
	if url == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", RepoOwner, RepoName)
	}

	var tags []Tag

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return tags, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return tags, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tags, err
	}

	err = json.Unmarshal(body, &tags)
	if err != nil {
		return tags, err
	}

	return tags, err
}

// Updates if latest version > current version
func Update() {

	v := Ver()
	log.Printf("Current version: %v\n", v)

	log.Println("Checking for updates")

	t, err := LatestTag()
	if err != nil {
		log.Println("ERR: Failed to fetch latest Release", err)
		log.Println("Skipping Update check")
		return
	}

	log.Printf("Latest release: %v, v:%v\n", t.Name, t.Version())

	if v >= t.Version() {
		log.Println("Already on the latest version")
		return
	}

	log.Println("Downloading the latest version")

	err = DownloadTag(t)
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
