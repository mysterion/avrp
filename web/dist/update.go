package dist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"
)

const RepoOwner = "mysterion"
const RepoName = "aframe-vr-player"

// curl https://api.github.com/repos/mysterion/aframe-vr-player/commits?per_page=1
type Commit struct {
	Sha     string        `json:"sha"`
	Details CommitDetails `json:"commit"`
}

type CommitDetails struct {
	Sha    string `json:"sha"`
	Author struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	} `json:"author"`
	Message string `json:"message"`
}

func GetLatestCommit() (Commit, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=1", RepoOwner, RepoName)

	var commit Commit
	var commits []Commit

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return commit, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return commit, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return commit, err
	}
	err = json.Unmarshal(body, &commits)
	if err != nil {
		return commit, err
	}
	if len(commits) < 1 {
		return commit, fmt.Errorf("No commits found 🤨")
	}
	return commits[0], err
}

func GetCommit(sha string) (CommitDetails, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/%s", RepoOwner, RepoName, sha)

	var commit CommitDetails

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return commit, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return commit, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return commit, err
	}
	err = json.Unmarshal(body, &commit)
	if err != nil {
		return commit, err
	}
	return commit, nil
}

func Update(sha string) error {

	if sha == "latest" {
		v := Ver()
		log.Printf("Current version: %v\n", v)
		log.Println("Checking for updates")
		c, err := GetLatestCommit()
		if err != nil {
			log.Println("ERR: Failed to fetch latest version", err)
			return err
		}

		log.Printf("Commit Details :-\n\nDate: %v\nMessage:\n%v\n\n", c.Details.Author.Date, c.Details.Message)

		if Valid() {
			if v == c.Sha {
				o := "this"
				if sha == "latest" {
					o = "the latest"
				}
				log.Printf("Already on %s version\n", o)
				return nil
			}
		}

		sha = c.Sha

	} else {
		c, err := GetCommit(sha)
		if err != nil {
			log.Printf("ERR: Failed to get %s of aframe-vr-player\n", sha)
			log.Printf("ERR: %v\n", err.Error())
			return err
		}
		log.Println(c, err)
		log.Printf("Commit Details :-\n\nDate: %v\nMessage:\n%v\n\n", c.Author.Date, c.Message)
	}

	fmt.Println("\n\nPress Enter to update")
	fmt.Scanln()

	err := DownloadCommit(sha)
	if err != nil {
		log.Println("ERR: Failed to fetch latest Release", err)
		if !Valid() {
			if err := Delete(); errors.Is(err, fs.ErrNotExist) {
				panic(err)
			}
		}
		return err
	}
	return nil
}
