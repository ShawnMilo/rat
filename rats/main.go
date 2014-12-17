package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

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
	log.Printf("%q@%q: %q\n", clean(r.PostForm["user"]), clean(r.PostForm["host"]), clean(r.PostForm["message"]))
}

func clean(val []string) string {
	return strings.Join(val, " ")
}
