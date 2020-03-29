package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	SetupUsers()
	loadGCP()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHero)
	http.HandleFunc("/capture", serveCapture)
	http.HandleFunc("/predict", servePredict)
	http.HandleFunc("/blog", serveBlog)
	http.HandleFunc("/log", serveLog)
	http.HandleFunc("/test", serveTest)


	lg(map[string]string{"msg":"Listening..."})
	http.ListenAndServe(":8080", nil)
}

func serveBlog(w http.ResponseWriter, r *http.Request) {
	lg(map[string]string{"msg":"serving blog"})
	http.ServeFile(w, r, filepath.Join("templates", "blog.html"))
}

func serveHero(w http.ResponseWriter, r *http.Request) {
	lg(map[string]string{"msg":"serving landing"})
	http.ServeFile(w, r, filepath.Join("templates", "hero.html"))
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	lg(map[string]string{"msg":"serving test"})
	http.ServeFile(w, r, filepath.Join("templates", "test.html"))
}

func serveCapture(w http.ResponseWriter, r *http.Request) {
	lg(map[string]string{"msg":"serving capture"})
	http.ServeFile(w, r, filepath.Join("templates", "capture.html"))
}

func servePredict(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lg(map[string]string{"msg":"serving predict"})
	pred, err := retrievePrediction(r)
	if err != nil {
		lgError(err)
		errorResponse(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	lg(map[string]string{"prediction time": time.Since(start).String()})
	json.NewEncoder(w).Encode(pred)
}

func serveLog(w http.ResponseWriter, r *http.Request) {
	lg(map[string]string{"msg":"serving log"})
	err := r.ParseForm()
	if err != nil {
		lgError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Each key has an array of values because each key can be declared more than once in the post form.
	values := make(map[string]string)
	for key, val := range r.PostForm {
		values[key] = val[0]
	}
	values["msg"] = "client log"

	lg(values)
	w.WriteHeader(http.StatusOK)
}
