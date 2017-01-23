package time_and_distance

import (
	"fmt"
)

var 	routeLeftNode		Location
var	routeRightNode 		Location

func	getRouteEnds ( ) {

	fmt.Println ( "Get Start / End Nodes" )

	findEnd := EdgeNodesLL [ routeRightNode ].RightNode
	findStart := EdgeNodesLL [ routeLeftNode ].LeftNode

	routeEnd = routeRightNode
	routeStart = routeLeftNode

	endNode := Location{}

	for { if findEnd == endNode { break }

	    routeEnd = findEnd
	    findEnd = EdgeNodesLL [ findEnd ].RightNode
	}

	for { if findStart == endNode { break }

	    routeStart = findStart
	    findStart = EdgeNodesLL [ findStart ].LeftNode
	}
}


