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

/*
	Likely Endpoints? ( Or, is some of this the internal setup of the SWIM Data Ingest )

	Ingesting the Source Streams in real-time
	The taps to each Source that only get the data
	** Handles all direct connections
	** How is it "informed" of totally new Data Source?
	** Create one ( at least ) per SWIM Source

	Subscription "Manager"
	** Handle / Manage Requests
	** Create Pattern Detectors?
	** Launch / Manage Pattern Detectors ( is this where the topic response happens? )
	** OR, does a detection get passed to Responders

	What are the "subscribe to" arguments / fields?
	** Topic, sub-topic ?
	** Valid Topics table ( where they fit in the hierarchy ( level ) )
	** Maybe they send in a topic=value ( topic could be <1 level, and the value is the endpoint
		( e.g., /airport/carrier/delta gets all airline data @ all airports for Delta? )
		( maybe it's /airport/delta, or /ATL/delta ... does that work? )
		( seems like that would not work; it requires a lot of specialized endpoints, right? )

		Maybe it's not an endpoint, but level1/level2/ ... are "fixed" as L1 = airport
		This is definitely unclear!
	**

	Monitors ( Detectors )
	** Some of what is described above?
	** Pattern can have a "must appear by" attribute

	History ( Persist Source Data )
	** Denormalized
	** Aggregated

	Analytics
	** Long-running Pattern Detection ( or persistent filter )
	** and Management
	** Real-time ( one-off ) Queries against historical content


	Responder ( delivers / ensures delivery )


	"Register" Callbacks ( Let Solace know endpoints for each Content Type ( Weather, Cancels, Airport Stats, ... )
	This is one endpoint for registering the function that handles the specific GET
	OR, does each "GET" provide the callback?  ( Per request versus by content type )

	"Retrieve" ( One for each Content Type ( what a registered callback will handle the Response Body )

		GET Weather
		GET Airline Performance
		GET Flight Data

	**	Content Data Models ( one per Response Body )
		Not the endpoint, but the expected return structure

		What data to we start with?
		What are the Sources w/i SWIM?
		What is the structure of that data?

	"Establish Connection" ( Client opens connections?  Identify yourself to Solace SWIM? )

 	"Monitor / Set Alert / Subscribe  ( POST with Subscription Detail )

 	** Topic ( one or more "interested" )
 	** Frequency ( including live streaming ),
 	** Pattern w/i Topic ( like only send me content from the Topic when these conditions are true )

 	** Do they provide an agent / endpoint / queue to feed data into?
 	** Can more than one "agent" ( endpoint be subscribed to the same feed )
 	** Maybe the patterns w/i a topic can be unique to subscriber ( just >1 call above )

 */

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

//	Do we need the same authorizations as Fly Dubai?

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
