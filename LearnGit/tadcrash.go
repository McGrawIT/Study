package main

import (
	"fmt"
//	"github.build.ge.com/aviation-intelligent-airport/time-and-distance.git/tad_svc"
	"tad_svc"
	"net/http"
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

	http.HandleFunc ( "/Distance",			tad_svc.RouteDistance )
	http.HandleFunc ( "/Time",			tad_svc.RouteDistance )

	http.HandleFunc ( "/Distance/Destination",	tad_svc.RouteDistance )
	http.HandleFunc ( "/Time/Destination",		tad_svc.RouteDistance )

	http.HandleFunc (  "/Distance/Route",		tad_svc.RouteDistance )
	http.HandleFunc (  "/Distance/RouteSegment",	tad_svc.RouteDistance )

	http.HandleFunc (  "/Time/Route",		tad_svc.RouteDistance )
	http.HandleFunc (  "/Time/RouteSegment",	tad_svc.RouteDistance )

//	-----------
//	Conversions
//	-----------

	http.HandleFunc (  "/Convert",			conversions.Convert )

//	------
//	Status
//	------

	http.HandleFunc ("/Info", tad_svc.InfoRequest )
	http.HandleFunc ("/Ping", tad_svc.Ping )

	http.HandleFunc ("/TestLoad", tad_svc.LoadSegmentTesting )

	fmt.Println ( "---------------------------------------" )
	fmt.Println ( "Time & Distance Service Server Starting" )
	fmt.Println ( "---------------------------------------" )

	tad_svc.Time_and_Distance_Service_Server ()
}
