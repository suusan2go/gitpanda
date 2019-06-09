package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string
)

var port int

var isPrintVersion bool

func init() {
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()
}

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	checkEnv("GITLAB_API_ENDPOINT")
	checkEnv("GITLAB_BASE_URL")
	checkEnv("GITLAB_PRIVATE_TOKEN")
	checkEnv("SLACK_OAUTH_ACCESS_TOKEN")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("gitpanda started: port=%s\n", port)
	http.HandleFunc("/", normalHandler)
	http.ListenAndServe(":"+port, nil)
}

func checkEnv(name string) {
	if os.Getenv(name) != "" {
		return
	}

	log.Printf("[ERROR] %s is required\n", name)
	fmt.Println("")
	printUsage()
	os.Exit(1)
}

func printVersion() {
	fmt.Printf("gitpanda %s, build %s\n", Version, Revision)
}

func printUsage() {
	fmt.Println("[Usage]")
	fmt.Println("  PORT=8000 \\")
	fmt.Println("  GITLAB_API_ENDPOINT=https://your-gitlab.example.com/api/v4 \\")
	fmt.Println("  GITLAB_BASE_URL=https://your-gitlab.example.com \\")
	fmt.Println("  GITLAB_PRIVATE_TOKEN=xxxxxxxxxx \\")
	fmt.Println("  SLACK_OAUTH_ACCESS_TOKEN=xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000 \\")
	fmt.Println("  ./gitpanda")
}

func normalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("It works"))

	case http.MethodPost:
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := strings.TrimSpace(buf.String())

		s := NewSlackWebhook(
			os.Getenv("SLACK_OAUTH_ACCESS_TOKEN"),
			&GitLabURLParserParams{
				APIEndpoint:  os.Getenv("GITLAB_API_ENDPOINT"),
				BaseURL:      os.Getenv("GITLAB_BASE_URL"),
				PrivateToken: os.Getenv("GITLAB_PRIVATE_TOKEN"),
			},
		)
		response, err := s.Request(
			body,
			false,
		)

		if err != nil {
			log.Printf("[ERROR] body=%s, response=%s, error=%v", body, response, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte(response))
	}
}
