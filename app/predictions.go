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
	img := r.FormValue("imgBase64")
	x0 := r.FormValue("x0")
	x1 := r.FormValue("x1")
	y0 := r.FormValue("y0")
	y1 := r.FormValue("y1")

	err := validateImage(img)
	if err != nil {
		return predictResponse{}, err
	}

	resp, err := http.PostForm(
		"http://league-nodeport-service/predict",
		url.Values{
			"imgBase64": {img},
			"x0": {x0},
			"x1": {x1},
			"y0": {y0},
			"y1": {y1},
		})
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
	if l > 350000 || l < 200000 {
		log.Println(fmt.Sprintf("invalid img size: %d", l))
		return errors.New("invalid img size")
	}

	return nil
}
