package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Payload int32  `json:"payload"`
}

func ErrResponse(code int, w http.ResponseWriter, msg string) {
	type Error struct {
		Error string `json:"error"`
	}

	JsonResponse(code, w, Error{
		Error: msg,
	})

	log.Println(msg)
}

func JsonResponse(code int, w http.ResponseWriter, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("An error occure in parsing json ", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
