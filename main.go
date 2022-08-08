package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	usage = `
/start : start sending 20 log/s for 1 minute
/stop  : stop sending logs
/*     : give this usage
`
)

type app struct {
	end time.Time
}

func (a *app) start(w http.ResponseWriter, r *http.Request) {
	a.end = time.Now().Add(1 * time.Minute)
	w.WriteHeader(200)
}

func (a *app) stop(w http.ResponseWriter, r *http.Request) {
	a.end = time.Now()
	w.WriteHeader(200)
}

func (a *app) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(usage))
	w.WriteHeader(200)
}


func (a *app) bg() {
	for {
		if time.Now().After(a.end) {
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
		end: time.Now(),
	}
	go a.bg()

	router := mux.NewRouter()
	router.HandleFunc("/start", a.start)
	router.HandleFunc("/stop", a.stop)
	router.NotFoundHandler = http.HandlerFunc(a.defaultHandler)
	router.HandleFunc("/", a.stop)
	fmt.Printf("listening on :8080\n")
	fmt.Printf("%s", usage)
	panic(http.ListenAndServe(":8080", router))
}
