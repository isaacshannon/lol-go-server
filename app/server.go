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
	http.HandleFunc("/", serveHero)
	http.HandleFunc("/capture", serveCapture)
	http.HandleFunc("/predict", servePredict)
	http.HandleFunc("/blog", serveBlog)
	http.HandleFunc("/log", serveLog)


	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

func serveBlog(w http.ResponseWriter, r *http.Request) {
	log.Println("serving blog")
	http.ServeFile(w, r, filepath.Join("templates", "blog.html"))
}

func serveHero(w http.ResponseWriter, r *http.Request) {
	log.Println("serving landing")
	http.ServeFile(w, r, filepath.Join("templates", "hero.html"))
}

func serveCapture(w http.ResponseWriter, r *http.Request) {
	log.Println("serving capture")
	http.ServeFile(w, r, filepath.Join("templates", "capture.html"))
}

func servePredict(w http.ResponseWriter, r *http.Request) {
	log.Println("serving predict")
	pred, err := retrievePrediction(r)
	if err != nil {
		log.Println(err)
		errorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pred)
}

func serveLog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Each key has an array of values because each key can be declared more than once in the post form.
	values := make(map[string]string)
	for key, val := range r.PostForm {
		values[key] = val[0]
	}

	marshalled, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(string(marshalled))
	w.WriteHeader(http.StatusOK)
}
