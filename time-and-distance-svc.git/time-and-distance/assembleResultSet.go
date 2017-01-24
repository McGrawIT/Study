package time_and_distance

import (
	"fmt"
)

var	(

	routeResultSet 	[]OverallRoute
	resultSet 	OverallRoute
)
	

func assembleResultSet ( route []routeNodeDistance ) ( routeResultSet []OverallRoute ) {

	routeResultSet = []OverallRoute{}
	resultSet = OverallRoute{}

//	------------------------------------------------------------------

	resultSet.AirportCode	= InputRoute.AirportCode
	resultSet.RouteID	= InputRoute.RouteID
	resultSet.RouteName	= InputRoute.RouteName

	resultSet.DefaultRoute	= InputRoute.DefaultRoute
	resultSet.NodeDefaults	= InputRoute.NodeDefaults


	resultSet.RouteStart	= LeftEndNode
	resultSet.RouteEnd	= RightEndNode

	resultSet.BestETA	= BestRouteETA
	resultSet.AverageSpeed	= TotalTraverseLength / BestRouteETA 

//	------------------------------------------------------------------


//	Linear Route Length
//	-------------------

	totalRouteLength = 0.0
	for _, legLength := range RouteLengths {

	    totalRouteLength = totalRouteLength + legLength
	}

	resultSet.RouteLength = totalRouteLength
	theRoute.RouteLength = totalRouteLength

//	Geo Route Length
//	----------------

	totalRouteLength = 0.0
	for _, legLength := range RouteGeoLengths {

	    totalRouteLength = totalRouteLength + legLength
	}

	resultSet.GeoLength = totalRouteLength
	theRoute.GeoLength = totalRouteLength

//	------------------------------------------------------------------

//	resultSet.NodeDistances = RouteDistances 
//	theRoute.NodeDistances = RouteDistances 

//	resultSet.OrderedRoute = route
	theRoute.OrderedRoute = route
	theRoute.TotalNodes = len ( Nodes )

//	------------------------------------------------------------------


/*	Once Route is built, walk in reverse to compute "to Route End"
	information.  ( It is not possible to know beond the immediate
	Next Node while building Route from Start to End )
*/
	ComputeReverseRouteValues () 

/*	With a completed Route ( all Time and Distance values computed ),
	create a Route Result Set for Overall Result Set
*/

	nodeResultSet = []NodeInstance{}

	revisitNode := Nodes [ LeftEndNode ]
	nextNode := LeftEndNode

	for {

	    nodeResultSet = append ( nodeResultSet, revisitNode )

	    if nextNode == RightEndNode { break }

	    nextNode = revisitNode.RightNode
	    revisitNode = Nodes [ nextNode ]
	}


//	------------------------------------------------------------------


//	Need these?  ( Asseembled elsewhere; need more or less ? )

	resultSet.RouteNodes = nodeResultSet	// []NodeInstance

	resultSet.AssetsClosest = allAssets	// []ClosestNodesByAsset

	resultSet.BasicRoute = theRoute


//	------------------------------------------------------------------

	fmt.Println ( "-----------------------------------" )

//	Does this remain here, or folded into Assets, or ?

	NodeAssetTimes := []NodeAssetDetails{}
	NodeAssetTime := NodeAssetDetails{}

	for key, TimeDetails := range NodeAssets {

	    NodeAssetTime.Node = key.Node
	    NodeAssetTime.Asset = key.Asset
	    NodeAssetTime.TimeDetails = TimeDetails

	    NodeAssetTimes = append ( NodeAssetTimes, NodeAssetTime )
	}
	resultSet.NodeAssetTimes = NodeAssetTimes 

	fmt.Println ( "-----------------------------------" )


//	------------------------------------------------------------------


//	Flatten Assets map[] into []AssetInstance for Result Set

/*	See what should be added / removed
	Rather than a simple append, I could create a Asset Result record
*/

	routeAssets := []AssetInstance{}

	for _, Asset := range Assets { 

	    routeAssets = append ( routeAssets, Asset )
	}
	resultSet.RouteAssets = routeAssets 


//	------------------------------------------------------------------


	routeResultSet = append ( routeResultSet, resultSet )

	return
}

