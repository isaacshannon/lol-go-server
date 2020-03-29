package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
		lgError(err)
		return predictResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lgError(err)
		return predictResponse{}, err
	}

	pred := predictResponse{}
	err = json.Unmarshal(body, &pred)
	if err != nil {
		lgError(err)
		return predictResponse{}, err
	}

	if userID != "" {
		lg(map[string]string{"msg":"saving image"})
		go saveImage(img, userID)
	}

	return pred, nil
}
