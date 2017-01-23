package time_and_distance

import (
	"fmt"
	"os"
	"encoding/json"
	"net/http"

//	au "github.build.ge.com/AviationRecovery/go-oauth.git"
)

var	port		string

func Time_and_Distance_Service_Server ( ) {

	port = os.Getenv ( "PORT" )
	if port == "" { port = "4545" }

	fmt.Println ( "--------------------------------------------" )
	fmt.Println ( "Time and Distance Listening on Port:", port )
	fmt.Println ( "--------------------------------------------" )

	http.ListenAndServe ( ":"+port, nil )
}

var	AssetSpeed	float64		// Overall Speed for Route


/**************************************************************/

func RouteDistance ( rw http.ResponseWriter, req *http.Request ) {

	var	routeParameter 		Route

	decoder := json.NewDecoder ( req.Body )
	err := decoder.Decode ( &routeParameter )
	if err != nil { }

	InputRoute = routeParameter

	displayRouteLoadStart ( routeParameter )

//	Modify later to support >1 Asset

	AssetSpeed = routeParameter.Speed

	fmt.Println ( "----------------------------" )

	TotalRouteEdges = len ( routeParameter.Segments )

	fmt.Println ( "Input:", len ( routeParameter.Segments ), "Segments" )

	for _, segment := range routeParameter.Segments {

		fmt.Println ( "Route Leg:", segment )
	}
	fmt.Println ( "----------------------------" )

	InputAssets = []InputAsset{}

	LoadPostBodyAssets ( routeParameter )

	createSegmentList ( routeParameter )
	getRouteEnds ( )

	fmt.Println ( "Route [ Start:", routeStart, "End:", routeEnd, "]" )

	PrimaryAssetLocation = routeParameter.AssetLocation
	PrimaryAssetDestination = routeParameter.Destination
	PrimaryAssetSpeed = routeParameter.Speed 

	LoadSegments ( rw )

	fmt.Println ( "================================================" )
	fmt.Println ( "--------------------- DONE ---------------------" )
	fmt.Println ( "================================================" )

}

/**************************************************************/

func	displayRouteLoadStart ( routeParameter Route ) {

	fmt.Println ( "----------------------------" )
	fmt.Println ( "Name:", routeParameter.RouteName )
	fmt.Println ( "Route Origin:", routeParameter.FullRouteStart )
	fmt.Println ( "Asset Location:", routeParameter.AssetLocation )
	fmt.Println ( "Destination:", routeParameter.Destination )
	fmt.Println ( "Node Defaults:", routeParameter.NodeDefaults )
	fmt.Println ( "Route Defaults:", routeParameter.DefaultRoute )
	fmt.Println ( "----------------------------" )
}
