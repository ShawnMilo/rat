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
	}else{
        username = user.Username
    }
	params.Set("host", host)
	params.Set("user", username)
	resp, err := http.Post("http://aoeus.com:8000", "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Println("%s", resp.Status)
	}
}
