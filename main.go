package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/convyr/http/que"
)

var (
	channel string
	port    string
	synch   bool
)

func init() {
	que.Flags()
	flag.StringVar(&port, "port", "80", "Port to listen on")
	flag.StringVar(&channel, "channel", os.Getenv("CONVYR_CHANNEL"), "channel for convyr communication")
	flag.BoolVar(&synch, "synch", false, "Make request synchronous")
}

func main() {
	flag.Parse()
	q, err := que.New(channel)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()
	h := Handler{q}
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), &h)
	if err != nil {
		log.Fatal(err)
	}
}

type Handler struct {
	*que.Que
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusTeapot)
		log.Print(http.StatusText(http.StatusTeapot))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if synch {
		ret, err := h.Que.Sync(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		}
		w.Write(ret)
		return
	}
	err = h.Que.Async(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
