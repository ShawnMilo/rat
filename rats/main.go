package main

import (
	"fmt"
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
	user, host, message string
	stamp               time.Time
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
		user:    clean(r.PostForm["user"]),
		host:    clean(r.PostForm["host"]),
		message: clean(r.PostForm["message"]),
		stamp:   time.Now(),
	}
	addMessage(msg)
	log.Printf("%s@%s: %q\n", msg.user, msg.host, msg.message)
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
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n", m.stamp.Format("2006-01-02 15:04:05"), m.user, m.host, m.message)
	}
	fmt.Fprintln(w, "</table></body></html>")
}
