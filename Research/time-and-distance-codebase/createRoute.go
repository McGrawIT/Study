package time_and_distance

import (
	"fmt"
	"net/http"
	"encoding/json"
)

/*	Create Route from Edges, Nodes ( Route Legs ) in one pass ( O(n) )
	
	Call ComputeCurrentNodeDistances () for each Node.  Thie computes
	all Route Distance information ( Node-Node, Node-Asset, Total Route,
	Asset-Route Leg, Split Leg ( Left, Right ) )

	Subseuent /Distance calls will use this "Computed Once" information

	Calls to /Time endpoints will use the Ordered Route ( Distances )
	to calculate Asset to Destination Times

	Route can be reused by subsequent calls w/o a Route

	Aseet and Destination are Optional, their presence ( or absence )
	simply alters the Result Set
*/

var	(

	allAssets 		[]ClosestNodesByAsset

	nodeResultSet		[]NodeInstance
	currentRouteResult 	OverallRoute

	BestRouteETA 		float64
	TotalTraverseLength 	float64
)

func createRoute ( rw http.ResponseWriter, 
		   routeStart, routeEnd Location, 
		   Asset, Destination Location ) ( []routeNodeDistance ) {

	fmt.Println ( "Route Start:", routeStart, "Route End:", routeEnd )
	fmt.Println ( "Primary Destination:", Destination, "Asset:", Asset )

	route := []routeNodeDistance{}		// Initalize Route

	RouteLengths = []float64{}
	RouteGeoLengths = []float64{}

//	Start Overall Route with Node that does not have a LeftNode
//	End Overall Route with Node that does not have a RightNode

	currentNode := routeStart
	previousNode := routeStart

	nextNode := EdgeNodesLL [ routeStart ].RightNode

	LeftEndNode = routeStart
	RightEndNode = routeEnd
	AssetDestination = Destination


	ac := InputRoute.AirportCode
	rID := InputRoute.RouteID
	rn := InputRoute.RouteName

	fmt.Println ( "Airport Code:", ac, "Route ID:", rID, "Name:", rn )

	nodeResultSet = []NodeInstance{}

	theRoute.RouteStart = LeftEndNode
	theRoute.RouteEnd = RightEndNode

//	-----------
//	ROUTE NODES
//	-----------

//	Walk the Linked EdgeNodes to create Route Nodes

	printLoopNext ( "First[P|C|N]", previousNode, currentNode, nextNode )

	newNode := NodeInstance{}

	allAssets = []ClosestNodesByAsset{}

	BestRouteETA = 0.0
	TotalTraverseLength = 0.0

	debug := 0
	for {

//	    -----------------------------
//	    Build the Route Node Instanes
//	    -----------------------------

	    NodeAssetIndex.Node = currentNode

//	    Node Location
//	    -------------

	    newNode.NodeLocation = currentNode

//	    Node-specific Detail
//	    --------------------

	    currentEdge := EdgeNodesLL[ currentNode ]

	    newNode.NodeType = currentEdge.NodeType

	    newNode.MinTraverse = currentEdge.MinTraverse
	    newNode.NodeDelayLeft = currentEdge.LeftNodeDelay

//	    newNode.DistanceMinimum = currentEdge.DistanceMinimum
	    newNode.SpeedLimit = currentEdge.MaximumSpeed

//	    Nodes to Left and Right of Current Node; End Nodes
//	    --------------------------------------------------

	    newNode.LeftNode = previousNode
	    newNode.RightNode = nextNode

	    newNode.LeftEndNode = LeftEndNode
	    newNode.RightEndNode = RightEndNode

//	    Distances
//	    ---------
//	    <<--- To Left ( Previous ) Node ---

	    toLeft := adjacentNodesLength ( previousNode, currentNode )
	    toGeoLeft := adjacentGeoNodesLength ( previousNode, currentNode )

	    newNode.ToLeftNode = toLeft
	    newNode.ToLeftGeoNode = toGeoLeft

	    RouteLengths = append ( RouteLengths, toLeft )
	    RouteGeoLengths = append ( RouteGeoLengths, toGeoLeft )

//	    Maintain distance to Left End ( Trailing Distance )
//	    Current Node to Left End ( Leftmost ) Node ( Route Start )

	    LeftEndDistance = LeftEndDistance + toLeft 
	    newNode.ToLeftEndNode = LeftEndDistance

//	    --- To Right ( Next ) Node --->>

	    toRight := adjacentNodesLength ( currentNode, nextNode )
	    toGeoRight := adjacentGeoNodesLength ( currentNode, nextNode )

	    newNode.ToRightNode = toRight
	    newNode.ToRightGeoNode = toGeoRight
 
//	    Based on Maximum Speed across the Leg and Leg Length, compute
//	    the shartest time to traverse the Leg

	    FastestLegTraverse := 0.0
	    if newNode.SpeedLimit > 0.0 {
		FastestLegTraverse = toRight / newNode.SpeedLimit
	    } else {
		fmt.Println ( "ROute Leg Speed Limit: ZERO" )
	    }
	    newNode.FastestLegTraverse = FastestLegTraverse

//	    Maintain Total Length of Traverse ( Route Length )

	    fmt.Println ( "**********************************************" )
	    fmt.Println ( "Speed Limit:", newNode.SpeedLimit )
	    fmt.Println ( "Distance to Right:", toRight )
	    fmt.Println ( "Fastest Time:", FastestLegTraverse )
	    fmt.Println ( "**********************************************" )

	    BestRouteETA = BestRouteETA + FastestLegTraverse
	    TotalTraverseLength = TotalTraverseLength + toRight

	    fmt.Println ( "Total Route Length:", TotalTraverseLength )
	    fmt.Println ( "Fastest:", BestRouteETA )
	    fmt.Println ( "**********************************************" )

	    newNode.ToRightNode = toRight
	    newNode.ToRightGeoNode = toGeoRight

//	    Create next entry in Ordered Route

	    nodeToRight := routeNodeDistance{}
	    nodeToRight.Node = currentNode
	    nodeToRight.Distance = toRight
	    nodeToRight.FastestLegTraverse = FastestLegTraverse
	    nodeToRight.GeoDistance = toGeoRight

	    route = append ( route, nodeToRight )	// Route Order

//	    Logging / Reporting

	    adjacent ( previousNode, currentNode, nextNode, toLeft, toRight )
	    adjacent ( previousNode, currentNode, nextNode, toGeoLeft, toGeoRight )

/*	    ---------------------------------------------------------
	    Do the Asset Processing on thie Node before moving to the
	    next one.  The Asset map is tied to a single Node
*/
//	    The distance ( between the Nodes ) does not change 
//	    Use the unique Asset Speed for each Asset
//	    ---------------------------------------------------------

	    CurrentNodeAssetIndex := NodeAsset{}

	    CurrentNodeAssetIndex.Node = currentNode

	    for NextAsset, _ := range Assets {

		incorporateAsset ( previousNode, currentNode, nextNode, NextAsset, newNode )
	    }

//	    -----------------------------------------------
//	    Add Current Node to Nodes[] and Node Result Set
//	    -----------------------------------------------
		
	    Nodes [ currentNode ] = newNode
	    nodeResultSet = append ( nodeResultSet, newNode )

	    showAdded ( previousNode, newNode  )

//	    -------------------------------------------------------
/*	    Once the Current Node is complete, and there are more 
	    Nodes to compute, set Previous, Current, and Next Nodes
//	    -------------------------------------------------------
*/

//	    First, use last two Nodes to calculate Distance to Leg

	    if currentNode != previousNode {

		legDistance ( previousNode, currentNode )
	    }

	    if currentNode == RightEndNode { break } 

	    if debug >= TotalRouteEdges { break }
	    debug++

	    previousNode = currentNode
	    currentNode = nextNode

//	    At Right End, there is no Next Node

	    if currentNode == RightEndNode {

		nextNode = RightEndNode

	    } else { nextNode = EdgeNodesLL [ currentNode ].RightNode }

	    printLoopNext ( "Next[P|C|N]", previousNode, currentNode, nextNode )
	}


	fmt.Println ( "=================================================" )
	fmt.Println ( "Total Length:", TotalTraverseLength )
	fmt.Println ( "Best Route Traverse:", BestRouteETA )
	fmt.Println ( "=================================================" )
	

//	Store each Asset's Closest Node time & distance detail
//	It will facilitate future Asset route decisions

	collectClosestNodes ()

/*	Once Route is built, walk in reverse to compute "to Route End"
	information.  ( It is not possible to know beond the immediate
	Next Node while building Route from Start to End )
*/
	ComputeReverseRouteValues () 

/*	With a completed Route ( all Time and Distance values computed ),
	create a Route Result Set for Overall Result Set
*/
/*

	nodeResultSet = []NodeInstance{}

	revisitNode := Nodes [ LeftEndNode ]
	nextNode = LeftEndNode

	for {

	    nodeResultSet = append ( nodeResultSet, revisitNode )

	    if nextNode == RightEndNode { break }

	    nextNode = revisitNode.RightNode
	    revisitNode = Nodes [ nextNode ]
	}
*/

	fmt.Println ( "-------------- [ Node, Asset ] Pairs -----------------" )

	fmt.Println ( "Nodes:", len ( EdgeNodesLL ), "Assets:", len ( Assets ) )
	for node, _ := range NodeAssets {

	    fmt.Println ( node.Node, node.Asset, ":", NodeAssets [ node ] )
	}
	
	fmt.Println ( "-----------------------------------" )
	fmt.Println ( "Route Lengths, Linear:", RouteLengths )
	fmt.Println ( "Route Lengths, Geo:", RouteGeoLengths )
	fmt.Println ( "-----------------------------------" )

	NodeAssetTimes := []NodeAssetDetails{}
	NodeAssetTime := NodeAssetDetails{}

	for key, TimeDetails := range NodeAssets {

	    NodeAssetTime.Node = key.Node
	    NodeAssetTime.Asset = key.Asset
	    NodeAssetTime.TimeDetails = TimeDetails

	    NodeAssetTimes = append ( NodeAssetTimes, NodeAssetTime )
	}
	currentRouteResult.NodeAssetTimes = NodeAssetTimes 

	fmt.Println ( "-----------------------------------" )

//	Compute Asset -> Destination Route(s)

	for _, Asset := range Assets { 

	    getAssetRoute ( Asset.AssetLocation, Asset.Destination )
	}

/*	Need to add Asset Result Set to final output

	AssetDetail		[]AssetResult	`"json:"Assets"`
	RouteTime		float64		`"json:"RouteTime"`
*/

	routeResultSet := assembleResultSet ( route )

	result, err := json.MarshalIndent ( routeResultSet , "", "  " )

	if err != nil { fmt.Println ( "JSON Marshall Failed:", err ) }

	fmt.Fprintf ( rw, string ( result ) ) 
//
	return route
}

func 	printLoopNext ( which string, previous, current, next Location ) {

	bars := "===================================="
	bars = bars + bars

	fmt.Println ( "" )
	fmt.Println ( bars )
	fmt.Println ( which, "[", previous, "|", current, "|", next, "]" )
	fmt.Println ( bars )
	fmt.Println ( "" )
}


func	adjacent ( 	previous, current, next Location, 
			Left, Right float64 ) {

	bars := "------------------------------------"
	bars = bars + bars

	l := "<===" 
	m := "===" 
	r := "===>"

	fmt.Println ( "" )
	fmt.Println ( bars )
	fmt.Println ( previous, l, Left, m, current, m, Right, r, next )
	fmt.Println ( bars )
	fmt.Println ( "" )
}

func	showAdded ( previousNode Location, newNode NodeInstance ) {

	bars := "--------------------------------------------------------" 

	fmt.Println ( bars )
	fmt.Println ( "Previous Node:", previousNode, Nodes [ previousNode ] )

	fmt.Println ( bars )
	fmt.Println ( "Added:", newNode )
	fmt.Println ( "Total:", len ( Nodes ), "Nodes" )
	fmt.Println ( bars )
}


func	legDistance ( previous, current Location ) {

//	Temporarily Disable... chase other bugs first
//	It works ( at least it used to ), but reducing moviing parts

	return

	firstAsset := AssetInstance{}

	for _, Asset := range Assets {

	    firstAsset = Asset
	    break
	}

	fmt.Println ( "Calculate Leg Distances from", previous, "to", current )

	DistanceToLeg ( Nodes [ previous ], Nodes [ current ], firstAsset )

	fmt.Println ( "Distance from Asset to Route Leg:", Asset_Leg )
	fmt.Println ( "Distance to Left Node @ Split:", Base_Left )
	fmt.Println ( "Distance to Right Node @ Split:", Base_Right )

	currentAssetInstance := Assets[ current ]

	currentAssetInstance.ToSegment = Asset_Leg
	currentAssetInstance.ToLeftLeg = Base_Left
	currentAssetInstance.ToRightLeg = Base_Right

	Assets[ current ] = currentAssetInstance
}
