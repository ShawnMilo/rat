package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var messages []*Message
var lock = sync.RWMutex{}
var maxMessages int = 1000

type Message struct {
	User, Host, Message string
	Stamp               time.Time
}

func init() {
	log.SetOutput(os.Stdout)
	messages = make([]*Message, 0, maxMessages)
}

func main() {
	http.HandleFunc("/log", history)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "bad method (POST only)", 405)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form data", 400)
		return
	}
	msg := &Message{
		User:    clean(r.PostForm["user"]),
		Host:    clean(r.PostForm["host"]),
		Message: clean(r.PostForm["message"]),
		Stamp:   time.Now(),
	}
	addMessage(msg)
	log.Printf("%s@%s: %q\n", msg.User, msg.Host, msg.Message)
}

// Just join for now. Maybe more later.
func clean(val []string) string {
	return strings.Join(val, " ")
}

func addMessage(msg *Message) {
	lock.Lock()
	defer lock.Unlock()
	messages = append(messages, msg)
	if len(messages) >= maxMessages {
		var tenth int = maxMessages / 10
		m := make([]*Message, 0, maxMessages)
		old := messages[maxMessages-tenth:]
		for i := 0; i < tenth; i++ {
			m[i] = old[i]
		}
		messages = m
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()
	first := len(messages) - 1
	if first < 0 {
		fmt.Fprintln(w, "no data")
		return
	}
	last := first - 100
	if last < 0 {
		last = 0
	}
	fmt.Fprintln(w, "<html><head><title>Rat Log</title></head><body>")
	fmt.Fprintln(w, "<table><tr><th>stamp</th><th>user</th><th>host</th><th>message</th></tr>")

	// Mon Jan 2 15:04:05 MST 2006
	for i := first; i >= last; i-- {
		m := messages[i]
		t, err := template.New("log").Parse("<tr><td>{{.Stamp}}</td><td>{{.User}}</td><td>{{.Host}}</td><td>{{.Message}}</td></tr>\n")
		if err != nil {
			http.Error(w, "internal error", 500)
			return
		}

		err = t.Execute(w, m)
		if err != nil {
			http.Error(w, "internal error", 500)
		}
	}
	fmt.Fprintln(w, "</table></body></html>")
}
