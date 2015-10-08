package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"
)

func main() {
	flag.Parse()
	message := strings.Join(flag.Args(), " ")
	params := url.Values{}
	params.Set("message", message)
	host, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	user, err := user.Current()
	var username string
	if err != nil {
		log.Println(err)
		username = ""
	} else {
		username = user.Username
	}

	server := os.Getenv("RATSERVER")
	if server == "" {
		server = "localhost"
	}
	port := os.Getenv("RATPORT")
	if port == "" {
		port = "8000"
	}

	params.Set("host", host)
	params.Set("user", username)
	resp, err := http.Post("http://"+server+":"+port, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Println("status:", resp.Status)
	}
}
