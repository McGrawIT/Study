package main

import (
	"fmt"
	"github.build.ge.com/aviation-intelligent-airport/configuration-manager-svc/config"
	"github.build.ge.com/aviation-intelligent-airport/configuration-manager-svc/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"
)

const (
	DEBUG      bool   = false

	// Trailing slashes go against common convention.
	// Forcing the appender to use a beginning slash results in much more readable paths.
	base_path = "/api/v1"
)

var (
	OAUTH2_ENABLED bool = false
	responsePath   string
	startTime      time.Time
)

func main() {

	fmt.Println("[---------------------------------]")
	fmt.Println("[ Starting SWIM Ingestion Service ]")
	fmt.Println("[---------------------------------]")

	fmt.Printf("DEBUG: %t\n", DEBUG)

	startTime = time.Now().UTC()

	if DEBUG {
		fmt.Println("Configuring REST Server Interface at ", startTime.Format(time.RFC850))
	}

	var err error
	err = config.LoadFromVcap()
	if err != nil {
		fmt.Println(err)
		fmt.Println("The service failed to load its configuration from VCAP.")
		fmt.Println("Attempting to load the configuration from file, instead.")
		err = config.LoadFromFile("configuration-manager-svc.config.json")
		if err != nil {
			fmt.Println(err)
			fmt.Println("The service failed to load its configuration from file.")
			fmt.Println("This failure is not recoverable. Please check your configuration sources.")
			os.Exit(1)
		}
	}

	// Register a handler for each route pattern
	router := mux.NewRouter()

	// Add a trivial handler for INFO
	router.Methods("GET").Path(base_path + "/info").HandlerFunc(handleInfoGet)
	router.Methods("GET").Path(base_path + "/ping").HandlerFunc(handlePingGet)


	// This is not an awesome thing to have. Find a way around it, like manual DB administration, even.
	router.Methods("PATCH").Path(base_path + "/db").HandlerFunc(handleDbPatch)

	disable_oauth := os.Getenv("DISABLE_OAUTH")
	if strings.ToLower(disable_oauth) == "true" && strings.ToLower(util.GetPredixSpace()) == "ia-dev" {
		OAUTH2_ENABLED = false
	}

	attachProfiler(router, OAUTH2_ENABLED)

	// Get Cloud Foundry assigned port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Port number is invalid")
		port = "8080"
	}

	// Determine response path.
	uri := util.GetServiceHostName("swim_ingestion_url")
	if uri == "" {
		responsePath = "localhost:" + port + base_path
	} else {
		responsePath = uri + base_path
	}

	if DEBUG {
		fmt.Println("Starting REST Server Interface for", uri)
		fmt.Printf("response_path = %s\n", responsePath)
		fmt.Printf("Port = %s\n", port)
	}

	// Start listening on the configured port.
	// ListenAndServe is not expected to return, so we wrap it in a log.Fatal
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func attachProfiler(router *mux.Router, oauth2_enabled bool) {
	if oauth2_enabled {
		router.PathPrefix("/debug/pprof").Handler(http.HandlerFunc(pprof.Index))
		router.PathPrefix("/debug/pprof/cmdline").Handler(http.HandlerFunc(pprof.Cmdline))
		router.PathPrefix("/debug/pprof/profile").Handler(http.HandlerFunc(pprof.Profile))
		router.PathPrefix("/debug/pprof/symbol").Handler(http.HandlerFunc(pprof.Symbol))

		// Manually add support for paths linked to by index page at /debug/pprof/
		router.PathPrefix("/debug/goroutine").Handler(pprof.Handler("goroutine"))
		router.PathPrefix("/debug/heap").Handler(pprof.Handler("heap"))
		router.PathPrefix("/debug/threadcreate").Handler(pprof.Handler("threadcreate"))
		router.PathPrefix("/debug/block").Handler(pprof.Handler("block"))

	} else {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		// Manually add support for paths linked to by index page at /debug/pprof/
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}
}
