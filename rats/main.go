package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	user, host, message string
	stamp               time.Time
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "bad method (POST only)", 405)
	}
	r.ParseForm()
	m := Message{
		user:    clean(r.PostForm["user"]),
		host:    clean(r.PostForm["host"]),
		message: clean(r.PostForm["message"]),
		stamp:   time.Now(),
	}
	log.Printf("%s@%s: %q\n", m.user, m.host, m.message)
}

// Just join for now. Maybe more later.
func clean(val []string) string {
	return strings.Join(val, " ")
}
