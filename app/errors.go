package main

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, err error) {
	data := map[string]string{"error": err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(data)
}
