package time_and_distance

import (
	"fmt"
	"math"
//	"time"
	"net/http"
//	"encoding/json"
)

var	EdgeNodes	map [ Location ]leg
var	EdgeNodesLL	map [ Location ]RouteEdge

type	leg	struct {

	firstNode	Location
	secondNode	Location

	DirectionLeft	bool		// Right Edge to Left Edge Valid
	DirectionRight	bool		// Left Edge to Right Edge Valid

	MaxSpeed	float64		// Max Speed across Route Leg
	MinTraverse	float64		// Minimum Time to traverse Route Leg

	LeftNodeDelay	float64		// Delay @ Node ( e.g., Traffic light )
	LeftNodeBuffer	float64		// Node Blocked, Use Route Leg only
}


var	totalRouteLength 	float64
var	routeStart		Location
var	routeEnd 		Location

func LoadSegments ( rw http.ResponseWriter ) {

	Nodes = make ( map[Location]NodeInstance )
	Assets = make ( map[Location]AssetInstance )

	NodeAssets = make ( map[NodeAsset]RouteAssetDetail )
	AssetFromClosest = make ( map[NodeAsset]AssetClosestNode )

	fmt.Println ( "--------------------------------" )
	fmt.Println ( "--------- Load Assets ----------" )
	fmt.Println ( "--------------------------------" )

	LoadAssets ()

//	Use until >1 Destination and >1 Asset are fully working

	Destination := PrimaryAssetDestination
	Asset := PrimaryAssetLocation


	fmt.Println ( "------------------------------" )
	fmt.Println ( "Edges:", len ( EdgeNodesLL ) )

//	Save Ends of Route ( Caller will indicate Destination )
//	Arbitrarily set the first one to the Start of the Route

	LeftEndDistance = 0
	LeftEndTime = 0
	RightEndDistance = 0

	startType := EdgeNodesLL [ routeStart ].NodeType
	endType := EdgeNodesLL [ routeEnd ].NodeType

	fmt.Println ( "Route[", "Start:", routeStart, "(", startType, ") End:", routeEnd, "(", endType, ") ]" )

	fmt.Println ( "-------------------------------------" )
	fmt.Println ( "Test One-and-done Load of Route Input" )
	fmt.Println ( "-------------------------------------" )

	fmt.Println ( "Airport Code:", InputRoute.AirportCode )
	fmt.Println ( "Route Name:", InputRoute.RouteName )
	fmt.Println ( "Route ID:", InputRoute.RouteID )
	fmt.Println ( "OffCourse:", InputRoute.OffCourse )
	fmt.Println ( "Route Origin:", InputRoute.FullRouteStart )
	fmt.Println ( "Route Node Defaults:", InputRoute.NodeDefaults )
	fmt.Println ( "Route Defaults:", InputRoute.DefaultRoute )
	fmt.Println ( "-------------------------------------" )

	route := createRoute ( rw, routeStart, routeEnd, Asset, Destination ) 
	fmt.Println ( "Route:", route )

	fmt.Println ( "Node<->Node ( Leg ) Distances:", RouteLengths )
	fmt.Println ( "Node<->Node ( Leg ) Geo Distances:", RouteGeoLengths )

	ComputeReverseRouteValues () 

	for Asset, _ := range Assets {

	    closest := Assets [ Asset ].ClosestNode
	    index := NodeAsset{}

	    index.Node = closest		// Closest Node
	    index.Asset = Asset

	    timeToLeftEnd := NodeAssets [ index ].AssetToLeftEndTime

	    fmt.Println ( "Asset:", Asset, "Closest Node:", closest, "Time to Left End:", timeToLeftEnd )
	}


//	Asset to Destination Times

	for asset, _ := range Assets {

//	    Get Asset-Specific Time for the Closest Node

	    closest := Assets [ asset ].ClosestNode

	    index := NodeAsset{}
	    index.Node = closest
	    index.Asset = asset

/*	    Get times for Asset to Closest Node, Closest Node to the
	    Left End Node, and Destination to the Left End Node
*/
	    AssetTime := NodeAssets [ index ].AssetToNodeTime
	    NodeTime := NodeAssets [ index ].NodeToLeftEndTime

//	    Destination to Left End Node

	    Destination := NodeAssets [ index ].Destination

	    index.Node = Destination

	    DestinationTime := NodeAssets [ index ].NodeToLeftEndTime

/*	    The Distance from the Closest Node to the Destination is 
	    calculated from those two Nodes distances to the Left End
*/

//	    Default to Destination toward the Left

	    NodeToDestination := NodeTime - DestinationTime 

//	    Switch if the Destination is to the Right

	    if DestinationTime > NodeTime { 
		NodeToDestination = DestinationTime - NodeTime 
	    }

//	    Asset to Destination is in two parts

	    AssetToDestination := AssetTime + NodeToDestination

	    fmt.Println ( "Asset", asset, "to Destination", Destination, "is", AssetToDestination )

	}

	fmt.Println ( "-----------------------------------------------" )
	fmt.Println ( "-----------------------------------------------" )

	totalRouteLength = 0.0
	for _, legLength := range RouteLengths {

		totalRouteLength = totalRouteLength + legLength
	}
	fmt.Println ( "Total Route Length:", totalRouteLength )
	fmt.Println ( "-----------------------------------------------" )

//	fmt.Println ( "Total Route Time:", LeftEndTime )
	computeRouteTimes ()

	fmt.Println ( "-----------------------------------------------" )

	displayAssetDistances ()

	fmt.Println ( "-----------------------------------------------" )
	fmt.Println ( "-----------------------------------------------" )

/*
	for Asset, _ := range Assets {

		displayAssetToDestination ( Asset, Destination )

	}
*/

}

func	displayAssetToDestination ( Asset, Destination Location ) {


	fmt.Println ( "Assets[", Asset, "] :", Assets [ Asset ] )

	ClosestNode := Assets [ Asset ].ClosestNode
	AssetLeg := Assets [ Asset ].NodeDistance	// Asset to Node

	fmt.Println ( "Closest Node:", ClosestNode )
	fmt.Println ( "Asset Leg:", AssetLeg )

	NodeAssetIndex.Node = ClosestNode
	NodeAssetIndex.Asset = Asset

	fmt.Println ( "Times per Asset:", NodeAssets [ NodeAssetIndex ] )

	AssetLegTime := NodeAssets [ NodeAssetIndex ].AssetToNodeTime

	DestToLE := Nodes [ Destination ].ToLeftEndNode

	NodeAssetIndex.Node = Destination
	
	DestTimeToLE := NodeAssets [ NodeAssetIndex ].NodeToLeftEndTime

	fmt.Println ( "---------------------------------------------" )
	fmt.Println ( "Destination Node:", Nodes [ Destination ] )
	fmt.Println ( "---------------------------------------------" )


	NodeToLE := Nodes [ ClosestNode ].ToLeftEndNode

	NodeAssetIndex.Node = ClosestNode
	
	NodeTimeToLE := NodeAssets [ NodeAssetIndex ].NodeToLeftEndTime

	fmt.Println ( "---------------------------------------------" )
	fmt.Println ( "Closest Node:", Nodes [ ClosestNode ] )
	fmt.Println ( "---------------------------------------------" )

	fmt.Println ( "Asset:", Asset )
	fmt.Println ( "Closest Node:", ClosestNode )
	fmt.Println ( "Destination:", Destination )

	fmt.Println ( "Asset to Closest Node:", AssetLeg )
	fmt.Println ( "Destination to Left End Node:", DestToLE )
	fmt.Println ( "Asset's Closest Node to Left End Node:", NodeToLE )

	fmt.Println ( "Asset to Closest Node Time:", AssetLegTime )
	fmt.Println ( "Destination to Left End Node Time:", DestTimeToLE )
	fmt.Println ( "Asset's Closest Node to Left End Node Time:", NodeTimeToLE )

	NodeToDestinationDistance := DestToLE - NodeToLE
	AssetToDestionDistance := math.Abs ( NodeToDestinationDistance ) + AssetLeg

	fmt.Println ( "---------------------------------------------" )
	fmt.Println ( "Destimation Time to Left End:", DestTimeToLE )
	fmt.Println ( "Node Time to Left End:", NodeTimeToLE )
	fmt.Println ( "Asset Leg Time:", AssetLegTime )

	NodeToDestinationTime := DestTimeToLE - NodeTimeToLE
	AssetTimeToDestion:= math.Abs ( NodeToDestinationTime ) + AssetLegTime

	fmt.Println ( "---------------------------------------------" )
	fmt.Println ( "Asset to Destination Distance:", AssetToDestionDistance )
	fmt.Println ( "Asset to Destination Time:", AssetTimeToDestion )
	fmt.Println ( "---------------------------------------------" )

	RouteNodes := []Location{}

//	This supports the Destionation to the "right" of the Closest Node
//	Must add logic to determine direciion and compute differently
	
	RouteNodes = append ( RouteNodes, Assets [ Asset ].AssetLocation )

	NextNode := Assets [ Asset ].ClosestNode
	fmt.Println ( "NextNode:", NextNode, "Asset:", Asset, "RouteNodes:", RouteNodes, "AssetLocation:", Assets [ Asset ].AssetLocation )

	stop := 0
	for {
		RouteNode := Nodes [ NextNode ].NodeLocation

		RouteNodes = append ( RouteNodes, RouteNode )

		fmt.Println ( "NextNode:", NextNode, "RouteNode:", RouteNode, "RouteNodes:", RouteNodes )

		if NextNode == Destination { break }

		NextNode = Nodes [ RouteNode ].RightNode

		stop++

		if stop > 15 { break }

		fmt.Println ( "(", stop, ")", "NextNode:", NextNode, "Destination:", Destination )
	}

}

func	displayAssetDistances () {

	fmt.Println ( "-------------- ASSETS -------------" )
	fmt.Println ( "" )

	fmt.Println ( "# of Assets:", len ( Assets ) )

	for NextAsset, _ := range Assets {

	    Asset := Assets [ NextAsset ]

	    fmt.Println ( "Location:", Asset.AssetLocation )
	    fmt.Println ( "Destination:", Asset.Destination )
	    fmt.Println ( "Speed:", Asset.Speed )

//	    Get Time to Closest Node

	    nodeAssetIndex := NodeAsset{}
	    nodeAssetIndex.Node = Asset.ClosestNode
	    nodeAssetIndex.Asset = NextAsset

	    nodeTime := NodeAssets [ nodeAssetIndex ].AssetToNodeTime

	    fmt.Println ( "Closest Node:", Asset.ClosestNode, "Distance:", Asset.NodeDistance, "Time to Node:", nodeTime )


//	    Get Time to Next Closest Node

	    nodeAssetIndex = NodeAsset{}
	    nodeAssetIndex.Node = Asset.NextClosestNode
	    nodeAssetIndex.Asset = NextAsset

	    nodeTime = NodeAssets [ nodeAssetIndex ].AssetToNodeTime

	    fmt.Println ( "Next Closest Node:", Asset.NextClosestNode, Asset.SecondDistance, "Time to Node:", nodeTime )

//	    fmt.Println ( "Closest Leg:", Asset.ClosestRouteLeg, Asset.ToSegment )
//	    fmt.Println ( "Leg Split: Going Right (", Asset.ToRightLeg, ") Going Left (", Asset.ToLeftLeg, ")" )
	    fmt.Println ( "" )
	    fmt.Println ( "-----------------------------------" )
	}
}

/*

	AssetRouteStart		// Starts @ Asset Location
	AssetRouteEnd		// Asset Destionation
	AssetRoute		// Legs ( Edges ) of Route to Destination
	AssetSpeed		
	RouteTime		// Asset to Destination travel Time
	RouteDistance		// Distance from Asset to Destination
*/

//	-----
//	TO DO:  Reinstate this logic
//	-----

//	The following form the three sides of the "Asset Triangle"

var 	Left_Right	float64		// Length of Route Leg ( Left<->Right )
var	Asset_Left	float64		// Distance ( Asset <-> Left Node )
var	Asset_Right	float64		// Distance ( Asset <-> Right Node )

//	Use SSS Rule (above) to get length of Asset to Base ( aka, Leg )

var	Asset_Leg	float64		// Asset to Route Leg Distance

var	Base_Left	float64		// Leg Spilt ( Distance to Left Node )
var	Base_Right	float64		// Leg Spilt ( Distance to Right Node )
	
func	DistanceToLeg ( Previous, Current NodeInstance, Asset AssetInstance ) {

/* 	---------------------------------------------------------------------	
	These calculations capture distances to Route Legs.  When a Node has
	a "Buffer" ( something that prevents Assets from taking a straight iine
	to the Node ), or we have "micro edges" ( fine-grained segments that
	allow support for non-straight line routes ).

	The need to get to a Route Leg is critical to using Routes correctly

	Use Current and Previous "Nodes" ( not all segment endpoints are Nodes,
	but all serve the same purpose: connect legs to form the Route ) and 
	each Asset Location to calculate Asset Triangles, and then compute the 
	straight line intersetion from Asset to Leg
 	---------------------------------------------------------------------	
*/
	Left_Right, Asset_Left, Asset_Right = AssetTriangle ( Previous, Current, Asset ) 

	Asset_Leg = GetAssetLegDistance ( Asset_Left, Left_Right, Asset_Right ) 

//	Compare Distance to Segment to Asset's Closest Segment
//	When there is new Closest Segment, compute the Split

//	wasCloser := UpdateClosestSegments ()

	Base_Left, Base_Right = SplitLeg ( Asset_Leg, Asset_Left, Asset_Right )

//	-----
//	TO DO
//	-----

/*	Store information for subswquent Time and Distance calculations 

		Intersection to Left Node ( Left Subsegment )
		Intersection to Right Node ( Right Subsegment )
		Distance from Asset to Leg ( Asset to Leg Segment )

	Asset to Destination decisions will consider ( in addition to
	the two closest Nodes ) both 2-segment ( Asset to Left / Right
	Nodes using the Leg Subsegment ) Route options
*/


//	HIDE1

/*	One or more Assets were supplied with Time and Distance request
	Compute Distance from those Asset(s) to this Node
	The Assets collection ( map[ Asset Location ] Asset Record )
	is populated during handling of the Time & Distance request Body
*/
/*	Asset(s) ( Distance / Time to each as well as the asset-specific times
	to each adjacent Node and Route End )
*/

}

/*	Create similar logic?  ( See Closest Node )
	Maybe save a collection of Closest Segments since there will be
	a great deal of them ( micro-segments )
*/

func	UpdateClosestSegments () (bool ) {

	return true

}
