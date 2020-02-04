package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	SetupUsers()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/capture", serveCapture)
	http.HandleFunc("/predict", servePredict)
	http.HandleFunc("/findmap", serveFindMap)


	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "index.html"))
}

func serveCapture(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "capture.html"))
}

func servePredict(w http.ResponseWriter, r *http.Request) {
	pred, err := retrievePrediction(r)
	if err != nil {
		errorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pred)
}

func serveFindMap(w http.ResponseWriter, r *http.Request) {
	pred, err := retrieveMap(r)
	if err != nil {
		errorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pred)
}
