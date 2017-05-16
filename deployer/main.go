package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("starting up...")
	http.HandleFunc("/", root)
	http.HandleFunc("/push", pushHandler)
	http.ListenAndServe(":80", nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

type GithubEvent struct {
	Ref        string
	Repository GithubRepository
	Pusher     GithubPusher
}

type GithubRepository struct {
	URL string
}

type GithubPusher struct {
	Name  string
	Email string
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	var event GithubEvent
	json.NewDecoder(r.Body).Decode(&event)
	fmt.Println(event)

	if err := execute("docker", "pull", "campbel/app"); err != nil {
		fmt.Println("error:", err)
	}

	if err := execute("docker", "rm", "-f", "app"); err != nil {
		fmt.Println("error:", err)
	}

	if err := execute("docker", "run", "--name", "app", "-d", "-p", "8080:80", "campbel/app"); err != nil {
		fmt.Println("error:", err)
	}
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
