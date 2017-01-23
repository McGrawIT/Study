package time_and_distance

import (
	"fmt"
	"math"

        "github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/conversions"

//	"conversions"
//	"time"
)

/*	Key Distance-Related Functions, including those that are used to make
	Route creation ( esp. Segment ( Route Leg ) selection decisions )

	GetAssetLegDistance ( Asset_Left, Left_Right, Asset_Right )
	SplitLeg ( Asset_Leg, Asset_Left, Asset_Right ) 

	Basic Functions:  Distance, Area, Height
*/

type 	Point 	struct {

	X 	float64
	Y 	float64
}

/*	Compute the length of the hypotenuse between two points.
	Forumula is the square root of (x2 - x1)^2 + (y2 - y1)^2
*/

func ( p Point ) Distance( p2 Point ) float64 {

	first := math.Pow( float64( p2.X - p.X ), 2 )
	second := math.Pow( float64( p2.Y - p.Y ), 2 )

	return math.Sqrt( first + second )
}

func New( x float64, y float64 ) Point {

	return Point{ x, y }
}


func minDistance() {

	p1 := New( 6, 1 )
	p2 := New( 6, 11 )

	dist := p1.Distance( p2 )

	fmt.Println( "Distance", dist )

//	-------------------------------

	Left := New( 7, 16 )			// Left Node
	Right := New( 4, 12 )			// Right Node
	Asset := New( 12, 11 )			// Asset Location

	Left_Right := Left.Distance( Right )	// Node to Node Distance
	Asset_Left := Asset.Distance( Left )	// Asset to Left Node Distance
	Asset_Right := Asset.Distance( Right )	// Asset to Right Node Distance

	fmt.Println ( "Asset to Left Node:", Asset_Left )
	fmt.Println ( "Asset to Right Node:", Asset_Right )
	fmt.Println ( "Left Node to Right Node:", Left_Right )

//	-------------------------------

	Asset_Leg := GetAssetLegDistance ( Asset_Left, Left_Right, Asset_Right ) 
	fmt.Println ( "From Asset to Leg:", Asset_Leg )

	Base_Left, Base_Right := SplitLeg ( Asset_Leg, Asset_Left, Asset_Right )

	fmt.Println ( "Base ( Left:", Base_Left, "Right:", Base_Right, ")" )
}

/*	----------------------------------------------------------------------
	Split the Route Leg into two distances, both from the Intersection of
	the Asset and the Route Leg.  ( One is to the Left Node, the other to
	the Right Node.  Use those two as Node-Node ( Leg ) distances.  ( They
	are the Second Leg of any Route the Asseet takes when chosing Closest
	Segment ( ignore ties ) over one of the Two Closest Nodes.  Using the 
	Two Closest Nodes is (a) way better than the first one, and (b) yields
	more than "good enough" accuracy. )
	----------------------------------------------------------------------
*/

func SplitLeg ( Asset_Leg, Asset_Left, Asset_Right float64 ) ( Base_Left, Base_Right float64 ) {

//	Use A-squared + B-squared = C-squared ( Sides of a Right Triangle )

	Asset_To_Leg := Asset_Leg * Asset_Leg	// Same for both Triangles

//	Left Portion of Base ( Left Split of Leg )

	Asset_To_Left := Asset_Left * Asset_Left
	Base_Left_Sqaured := Asset_To_Left - Asset_To_Leg 

	Base_Left = math.Sqrt ( Base_Left_Sqaured )

//	Right Portion of Base ( Right Split of Leg )

	Asset_To_Right := Asset_Right * Asset_Right
	Base_Right_Squared := Asset_To_Right - Asset_To_Leg 

	Base_Right = math.Sqrt ( Base_Right_Squared )

	return
}


func GetArea ( a, b, c float64 ) ( area float64 ) {

	p := ( a + b + c ) / 2		// Perimeter

	pa := p - a
	pb := p - b
	pc := p - c

	final := p * pa * pb * pc

	area = math.Sqrt( final )
	return
}

func GetHeight ( area, base float64 ) ( height float64 ) {

	height = ( ( 2 * area ) / base )
	return 
}

func GetAssetLegDistance ( Asset_Left, Left_Right, Asset_Right float64 ) ( height float64 ) {

	base := Left_Right
	area := GetArea ( Asset_Left, base, Asset_Right )
	height = GetHeight ( area, base )

	fmt.Println ( "Distance from Asset to Leg:", height )

	return
}




/******************************  MORE NEEDED FUNCTIONS *******************/


func AssetTriangle ( LeftNode, RightNode NodeInstance, AssetNode AssetInstance ) ( Left_Right, Asset_Left, Asset_Right float64 ) {

	LeftLocation := LeftNode.NodeLocation
	RightLocation := RightNode.NodeLocation
	AssetLocation := AssetNode.AssetLocation

	Left := New( LeftLocation.Lat, LeftLocation.Long )	// Left Node Location
	Right := New( RightLocation.Lat, RightLocation.Long )	// Right Node Location
	Asset := New( AssetLocation.Lat, AssetLocation.Long )	// Asset Location

	Left_Right = Left.Distance ( Right )	// Node to Node Distance
	Asset_Left = Asset.Distance ( Left )	// Asset to Left Node Distance
	Asset_Right = Asset.Distance ( Right )	// Asset to Right Node Distance

	fmt.Println ( "Node(", LeftLocation, ") to Node (", RightLocation, ") :", Left_Right )
	fmt.Println ( "Asset(", AssetLocation , ") to Left Node(", LeftLocation,") :", Asset_Left )
	fmt.Println ( "Asset(", AssetLocation , ") to Right Node(", RightLocation,") :", Asset_Right )

	return
}


/* 	Distance from Node to Node ( Segment (Route Leg) Distances )
	Ignoring Height ( for early Versions )

	Left_Right := SegmentLenght ( LegLength )
*/

func Node_AssetLength ( Node, Asset Location ) ( Node_Asset float64 ) {

	NodePoint := New( Node.Lat, Node.Long )		// Node Location
	AssetPoint := New( Asset.Lat, Asset.Long )	// Asset Location

	Node_Asset = NodePoint.Distance ( AssetPoint )	// Node<->Asset Distance

	return
}

func adjacentNodesLength ( LeftNode, RightNode Location ) ( Left_Right float64 ) {

	LeftPoint := New( LeftNode.Lat, LeftNode.Long )
	RightPoint := New( RightNode.Lat, RightNode.Long )

//	Point-to-Point Distance ( Node-to-Node; Asset-to-Node )

	Left_Right = LeftPoint.Distance ( RightPoint )	// Node to Node Distance

	return
}

func adjacentGeoNodesLength ( Left, Right Location ) ( float64 ) {

	leftX := Left.Lat
	leftY := Left.Long 

	rightX := Right.Lat
	rightY := Right.Long 

	geoSegmentLength := conversions.DriveGreatCircle ( leftX, leftY, rightX, rightY )

	return geoSegmentLength
}


/*	Direction used to dtermine if Left or Right Node distances are 
	used to Walk Asset --> Destinaion ( for Route Distance / Time )
*/

func towardStart ( Asset AssetInstance, Destination NodeInstance ) ( bool ) {

/*	A = Length ( Start <-- Destination )
	B = Length ( Start <-- Closest Node )

	If A > B, Destination is Closer to End ( heading toward End Node )
*/

	destinationToStart := Destination.ToLeftEndNode
	closestToStart := Nodes[ Asset.ClosestNode  ].ToLeftEndNode

	if destinationToStart > closestToStart {

	    return false	// Destination --> End
	}
	return true		// Start <-- Destination
}

//	Node <--> Node Length ( not Asset <--> Node )

func routeSegmentLength ( startNode, endNode NodeInstance ) ( length float64 ) {

//	It doesn't matter which NOde is considered the Start / End

	endLength := endNode.ToLeftEndNode
	startLength := startNode.ToLeftEndNode

	length = math.Abs ( endLength - startLength )

	return
}

//	Node <--> Node Length ( not Asset <--> Node )

func destinationDistance ( asset 	AssetInstance, 
			   destination 	NodeInstance ) ( distance float64 ) {

	assetLength := asset.NodeDistance

	closestNode := Nodes[ asset.ClosestNode ]
	routeLegLength := routeSegmentLength ( closestNode, destination )

	distance = assetLength + routeLegLength

	return
}

func UpdateSegmentClosest ( ) { } 

/*	Compare Currnt Node to Asset Distance with Asset's Closest
	and Next Closest Nodes;  Update if Closer
*/

func UpdateAssetClosest ( Asset, PotentialNode Location, PotentialBetterDistance float64 ) ( Location, float64, Location, float64 ) {

	Closest := Assets[ Asset ].NodeDistance
	ClosestNode := Assets[ Asset ].ClosestNode

	nextClosestNode := Assets[ Asset ].NextClosestNode
	nextClosest := Assets[ Asset ].SecondDistance

	SecondDistance := nextClosest
	NodeDistance := Closest

	fmt.Println ( "---------------" )
	fmt.Println ( "Checking:", Asset, "Closest Node:", Closest, "Next:", nextClosest )
	fmt.Println ( "Is", PotentialNode, "(", PotentialBetterDistance , ") closer?", )


//
//	Handle Initial State ( 0, 0 )

	if Closest == 0.0 {
	
		fmt.Println ( "---------------" )
		fmt.Println ( "New Asset, Set to Current Node for Both" )
		fmt.Println ( "---------------" )
		return PotentialNode, PotentialBetterDistance, PotentialNode, PotentialBetterDistance 

	}

//	Compare and update whan Lower

	if PotentialBetterDistance < Closest { 

		fmt.Println ( PotentialNode, "is Closest" )

		SecondDistance = Closest
		nextClosestNode = ClosestNode

		NodeDistance = PotentialBetterDistance
		ClosestNode = PotentialNode

	} else if PotentialBetterDistance < nextClosest { 

		fmt.Println ( PotentialNode, "is Next Closest" )

		SecondDistance = PotentialBetterDistance 
		nextClosestNode = PotentialNode
	}
	fmt.Println ( "---------------" )
	fmt.Println ( "Results:" )
	fmt.Println ( "Closest Node:", ClosestNode, "Distance:", NodeDistance )
	fmt.Println ( "Next Closest:", nextClosestNode, "Next Distance:", SecondDistance )

	fmt.Println ( "---------------" )

	return ClosestNode, NodeDistance, nextClosestNode, SecondDistance 
}

