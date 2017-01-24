package time_and_distance

import (

)

/*	--------------------------------------------------------------
	Common Location along the Route ( aka, Node / Vertex )
	Will likely use <x,y,z> coordinate system with reference point
	--------------------------------------------------------------
*/

type 	Location struct {

	Lat  		float64 		`json:X"`
	Long 		float64 		`json:Y"`

	Height 		int 			`json:Z"`

/*	Other possible Location Dimension Names

	PointX		float64		`json:PointX"`
	PointY		float64		`json:PointY"`

	AxisX		float64		`json:"AxisX"`
	AxisY		float64		`json:"AxisY"`
*/
}

/*	Being replaced by RouteEdge type ( backward-compatibilty until all
	functions using "old" variable names are changed )
*/

/*
type 	Edge 	struct { // Route Leg ( or Segment )

	leftEdgeNode 	Location
	rightNode    	Location

	DirectionLeft  	bool 		// Right Edge to Left Edge Valid
	DirectionRight 	bool 		// Left Edge to Right Edge Valid

	MaxSpeed 	float64 	// Max Speed across Route Leg
	MinTraverse 	float64 	// Minimum Time to traverse Route Leg

	LeftNodeDelay  	float64		// Delay @ Node ( e.g., Traffic Light )
	LeftNodeBuffer 	float64		// Node Blocked, Use Route Leg only
}
*/


//	*** TEMP : Quick Hack to bring back basic POST Reponses

type Leg struct {

	LeftNode    	Location 	`json:"LeftNode"`
	RightNode   	Location 	`json:"RightNode"`
	MaxSpeed 	float64  	`json:"MaxSpeed"`
}

var (
//	InputSegments 		[]Edge

	Nodes             	map[Location]NodeInstance
	Ends              	map[Location]int

	AllEdges          	map[EdgeID]RouteEdge

	OverallRouteBegin 	Location
	OverallRouteEnd   	Location

	LeftEndNode       	Location
	RightEndNode      	Location

	AssetDestination  	Location
)

//	Using both, but could probably use DistanceFromEnd since we only
//	compute one / direction

var (

	LeftEndDistance 	float64 	// Distance from Left End
	LeftEndTime 		float64     	// Time from Left End

	RightEndDistance 	float64 	// Distance from Right End

	RouteLengths 		[]float64
	RouteGeoLengths 	[]float64
)

//	Last-minute hack until >1 Asset Destination is supported

var	PrimaryAssetLocation 		Location
var	PrimaryAssetDestination 	Location
var	PrimaryAssetSpeed 		float64

//	Input: "Bags" of Segments ( Route Legs ) [ Set of Edges ]

var	TotalRouteEdges	int

var	InputRoute	Route

type	Route		struct {

	AirportCode	string		`json:"AirportCode"`
	RouteName	string		`json:"RouteName"`
	RouteID		string		`json:"RouteID"`

	OffCourse	int		`json:"OffCourse"`
	Segments	[]RouteEdge	`json:"Segments"`

	AssetsInRoute	[]InputAsset	`json:"Assets"`

//	Not needed; calculated during ordering of Edges

	FullRouteStart	Location	`json:"FullRouteStart"`
	FullRouteEnd	Location	`json:"FullRouteEnd"`
	RouteReference	Location        `json:"RouteReferencePoint"`

//	Node Defaults

	NodeDefaults	Defaults	`json:"NodeDefaults"`

//	Primary Asset ( will also be in Assets List )

	DefaultRoute	RouteDefaults	`json:"RouteDefaults"`

	AssetLocation	Location	`json:"AssetLocation"`
	Destination	Location	`json:"Destination"`
	Speed 		float64 	`json:"AssetSpeed"`

}

type	RouteDefaults struct {

	AssetLocation	Location	`json:"AssetLocation"`
	RouteStart	Location	`json:"RouteStart"`
	RouteEnd	Location	`json:"RouteEnd"`
	Destination	Location	`json:"Destination"`
	Speed	 	float64 	`json:"Speed"`
}

type	Defaults	struct {

	MaxSpeed	float64		`json:"MaxSpeed"`
	MaxHeight	float64		`json:"MaxHeight"`
	MaxWeight	float64		`json:"MaxWeight"`
	MaxWidth	float64		`json:"MaxWidth"`
	NodeDelays	[]NodeDelay	`json:"NodeDelays"`
}

type	NodeDelay	struct {

	NodeType	string		`json:"NodeType"`
	NodeDelay	float64		`json:"NodeDelay"`
}

type	RouteAsset	struct {

	AssetLocation 		Location	`json:"Location"`
	Destination 		Location	`json:"Destination"`
	Speed 			int 		`json:"Speed"`

	AssetType 		int 		`json:"AssetType"`
	NodeDistanceMax 	int 		`json:"NodeDistanceMax"`
}



//	-------------------------------------------------------------------
//				NODES
//	-------------------------------------------------------------------

type 	NodeInstance struct {

	NodeType     		string   	`json:"NodeType"`
	NodeLocation 		Location 	`json:"NodeLocation"`
//	Zones			[]int		`json:"Zones"`
	Zones 			int 		`json:"Zones"`

	MinTraverse 		float64 	`json:"TraVerseMinimum"`

//	Maximum Speed Asset can travel from Current Node to Right Node
//	Based on Speed Limit and Leg Length, compute Fastest Time across Leg

	SpeedLimit 		float64 	`json:"SpeedLimit"`
	FastestLegTraverse 	float64 	`json:"FastestLegTraverse"`

/*	Location of Node, Distance to its Left, Right Nodes, any Node Delays
	when crossing this Node ( going Left-to-Right and Right-to-Left ), and
	the Maximmum Speed for Assets traveling Left->Right and Right->Left
	do not change, regardless of each Asset's Speed / Speed

	Maximumm Asset Dimensions ( Height, Width and Weight ) are also fixed
	and are directional:  one direction thorugh a Node may have different
	restrictions.

	Max / Min values of Zero ( not supplied ) are taken as no retrictions
*/

//	Current Node --> Left Node

	LeftNode      		Location `json:"LeftNode"`
	ToLeftNode    		float64  `json:"DistanceToPreviousNode"`
	ToLeftGeoNode  		float64  `json:"DistanceToPreviousGeoNode"`
	NodeDelayLeft 		float64  `json:"NodeDelayLeft"`

//	Capture Node Characteristics ( for Future Use, may only want one
//	set for both directions @ Node )

	//MaxSpeedLeft 		float64 `json:"MaxSpeedLeft"`
	//MaxHightLeft    	float64 `json:"MaxHightLeft"`
	//MaxWeightLeft   	float64 `json:"MaxWeightLeft"`
	//MaxWidthLeft    	float64 `json:"MaxWidthLeft"`

//	Current Node --> Right Node

	RightNode      		Location `json:"RightNode"`
	ToRightNode    		float64  `json:"DistanceToNextNode"`
	ToRightGeoNode    	float64  `json:"DistanceToNextGeoNode"`
	NodeDelayRight 		float64  `json:"NodeDelayRight"`

//	Capture Node Characteristics ( for Future Use, may only want one
//	set for both directions @ Node )

	//MaxSpeedRight 	float64 `json:"MaxSpeedRight"`
	//MaxHightRight    	float64 `json:"MaxHightRight"`
	//MaxWeightRight   	float64 `json:"MaxWeightRight"`
	//MaxWidthRight    	float64 `json:"MaxWidthRight"`

/*	Capture Location of two Ends of the Route ( and distance to them )
	Knowing this makes very simple Route Segment calculations, esp. to
	determine distances from an Asset to its Destination, which will
	be a subset of the Overall Route almost every time.  ( Very useful
	since Distances never change for a Route )
*/
	LeftEndNode   		Location `json:"RouteStart"`
	ToLeftEndNode 		float64  `json:"DistanceFromRouteStart"`

	RightEndNode   		Location `json:"RouteEnd"`
	ToRightEndNode 		float64  `json:"DistanceToRouteEnd"`

/*	Added "LeftMost" in addition to Overall Route Start / End to support
	"one way" Edges.  The "end" Node is where an Asset can go no farther
	in that direction ( not supported ( yet ) )

	For Future Use
*/
	//LeftMostNode    	Location `json:"LeftMostNode"`
	//ToLeftMostNode  	float64  `json:"DistanceToLeftMostNode"`
	//RightModeNode   	Location `json:"RightModeNode"`
	//ToRightModeNode 	float64  `json:"DistanceToRightMostNode"`

/*	Assets outside the Range are too far away to consider this Node for
	Route Entry.  Assets within Rnage, but outside the Node Buffer, must
	use an Edge ( too far away to be sure something will not be in the way )

	For Future Use
*/
	//AssetRange      	float64 `json:"AssetRange"`
	//DistanceMinimum 	float64 `json:"DistanceMinimum"`
}

//	-------------------------------------------------------------------
//				NODE / ASSET
//	-------------------------------------------------------------------

/*	Only Node-Node times vary, and are determined by Speed of each Asset
	Time values ( since they vary by Asset ) are maintained in another
	strucutre indexed by < Node, Asset > pairs
*/

var 	NodeAssets 		map[NodeAsset]RouteAssetDetail
var 	NodeAssetIndex 		NodeAsset

type 	NodeAsset struct {

	Node  		Location	`json:"Node"`
	Asset 		Location	`json:"Asset"`
}

/*	Assets affect timing between Nodes along the Routes.  Distances are
	fixed, but traverse times are not.  Adjacent Node times are best held
	at Nodes for optimal Route selection.
*/

type	NodeAssetDetails	 struct {

	Node  		Location		`json:"Node"`
	Asset  		Location		`json:"Asset"`
	TimeDetails	RouteAssetDetail	`json:"TimeDetail"`
}

type 	RouteAssetDetail struct {

	Asset       		Location	`json:"Asset"`
	Destination 		Location	`json:"Destination"`
	Speed    		float64		`json:"Speed"`
	
	AssetToNode     	float64		`json:"AssetToNode"`
	AssetToNodeTime 	float64		`json:"AssetToNodeTime"`
	
	TimeToDestination 	float64		`json:"TimeToDestination"`

//	Time to Nodes ( from this Node ) are tied to Asset-by-Asset Speed

	NodeToLeftNodeTime  	float64		`json:"NodeToLeftNodeTime"`
	NodeToRightNodeTime 	float64		`json:"NodeToRightNodeTime"`

	NodeToLeftEndTime  	float64		`json:"NodeToLeftEndTime"`
	NodeToRightEndTime 	float64		`json:"NodeToRightEndTime"`

	AssetToLeftEndTime  	float64		`json:"AssetToLeftEndTime"`
	AssetToRightEndTime	float64		`json:"AssetToRightEndTime"`
}

type	ClosestNodesByAsset	struct {

	ClosestNode		NodeAsset
	Detail			AssetClosestNode
}

var	AssetFromClosest	map [ NodeAsset ]AssetClosestNode

type	AssetClosestNode	struct {

	Asset			Location	`json:"Asset"`
	ClosestNode		Location	`json:"ClosestNode"`
	Destination		Location	`json:"Destination"`

	ToClosestNode		float64		`json:"ToClosestNode"`
	ToDestination		float64		`json:"ToDestination"`
	ToLeftEnd		float64		`json:"ToLeftEnd"`
	ToRightEnd		float64		`json:"ToRightEnd"`

	TimeToClosest		float64		`json:"TimeToClosestNode"`
	TimeToDestination	float64		`json:"TimeToDestination"`
	TimeToLeftEnd		float64		`json:"TimeToLeftEnd"`
	TimeToRightEnd		float64		`json:"TimeToRightEnd"`
}



//	-------------------------------------------------------------------
//				ASSETS
//	-------------------------------------------------------------------

var 	Assets 		map[Location]AssetInstance

var 	InputAssets 	[]InputAsset

//	Route Asset ( also used as Asset portion of Result Set ( POST Body ) )

type 	AssetInstance struct {

	AssetID   	string 			`json:"AssetID"`
	AssetType 	string 			`json:"AssetType"`

	Speed        	float64  		`json:"Speed"`

	Weight 		int 			`json:"Weight"`
	Height 		int 			`json:"Height"`
	Width  		int 			`json:"Width"`

	AssetLocation   Location 		`json:"AssetLocation"`
	Destination 	Location 		`json:"Destination"`

	Route		DestinationRoute	`json:"RouteToDestination"`
	RouteLength	float64  		`json:"RouteLength"`
	RouteDuration	float64 		`json:"RouteDuration"`


//	Which Zones can this Asset support?  ( Empty set, no restrictionis )
//	Switch to a set of Valid Zones once "good" JSON test data is created

//	ValidZones	[]int			`json:"ValidZones"`
	ValidZones 	int 			`json:"ValidZones"`

//	Assets this far from the Node are "too far" to be considered for
//	by this Node ( e.g., too many path obstacles likely if that far away )

	AwayFromNodeMax 	float64 	`json:"AwayFromNodeMax"`

//	Are these needed?  ( Do they make sense or are implemented elsewhere? )

	DistanceToNodes 	[]NodeDistance        `json:"DistanceToNodes"`

//	Computed ( Based on Asset Location and Destination )
//	Type of Request ( /Distance or /Time affects computation

//	Destinations    	[]DestinationsByAsset `json:"Destinations"`

//	These three as just one Route to Destination
//	Once the above is implemented...
//	... or maybe this is the default ( chorest using closest? )


//	These Distances are Calculated @ each Node, saved @ Node
//	and checked against Asset's Current Lowest

	ClosestNode  		Location 	`json:"ClosestNode"`
	NodeDistance 		float64  	`json:"DistanceToClosestNode"`

	NextClosestNode 	Location 	`json:"NextClosestNode"`
	SecondDistance  	float64  	`json:"DistanceToNextClosest"`

//	ClosestRouteLeg 	Location 	`json:"ClosestRouteLeg"`

// 	Asset to Route Leg Distance ( straight line to triangle Base )
// 	Distance ( going Right ) @ Asset intersection after Route Leg split
// 	Distance ( going Left ) @ Asset intersection after Route Leg split

	ToSegment  		float64 	`json:"AssetToRouteLeg"`
	ToRightLeg 		float64 	`json:"LegSplityRightSide"`
	ToLeftLeg  		float64 	`json:"LegSplityLeftSide"`
}


type	RouteLeg struct {

	FromNode	Location	`json:"FromNode"`
	ToNode		Location	`json:"ToNode"`
}

type	AssetRouteLeg struct {

	Leg		RouteLeg	`json:"Leg"`

	Length		float64		`json:"Length"`
	Time		float64		`json:"Time"`
}

type	DestinationRoute struct {

	Asset       	Location		`json:"Asset"`
	Destination 	Location		`json:"Destination"`
	Speed		float64			`json:"Speed"`

	RouteLegs	[]AssetRouteLeg		`json:"RouteLegs"`
	RouteTime	float64			`json:"Duration"`
	RouteLength	float64			`json:"Length"`
}


//	Route Preferences:  "Time", "Distance", "Cost", "Closest Node"
//	Shortest ( Distance ) and Fastest ( Time ) are always returned

type 	InputAsset struct {

	AssetLocation    	Location 	`json:"Location"`
	Destination 		Location 	`json:"Destination"`
	RoutePreference 	string	 	`json:"RoutePreference"`
	AssetSpeed	    	float64  	`json:"Speed"`

	AssetID   		string 		`json:"AssetID"`
	AssetType 		string 		`json:"AssetType"`

//	AssetDimensions		Dimensions	`json:"Dimensions"`

	Weight 			int 		`json:"Weight"`
	Height 			int 		`json:"Height"`
	Width  			int 		`json:"Width"`

//	ValidZones		[]int		`json:"ValidZones"`
	ValidZones      	int     	`json:"ValidZones"`

	AwayFromNodeMax 	float64 	`json:"AwayFromNodeMax"`
}

type 	Dimensions struct {

	Weight 		int 		`json:"Weight"`
	Height 		int 		`json:"Height"`
	Width  		int 		`json:"Width"`
}

//	Deprecated ( for now ); only supporting one unique Destination / Asset

var 	AssetDestinations 	map[DestinationIndex]DestinationsByAsset

//	-------------------------------------------------------------------
//				ROUTES
//	-------------------------------------------------------------------

/*	Store and use one or more Routes.  Loaded when part of the POST Body.
	Remains active for later use by Route reference only.
*/

type 	SiteRoutes 	map[RouteKey]OverallRoute

type 	RouteKey struct {

	AirportCode 	string  	 `json:"AirportCode"`
	RouteID     	string  	 `json:"RouteID"`
	RouteName   	string   	`json:"RouteName"`
	RouteStart  	Location 	`json:"RouteStart"`
	RouteEnd    	Location 	`json:"RouteEnd"`
}

/*	POST Reponse structure;  Using embedded types, the composite Result Set
	will be used for ALL responses.  From trivial ( straight line distance )
	to complex ( mulitple Assets headed to multiple Destinations at varying
	Speeds along shortest, fastest, cheapest, and "closest" Routess )

	Additional available Route / Node / Edge / Asset informations is
	included for current and future application use cases / analytics
*/

type 	OverallRoute struct {

	AirportCode	string          	`json:"AirportCode"`
	RouteID		string          	`json:"RouteID"`
	RouteName	string          	`json:"RouteName"`
	DefaultRoute	RouteDefaults		`json:"RouteDefaults"`
	NodeDefaults	Defaults		`json:"NodeDefaults"`
	RouteReference	Location        	`json:"RouteReferencePoint"`

	BasicRoute	TheRoute		`json:"TheRoute"`

	RouteStart	Location        	`json:"RouteStart"`
	RouteEnd	Location        	`json:"RouteEnd"`
	RouteLength	float64         	`json:"RouteLength"`
	GeoLength	float64         	`json:"GeoLength"`

	BestETA		float64         	`json:"BestETA"`
	AverageSpeed	float64        		`json:"AverageSpeed"`
	TotalTravese	float64        		`json:"TotalTravese"`
	
	OffRoute	Location        	`json:"OffRoute"`
	RouteTimes	[]RouteTime     	`json:"RouteTimes"`

	RouteNodes	[]NodeInstance  	`json:"RouteNodes"`
	RouteAssets 	[]AssetInstance 	`json:"Assets"`

	NodeAssetTimes 	[]NodeAssetDetails	 `json:"NodeAssetTimes"`

	AssetsClosest	[]ClosestNodesByAsset	`json:"ClosestNodesByAsset"`    
}

var	theRoute	TheRoute

type	TheRoute struct {

	RouteStart	Location  		`json:"Start"`
	RouteEnd	Location  		`json:"End"`
	RouteLength	float64      		`json:"Length"`
	GeoLength	float64      		`json:"GeoLength"`

	TotalNodes	int	  		`json:"TotalNodes"`
	OrderedRoute	[]routeNodeDistance	`json:"Ordered"`
}

type	routeNodeDistance struct {

	Node			Location 	`json:"RouteNode"`
	Distance		float64 	`json:"DistanceToNextNode"`
	FastestLegTraverse	float64 	`json:"FastestLegTraverse"`
	GeoDistance		float64 	`json:"GeoDistanceToNextNode"`
}

//	-------------------------------------------------------------------
//			EDGES ( ROUTE LEGS )
//	-------------------------------------------------------------------

type 	EdgeID struct {

	LeftNode  	Location 	`json:"LeftNode"`
	RightNode 	Location 	`json:"RightNode"`
}

type 	RouteEdge struct {

	LeftNode  		Location 	`json:"LeftNode"`
	RightNode 		Location 	`json:"RightNode"`
	Length    		float64  	`json:"Length"`

// 	Left Node --> Right Node ia always Valid
//	If LeftToRightOnly is false, both Direcitons are Valid ( default )

	OnlyLeftToRight 	bool 		`json:"OnlyLeftToRight"`

/*	Is either ( or both ) ends of this Edge a Node?  ( Many "micro-edges"
	will exist to support non-straight paths between two Nodes. )
*/
	LeftIsNode  		bool 		`json:"LeftIsNode"`
	RightIsNode 		bool 		`json:"RightIsNode"`

/*	Edge characteristics are associated with Left Vertex ( some are Noded )
	since only one end of the Edge is needed for them
*/
//	EdgeInZone		[]int		`json:"LegInZones"`
	EdgeInZone 		int 		`json:"LegInZones"`

//	Route Leg dimension restrictions for each Asset

	MaxAssetWeight 		int 		`json:"MaxAssetWeight"`
	MaxAssetHeight 		int 		`json:"MaxAssetHeight"`
	MaxAssetWidth  		int 		`json:"MaxAssetWidth"`

// 	Maximum Asset Speed across Route Leg ( Edge )
// 	Minimum Time to traverse Route Leg ( Edge )

	MaximumSpeed 		float64 	`json:"MaximumSpeed"`
	SpeedLimit 		float64 	`json:"SpeedLimit"`
	MinTraverse  		float64 	`json:"MinimumTraverse"`

// 	Delay @ Node ( e.g., Traffic Light, Intersection )
// 	Node Blocked unless Asset w/i the "Buffer", Use Route Leg only

	LeftNodeDelay  		float64 	`json:"NodeDelay"`
	NodeType  		string	 	`json:"NodeType"`
	LeftNodeBuffer 		float64 	`json:"NodeBlockedDistance"`
}

type 	DestinationIndex struct {

	Asset       		Location 	`json:"Asset"`
	Destination 		Location 	`json:"Destination"`
}

type 	DestinationsByAsset struct {

	Destinations DestinationIndex

//	Routes ( Shortest, Quickest, Lowest Cost, Using Closest Node )

	ClosestNodeRoute 	[]RouteOption 	`json:"Closest"`
	QuickestRoute    	[]RouteOption 	`json:"Quickest"`
	ShortestRoute    	[]RouteOption 	`json:"Shortest"`
	LowestCostRoute  	[]RouteOption 	`json:"LowestCost"`
}

//	Generic Route ( Nodes, Time, Distance ( Location to Destination ) )

type 	RouteOption struct {

	Route         		[]Location 	`json:"Route"`
	RouteTime     		float64    	`json:"RouteTime"`
	RouteDistance 		float64    	`json:"RouteDistance"`
}

//	Distance to Nodes ( from an Asset )

type 	NodeDistance struct {

	Node           		Location 	`json:"Node"`
	DistanceToNode 		float64  	`json:"DistanceToNode"`
}

//	The Route Time is unique to each Asset's Speed

type RouteTime struct {

	Asset 		Location 		`json:"Asset"`
	Time  		float64  		`json:"RouteTime"`
}
