package main

import (

//	cl "github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/clients"
//	"github.build.ge.com/aviation-predix-common/vcap-support.git"

	tad "github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/time-and-distance"
	"github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/conversions"


	"github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/server"

//	tad "time-and-distance"
//	"conversions"
//	"server"

	"fmt"

	"net/http"
	"strings"

//	"os"
//	"runtime"
//	"time"
)


/*	Early ideas for Distance Service Endpoints 

	/Distance/Destination		// Asset to Destination
	/Distance/Asset			// Asset Origina to Destination
	/Distance/Route/Assets		// >1 Asset to Destination
	/Distance/Route			// Route ( Start to Finish ) 
	/Distance/Segment		// Segment Distance ( Point-to-Point )

	/Distance/GreatCircle		// Segment Great Circle ( P2P )

	/Time/Destination		// Asset to Destination
	/Time/Route			// Start to FInish
	/Time/Segment			// Max or given Velocity


//	More ( and same / similar to above ) Endpoints for Distance Service

	Each endpoint 
.	Each endpoint returns ordered routes
	If Velocity is given, Time and Distance will be returned
	If not, we coule return the "best possible" time ( use Leg Maxes )

	/Route/Assets/Time
	/Route/Assets/Distance

	/Asset/Routes/Time
	/Asset/Routes/Distance

	/Order
	/Route
	/Asset

	/Distance/Order { one or more Routes }
	/Distance/Route

	/Distance ( one or more routes )
	/Time ( one or more routes )
	
	/Convert

	Support calculations around <x,y,z> and Origin <x,y,z>

	Origin / Delta / Units ( start with X,Y,Z ( add Lat / Long conversion )

*/

func main() {

//	-----------------------------------
//	Time and Distance Service Endpoints
//	-----------------------------------

//	Time requests require a Route and Asset ( for Speed )
//	This relies on the same functionality as the Distance services, so
//	the Result Sets will also include all Distance Result Sets

	http.HandleFunc ( "/api/v1/Distance",		tad.RouteDistance )
	http.HandleFunc ( "/api/v1/distance",		tad.RouteDistance )
	http.HandleFunc ( "/Distance",			tad.RouteDistance )
	http.HandleFunc ( "/distance",			tad.RouteDistance )

	http.HandleFunc ( "/api/v1/Time",		tad.RouteDistance )
	http.HandleFunc ( "/api/v1/time",		tad.RouteDistance )
	http.HandleFunc ( "/Time",			tad.RouteDistance )
	http.HandleFunc ( "/time",			tad.RouteDistance )

	http.HandleFunc ( "/Distance/Destination",	tad.RouteDistance )
	http.HandleFunc ( "/Time/Destination",		tad.RouteDistance )

	http.HandleFunc (  "/Distance/Route",		tad.RouteDistance )
	http.HandleFunc (  "/Distance/RouteSegment",	tad.RouteDistance )

	http.HandleFunc (  "/Time/Route",		tad.RouteDistance )
	http.HandleFunc (  "/Time/RouteSegment",	tad.RouteDistance )

//	-----------
//	Conversions
//	-----------

	http.HandleFunc (  "/api/v1/Convert",		conversions.Convert )
	http.HandleFunc (  "/api/v1/convert",		conversions.Convert )
	http.HandleFunc (  "/Convert",			conversions.Convert )
	http.HandleFunc (  "/convert",			conversions.Convert )

//	------
//	Status
//	------

//	http.HandleFunc ("/Info", tad.InfoRequest )

	http.HandleFunc ("/api/v1/Ping", HandlePing )
	http.HandleFunc ("/api/v1/ping", HandlePing )
	http.HandleFunc ("/Ping", HandlePing )
	http.HandleFunc ("/ping", HandlePing )

	http.HandleFunc ("/api/v1/Info", server.HandleInfo )
	http.HandleFunc ("/api/v1/info", server.HandleInfo )
	http.HandleFunc ("/Info", server.HandleInfo )
	http.HandleFunc ("/info", server.HandleInfo )


	fmt.Println ( "---------------------------------------" )
	fmt.Println ( "Time & Distance Service Server Starting" )
	fmt.Println ( "---------------------------------------" )

	tad.Time_and_Distance_Service_Server ()
}



func HandlePing(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	// Always set content type and status code
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	message := []string{
		"{ }",
	}

// 	And write your response to w

//	if DEBUG { fmt.Printf("Ping Command\n%s", strings.Join(message, "")) }

	fmt.Fprintf(w, strings.Join(message, ""))

}

/*
func HandleInfo(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	vcapAppMap, _ := vcap.LoadApplication()
	guid := vcapAppMap.ID
	uri := GetResponsePath()
	predixSpace := GetPredixSpace()

	// Always set content type and status code
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// And write your response to w
	var ts = time.Now().UTC()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	host_name := cl.GetServiceHostName("uaa_url")

	oauth2_state := "Disabled"
	if OAUTH2_ENABLED {
		oauth2_state = "Enabled"
	}

	ma := float32(mem.Alloc) / 1024.0 / 1024.0                    //  convert to MB
	//rma := float32(dbg.GetRuntimeMaxAllocMem()) / 1024.0 / 1024.0 //  convert to MB
	ta := float32(mem.TotalAlloc) / 1024.0 / 1024.0               //  convert to MB
	sm := float32(mem.Sys) / 1024.0 / 1024.0                      //  convert to MB
	fm := float32(mem.Sys-mem.Alloc) / 1024.0 / 1024.0            //  convert to MB

	message := []string{

		fmt.Sprintf("Received Info on ... %s \n", ts.Format(time.RFC850)),
		fmt.Sprintf("Time and Distance Microservice - Copyright GE 2015-2016\n\n"),
		fmt.Sprintf("Predix Space: %s\n", predixSpace),
		fmt.Sprintf("Instance GUID: %s\n", guid),
		fmt.Sprintf("Instance IP: %s\n", os.Getenv("CF_INSTANCE_ADDR")),
		fmt.Sprintf("Instance #: %s\n", os.Getenv("CF_INSTANCE_INDEX")),
		fmt.Sprintf("Started: %s\n", StartTime.Format(time.RFC850)),
		fmt.Sprintf("Uptime:  %s\n", ts.Sub(StartTime).String()),
		fmt.Sprintf("\n"),
		//fmt.Sprintf("System Operation\n"),
		//fmt.Sprintf("  Number Usable CPUS:     %d\n", dbg.GetRuntimeNumCpus()),
		//fmt.Sprintf("  Number Go Threads:      %d\n", dbg.GetRuntimeGoThreads()),
		//fmt.Sprintf("  Number cGo Threads:     %d\n", dbg.GetRuntimeCGoThreads()),
		//fmt.Sprintf("  Max Number Go Threads:  %d\n", dbg.GetRuntimeMaxGoThreads()),
		//fmt.Sprintf("  Max Number cGo Threads: %d\n", dbg.GetRuntimeMaxCGoThreads()),
		//fmt.Sprintf("\n"),
		fmt.Sprintf("Memory Stats\n"),
		fmt.Sprintf("  Alloc:       %f  (MB allocated and not yet freed)\n", ma),
		//fmt.Sprintf("  Max Alloc:   %f  (MB allocated and not yet freed)\n", rma),
		fmt.Sprintf("  Total Alloc: %f  (MB allocated (even if freed))\n", ta),
		fmt.Sprintf("  Sys:         %f  (MB obtained from system)\n", sm),
		fmt.Sprintf("  Free Memory: %f  (MB Sys - Alloc)\n", fm),
		fmt.Sprintf("  Lookups:     %d  (number of pointer lookups)\n", mem.Lookups),
		fmt.Sprintf("  Mallocs:     %d  (number of mallocs)\n", mem.Mallocs),
		fmt.Sprintf("  Frees:       %d  (number of frees)\n", mem.Frees),
		fmt.Sprintf("\n"),
		fmt.Sprintf("Service State\n"),
		//fmt.Sprintf("  PostgreSQL:     %-s\n", db_state),
		fmt.Sprintf("  UAA:            %-s\n", oauth2_state),
		//fmt.Sprintf("  Oauth Version:  %-s\n", au.GetVersion()),
		//fmt.Sprintf("  Logger Version: %-s\n", dbg.GetVersion()),
		//fmt.Sprintf("\n"),
		fmt.Sprintf("Configuration Variables\n"),
		fmt.Sprintf("  URI:  %s\n", uri),
		fmt.Sprintf("  PORT: %s\n", os.Getenv("PORT")),
        fmt.Sprintf("  UAA_HOST_NAME:       %s\n", host_name),
	}

	if DEBUG { fmt.Printf("Info Command\n%s", strings.Join(message, "")) }
	fmt.Fprintf(w, strings.Join(message, ""))

}

*/


