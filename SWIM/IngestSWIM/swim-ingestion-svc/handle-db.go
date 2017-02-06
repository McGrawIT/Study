package main

import (
	"github.build.ge.com/aviation-intelligent-airport/configuration-manager-svc/data"
	"github.com/jackmanlabs/errors"
	"net/http"
)

func handleDbPatch(w http.ResponseWriter, r *http.Request) {
	err := data.Recreate()
	if err != nil {
		writeError(w, errors.Stack(err))
	}
}
