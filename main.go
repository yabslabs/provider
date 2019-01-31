package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yabslabs/yabs-exporter/storage"
)

const (
	url               = "https://api.github.com/orgs/yabslabs"
	gitAccessTokenKey = "GIT_ACCESS_TOKEN"
	gitUsernameKey    = "GIT_USERNAME"
	githubAcceptType  = "application/vnd.github.wyandotte-preview+json"
	acceptKey         = "Accept"
)

var (
	accessToken = ""
	username    = ""
)

func init() {
	if accessToken = os.Getenv(gitAccessTokenKey); accessToken == "" {
		log.Fatalf("access-token not provided ($%v)", gitAccessTokenKey)
	}
	if username = os.Getenv(gitUsernameKey); username == "" {
		log.Fatalf("username not provided ($%v)", gitUsernameKey)
	}
}

//main exports yabs
func main() {
	client := &http.Client{}
	storage := storage.NewStorage()

	repos, err := getRepos(client)
	if err != nil {
		log.Fatalf("repo-check failed: %v", err)
	}
	migrationID, err := startBackup(client, repos)
	if err != nil {
		log.Fatalf("create backup failed: %v", err)
		os.Exit(1)
	}

	if err = awaitBackup(client, migrationID); err != nil {
		log.Fatalf("await backup failed: %v", err)
		os.Exit(1)
	}

	if err = downloadExport(client, migrationID, storage); err != nil {
		log.Fatalf("download backup failed: %v", err)
		os.Exit(1)
	}
	log.Printf("everything is fine len: %v", len(repos))
}

func awaitBackup(client *http.Client, migrationID int) error {
	req, err := createGETRequest(fmt.Sprintf("%v/migrations/%v", url, migrationID), nil)
	if err != nil {
		return err
	}
	for {
		time.Sleep(1 * time.Second)
		mig := Migrations{}
		if err = sendRequest(client, req, &mig); err != nil {
			return err
		}
		fmt.Println(mig.State)
		if strings.ToLower(mig.State) == "failed" {
			return fmt.Errorf("backup failed")
		}
		if strings.ToLower(mig.State) == "exported" {
			return nil
		}
	}
}

func getRepos(client *http.Client) (Repos, error) {
	req, err := createGETRequest(fmt.Sprintf("%v/repos", url), nil)
	if err != nil {
		log.Fatalf("unable to create get-request: %v", err)
	}
	repos := make(Repos, 0)
	if err = sendRequest(client, req, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func startBackup(client *http.Client, repos Repos) (int, error) {
	repoList := make([]string, 0, len(repos))
	for _, repo := range repos {
		repoList = append(repoList, repo.Name)
	}
	repositories := &Repositories{Repositories: repoList}
	req, err := createPOSTRequest(fmt.Sprintf("%v/migrations", url), repositories)
	if err != nil {
		return 0, err
	}
	migrations := &Migrations{}
	err = sendRequest(client, req, migrations)
	return migrations.ID, err
}

func downloadExport(client *http.Client, migrationID int, storage storage.Storage) error {
	req, err := createGETRequest(fmt.Sprintf("%v/migrations/%v/archive", url, migrationID), nil)
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	defer func() {
		if err = response.Body.Close(); err != nil {
			log.Println("YABS--jVf8: unable to close body: ", err)
		}
	}()
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return storage.Save(os.TempDir(), "github.bak", body)
}
