package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func errResponse(code int, w http.ResponseWriter, msg string) {
	type Error struct {
		Error string `json:"error"`
	}

	log.Fatalf(msg)
	jsonResponse(code, w, Error{
		Error: msg,
	})

}

func jsonResponse(code int, w http.ResponseWriter, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("An error occure in parsing json ", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
