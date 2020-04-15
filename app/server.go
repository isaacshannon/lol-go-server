package main

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	SetupUsers()
	loadGCP()
	//show lines on logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//Init Sentry
	err := sentry.Init(sentry.ClientOptions{Dsn: "https://9ff1155e0ae94b668ea3a713d438bae6@o370311.ingest.sentry.io/5196659"})
	if err != nil {log.Fatalf("sentry.Init: %s", err)}
	defer sentry.Flush(2 * time.Second)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHero)
	http.HandleFunc("/capture", serveCapture)
	http.HandleFunc("/predict", servePredict)
	http.HandleFunc("/instructions", serveInstructions)
	http.HandleFunc("/log", serveLog)
	http.HandleFunc("/test", serveTest)

	log.Println("Listening...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func serveInstructions(w http.ResponseWriter, r *http.Request) {
	log.Println("serving blog")
	http.ServeFile(w, r, filepath.Join("templates", "instructions.html"))
}

func serveHero(w http.ResponseWriter, r *http.Request) {
	log.Println("serving landing")
	http.ServeFile(w, r, filepath.Join("templates", "hero.html"))
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	log.Println("serving test")
	http.ServeFile(w, r, filepath.Join("templates", "test.html"))
}

func serveCapture(w http.ResponseWriter, r *http.Request) {
	log.Println("serving capture")
	http.ServeFile(w, r, filepath.Join("templates", "capture.html"))
}

func servePredict(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("serving predict")
	pred, err := retrievePrediction(r)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println(fmt.Sprintf("prediction time: %s", time.Since(start).String()))
	err = json.NewEncoder(w).Encode(pred)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func serveLog(w http.ResponseWriter, r *http.Request) {
	log.Println("serving log")
	err := r.ParseForm()
	if err != nil {
		sentry.CaptureException(err)
		errorResponse(w, err)
		return
	}

	// Each key has an array of values because each key can be declared more than once in the post form.
	values := make(map[string]string)
	for key, val := range r.PostForm {
		values[key] = val[0]
	}
	values["msg"] = "client log"

	log.Println(values)
	w.WriteHeader(http.StatusOK)
}
