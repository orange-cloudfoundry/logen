package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

type app struct {
	started bool
}

func (a *app) start(w http.ResponseWriter, r *http.Request) {
	a.started = true
	w.WriteHeader(200)
}

func (a *app) stop(w http.ResponseWriter, r *http.Request) {
	a.started = false
	w.WriteHeader(200)
}

func (a *app) bg() {
	for {
		if !a.started {
			time.Sleep(1 * time.Second)
			continue
		}
		for i := 0; i < 20; i++ {
			fmt.Printf("{\"message\": \"here is some log %d\"}\n", i)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	kingpin.Version(version.Print("standard-app"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	a := &app{
		started: true,
	}
	go a.bg()

	router := mux.NewRouter()
	router.HandleFunc("/start", a.start)
	router.HandleFunc("/stop", a.stop)
	panic(http.ListenAndServe(":8080", router))
}
