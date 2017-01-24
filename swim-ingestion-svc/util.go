package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// This sets the content type and serializes the message to the writer.
func writeObjectAsJson(w http.ResponseWriter, msg interface{}) {

	// This is redundant.
	// w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	jsonBuffer := bytes.NewBuffer(nil)
	enc := json.NewEncoder(jsonBuffer)
	enc.SetIndent("", "\t")
	enc.Encode(msg)

	if DEBUG {
		pc, file, line, _ := runtime.Caller(1)
		func_ := runtime.FuncForPC(pc)
		funcChunks := strings.Split(func_.Name(), "/")
		funcName := funcChunks[len(funcChunks)-1]

		log.Printf("Caller:\t%s:%d\t(%s)\n", file, line, funcName)
		log.Printf("Message:\t%s\n", jsonBuffer.Bytes())
	}

	io.Copy(w, jsonBuffer)
}

// This sets the content type and serializes the message to the writer.
func writeRawJson(w http.ResponseWriter, msg []byte) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK) // This is redundant.

	var jsonBuffer *bytes.Buffer = bytes.NewBuffer(msg)

	if DEBUG {
		indentBuffer := bytes.NewBuffer(nil)
		json.Indent(indentBuffer, msg, "", "/t")
		jsonBuffer = indentBuffer

		pc, file, line, _ := runtime.Caller(1)
		func_ := runtime.FuncForPC(pc)
		funcChunks := strings.Split(func_.Name(), "/")
		funcName := funcChunks[len(funcChunks)-1]

		log.Printf("Caller:\t%s:%d\t(%s)\n", file, line, funcName)
		log.Printf("Message:\t%s\n", jsonBuffer.Bytes())
	}

	io.Copy(w, jsonBuffer)
}

func writeError(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	msg := ApiError{Message: err.Error()}

	jsonBuffer := bytes.NewBuffer(nil)
	enc := json.NewEncoder(jsonBuffer)
	enc.SetIndent("", "\t")
	enc.Encode(msg)

	if DEBUG {
		pc, file, line, _ := runtime.Caller(1)
		func_ := runtime.FuncForPC(pc)
		funcChunks := strings.Split(func_.Name(), "/")
		funcName := funcChunks[len(funcChunks)-1]

		log.Printf("Caller:\t%s:%d\t(%s)\n", file, line, funcName)
		log.Printf("Message:\t%s\n", jsonBuffer.Bytes())
	}

	io.Copy(w, jsonBuffer)
}

func write404(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)

	msg := ApiError{Message: fmt.Sprintf("Object could not be found.")}

	jsonBuffer := bytes.NewBuffer(nil)
	enc := json.NewEncoder(jsonBuffer)
	enc.SetIndent("", "\t")
	enc.Encode(msg)

	if DEBUG {
		pc, file, line, _ := runtime.Caller(1)
		func_ := runtime.FuncForPC(pc)
		funcChunks := strings.Split(func_.Name(), "/")
		funcName := funcChunks[len(funcChunks)-1]

		log.Printf("Caller:\t%s:%d\t(%s)\n", file, line, funcName)
		log.Printf("Message:\t%s\n", jsonBuffer.Bytes())
	}

	io.Copy(w, jsonBuffer)
}

type ApiError struct {
	Message string
}
