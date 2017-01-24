package main

import (
	"github.build.ge.com/aviation-intelligent-airport/swim-ingestion-svc/binding"
	"github.com/jackmanlabs/errors"
	"log"
)

func main() {
	ver := binding.Version()
	log.Println(ver)

	err := binding.Initialize()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}


}
