package main

import (
	"net/http"
)

func handlePingGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	msg := EmptyObject{}
	writeObjectAsJson(w, msg)
}

type EmptyObject struct{}
