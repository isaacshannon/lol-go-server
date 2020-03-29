package main

import "encoding/json"
import "log"

func lg(data map[string]string){
	marshalled, err := json.Marshal(data)
	if err != nil {
		lgError(err)
		return
	}
	log.Println(string(marshalled))
}

func lgError(err error) {
	log.Println(`{"error":`+err.Error()+"}")
}
