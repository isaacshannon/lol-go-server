package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type predictResponse struct {
	Predictions [][]interface{} `json:"predictions"`
}
func retrievePrediction(r *http.Request) (predictResponse, error) {
	userID := r.FormValue("userID")
	ip := r.Header.Get("X-FORWARDED-FOR")
	img := r.FormValue("imgBase64")

	err := validateUser(userID)
	if err != nil {
		return predictResponse{}, err
	}
	err = validateUserIP(userID, ip)
	if err != nil {
		return predictResponse{}, err
	}
	err = validateImage(img)
	if err != nil {
		return predictResponse{}, err
	}

	resp, err := http.PostForm(
		"http://league-nodeport-service/predict",
		url.Values{"imgBase64": {img}})
	if err != nil {
		sentry.CaptureException(err)
		return predictResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sentry.CaptureException(err)
		return predictResponse{}, err
	}

	pred := predictResponse{}
	err = json.Unmarshal(body, &pred)
	if err != nil {
		sentry.CaptureMessage(string(body))
		sentry.CaptureException(err)
		return predictResponse{}, err
	}

	// Post prediction tasks
	if userID != "" {
		log.Println("saving image")
		go saveImage(img, userID)
	}

	return pred, nil
}

func validateImage(imgURL string) error {
	if !strings.Contains(imgURL, "data:image/png;base64") {
		log.Println("invalid img format")
		return errors.New("invalid img format")
	}

	l := len(imgURL)
	if l > 500000 || l < 100000 {
		log.Println(fmt.Sprintf("invalid img size: %d", l))
		return errors.New("invalid img size")
	}

	return nil
}
