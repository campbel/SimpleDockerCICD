package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

const (
	buildContainer = "campbel/pipeline-build:latest"
)

var (
	dockerUsername = os.Getenv("DOCKER_USERNAME")
	dockerPassword = os.Getenv("DOCKER_PASSWORD")
	githubSecret   = os.Getenv("GITHUB_SECRET")
)

func main() {
	fmt.Println("starting up...")

	file, err := os.Open("pipeline.json")
	if err != nil {
		panic(err)
	}

	var config PipelineConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}

	for path, hook := range config.Hooks {
		fmt.Println("registering", path)
		http.Handle(path, githubAuthenticationWrapper(hook))
	}

	http.HandleFunc("/", root)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// wrap the builds in an authentication wrapper to verify the request came from github
func githubAuthenticationWrapper(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

//
// Github Events
//

type GithubEvent struct {
	Ref        string
	Repository GithubRepository
	Pusher     GithubPusher
}

type GithubRepository struct {
	URL      string
	CloneURL string `json:"clone_url"`
}

type GithubPusher struct {
	Name  string
	Email string
}

//
// Define the Build Config
//

type PipelineConfig struct {
	Hooks map[string]HookConfig
}

type HookConfig struct {
	Jobs []JobConfig
}

type JobConfig struct {
	Path        string
	Container   string
	Environment map[string]string
	Volumes     map[string]string
}

func (hook HookConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event GithubEvent
	json.NewDecoder(r.Body).Decode(&event)
	if err := jobSequence(event.Repository.CloneURL, hook.Jobs); err != nil {
		http.Error(w, "job failed", http.StatusInternalServerError)
		return
	}
}

func jobSequence(repository string, jobs []JobConfig) error {

	for _, job := range jobs {
		if err := jobExecute(repository, job); err != nil {
			return err
		}
	}

	return nil
}

func jobExecute(repository string, job JobConfig) error {

	args := []string{"run"}

	args = append(args,
		"-e", "REPOSITORY="+repository,
		"-e", "DOCKER_USERNAME="+dockerUsername,
		"-e", "DOCKER_PASSWORD="+dockerPassword)

	for key, value := range job.Environment {
		args = append(args, "-e", key+"="+value)
	}

	for key, value := range job.Volumes {
		args = append(args, "-v", key+":"+value)
	}

	args = append(args, job.Container)

	if err := execute("docker", args...); err != nil {
		return err
	}

	return nil
}

func execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
