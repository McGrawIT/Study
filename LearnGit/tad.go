package main

import (
	"fmt"
//	"github.build.ge.com/aviation-intelligent-airport/time-and-distance.git/tad_svc"
	"tad_svc"
	"net/http"
)


Primary POST Body to hit Endpoints:

Route ({ Header }{ Bag of Legs ( Node, Node ){ Node, Node}}  )
Asset { Node, Velocity, Other )
Destination { Detail }( Node )
Node ( Detail ) { X, Y } 

{{{{R{}}{L{}}{L{}}{L}}{{R{}}{L}{L}{L}}}
{{{{A{}}{L{}}{L{}}{L{}}}{{A{}}{L{}}{L{}}{L{}}}}
{{{{D{}}{L{}}{L}{}{L}}{{D{}{}}{L}{L}{L{}}}}

{ "Route:"{
	"F1":"Val1",
	"F2":"Val2",
},{},{},},
{"Legs":{
	{
	"LeftNode":{"X":"38","Y":"22"},{"X":"55","Y":"11"},
	"RightNode":{"X":"38","Y":"22"},{"X":"55","Y":"11"},
}
{Route...}
Asset:{Header}{Location}
Asset:{Header}{Location}
Asset:{Header}{Location}
Destination:{Header}{Location"
Destination:{Header}{Location"



/*	Endpoints for Distance Service

	/Distance/Destination		// Asset to Destination
	/Distance/Route			// Route ( Start to Finish ) 
	/Distance/Routes		// Coult include /Route
	/Distance/Asset			// Needed? /AssetRoute ?
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

	http.HandleFunc ("/Distance",tad_svc.ComputeDistance )
	http.HandleFunc ("/Distance/Calc",tad_svc.ComputeDistance )
	http.HandleFunc ("/Distance/Route",tad_svc.GetRoutes )
	http.HandleFunc ("/Time",tad_svc.ComputeTimeToDestination )

//	---------------------------------
//	Provide Detailed Slot Information
//	---------------------------------

	http.HandleFunc ("/Info", tad_svc.InfoRequest )
	http.HandleFunc ("/Ping", tad_svc.Ping )

	fmt.Println ( "---------------------------------------" )
	fmt.Println ( "Time & Distance Service Server Starting" )
	fmt.Println ( "---------------------------------------" )

	tad_svc.Time_and_Distance_Service_Server ()
}
