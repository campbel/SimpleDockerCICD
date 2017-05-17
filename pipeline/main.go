package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	buildContainer = "campbel/pipeline-build:latest"
)

var (
	dockerUsername = getSecret("DOCKER_USERNAME")
	dockerPassword = getSecret("DOCKER_PASSWORD")
	githubSecret   = getSecret("GITHUB_SECRET")
)

func getSecret(key string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	data, err := ioutil.ReadFile("/run/secrets/" + key)
	if err != nil {
		fmt.Println("couldn't find secret", key)
		return ""
	}
	return strings.TrimSpace(string(data))
}

const rootHTML = `
<html>
	<head>
		<link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">	
		<style>
			body {
				font-family: 'Roboto', sans-serif;
			}
			.title {
				text-align: center;
				margin: 40px 0px 20px;
				color: rgba(50,50,50,1);
			}
		</style>
	</head>
	<body>
		<h1 class="title">Simple Docker CI/CD</h1>
	</body>
</html>
`

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

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, rootHTML)
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

// wrap the builds in an authentication wrapper to verify the request came from github
func githubAuthenticationWrapper(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

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
	if err := jobSequence(hook.Jobs); err != nil {
		http.Error(w, "job failed", http.StatusInternalServerError)
		return
	}
}

func jobSequence(jobs []JobConfig) error {

	for _, job := range jobs {
		if err := jobExecute(job); err != nil {
			return err
		}
	}

	return nil
}

func jobExecute(job JobConfig) error {

	args := []string{"run"}

	args = append(args,
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
