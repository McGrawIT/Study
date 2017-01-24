package time_and_distance

import (
	"fmt"
)

func getAssetRoute ( asset, destination Location ) {

	fmt.Println ( "Getting:", asset, "--->", destination, "Distance" )
	destinationRoute := DestinationRoute{}

	legNode := RouteLeg{}

	routeLeg := AssetRouteLeg{}
	routeLegs := []AssetRouteLeg{}

//	First Route Leg is Asset --> Closest Node
//	Subsequent Legs are Node --> Node ( toward Route Start or Route End )

	closestNode := Assets[ asset ].ClosestNode

	currentNode := asset	
	nextNode := closestNode

//	Route starts with Asset -> Closest Node
//	Special Case ( Setup and Load )

	legNode.FromNode = asset
	legNode.ToNode = closestNode

	NodeAssetIndex.Node = nextNode 
	NodeAssetIndex.Asset = asset

	routeLegTimes := NodeAssets[ NodeAssetIndex ]	// for Asset --> Node

	fmt.Println ( "Node Asset Index:", NodeAssetIndex )
	fmt.Println ( "Node Asset Timess:", routeLegTimes )

	destinationRoute.Asset = routeLegTimes.Asset 
	destinationRoute.Destination = routeLegTimes.Destination
	destinationRoute.Speed = routeLegTimes.Speed	

	legLength := routeLegTimes.AssetToNode 
	legTime := routeLegTimes.AssetToNodeTime 

	routeLeg.Leg = legNode
	routeLeg.Time = legTime
	routeLeg.Length = legLength

	routeLegs = append ( routeLegs, routeLeg )


//	Initialize Route Time and Distance ( Length ) Totals

	routeTime := routeLeg.Time
	routeLength := routeLeg.Length

//	Determine Direction to Destination ( walk Left or Right nodes )

	headedToStart := towardStart ( Assets[ asset ], Nodes[ destination ] )

	fmt.Println ( "Heading toward Start Node:", headedToStart )

//	Set Variables for Closest Node --> Next Node ( to Left or Right )
//	We are already @ first Node ( Closets Node ), so use the same
//	Route Leg information ( just node times vs. asset time )

	if headedToStart {		// Use Left Nodes

	    nextNode = Nodes[ closestNode ].LeftNode
	    legLength = Nodes[ closestNode ].ToLeftNode
	    legTime = routeLegTimes.NodeToLeftNodeTime

	} else {			// Use Right Nodes

	    nextNode = Nodes[ closestNode ].RightNode
	    legLength = Nodes[ closestNode ].ToRightNode
	    legTime = routeLegTimes.NodeToRightNodeTime
	}

	for {

//	    If the Closest Node was also the Destination, stop @ the one
//	    Route Leg added above and ignore the first Node <-> Node Leg

	    if closestNode == destination { break }

//	    Load Route Leg 

	    legNode.FromNode = currentNode
	    legNode.ToNode = nextNode

	    routeLeg.Leg = legNode
	    routeLeg.Length = legLength
	    routeLeg.Time = legTime

	    routeLegs = append ( routeLegs, routeLeg )

	    fmt.Println ( "Leg:", legNode, "Length:", legLength, "Time:", legTime )

	    routeTime = routeTime + legTime
	    routeLength = routeLength + legLength


	    if currentNode == destination { break }	// Just did last Node

//	    Walk Ordered Route ( get next node ( left or right ) )

	    currentNode = nextNode

	    if headedToStart {		// Use Left Nodes

		nextNode = Nodes[ currentNode ].LeftNode
		legLength = Nodes[ currentNode ].ToLeftNode
		legTime = routeLegTimes.NodeToLeftNodeTime

	    } else {			// Use Right Nodes

		nextNode = Nodes[ currentNode ].RightNode
		legLength = Nodes[ currentNode ].ToRightNode
		legTime = routeLegTimes.NodeToRightNodeTime
	    }

	    NodeAssetIndex.Node = currentNode 
	    routeLegTimes := NodeAssets[ NodeAssetIndex ]


	    fmt.Println ( "Node Asset Index:", NodeAssetIndex )
	    fmt.Println ( "Node Asset Timess:", routeLegTimes )
	    fmt.Println ( "----------------------------------------------" )
	}

	destinationRoute.RouteLegs = routeLegs
	destinationRoute.RouteTime = routeTime
	destinationRoute.RouteLength = routeLength

	currentAsset := Assets[ asset ]

	currentAsset.RouteDuration = routeTime
	currentAsset.RouteLength = routeLength
	currentAsset.Route = destinationRoute

	Assets[ asset ] = currentAsset

}


