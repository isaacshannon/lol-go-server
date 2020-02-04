package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type predictResponse struct {
	Result string `json:"result"`
}
func retrievePrediction(r *http.Request) (predictResponse, error) {
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
		log.Println(err)
		return predictResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return predictResponse{}, err
	}

	pred := predictResponse{}
	err = json.Unmarshal(body, &pred)
	if err != nil {
		log.Println(err)
		return predictResponse{}, err
	}

	return pred, nil
}

type findResponse struct {
	Minimap string `json:"minimap"`
	X0 int64 `json:"x0"`
	X1 int64 `json:"x1"`
	Y0 int64 `json:"y0"`
	Y1 int64 `json:"y1"`

}

func retrieveMap(r *http.Request) (findResponse, error) {
	img := r.FormValue("imgBase64")
	resp, err := http.PostForm(
		"http://league-nodeport-service/findmap",
		url.Values{
			"imgBase64": {img},
		})
	if err != nil {
		log.Println(err)
		return findResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return findResponse{}, err
	}

	pred := findResponse{}
	err = json.Unmarshal(body, &pred)
	if err != nil {
		log.Println(err)
		return findResponse{}, err
	}

	return pred, nil
}
