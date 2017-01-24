package time_and_distance

import (

	"fmt"
)

func	createSegmentList ( routeParameter Route ) {

//	Create Doublely-Linked-List of Route Nodes / Edges ( Node-Node )
//	Node / Edge-specific information is supplied through Left Node

	EdgeNodesLL = make ( map [ Location ]RouteEdge )

//	InputSegments = []Edge{}

	emptyEdge := RouteEdge{}
	noLocation := Location{}

	routeLeftNode = Location{}
	routeRightNode = Location{}

	for _, routeEdge := range routeParameter.Segments {

//	    Assume Bi-Directional for first implemantion

//	    Read Both EdgeNodes ( for Segment Pair )
		
	    routeLeftNode = routeEdge.LeftNode
	    routeRightNode = routeEdge.RightNode

	    LeftEdge := EdgeNodesLL [ routeLeftNode ]
	    RightEdge := EdgeNodesLL [ routeRightNode ]

/*	    Each Node is only created once ( when it shows up as the
	    Left Node ( L, R ) or the Right Node ( L, R ) )
*/
	    if LeftEdge == emptyEdge {

//		Left Node is new, so add Route Edge ( it will have
//		all of the Ndoe-specific information ).  It will only
//		be missing its Left Node ( that will be added when
//		this Node appears on the Right Side of the Segment )

//		The Left Node may have already appeared as the Right
//		Node, which set the Left Node ( either way, get it )

		leftNode := EdgeNodesLL [ routeEdge.RightNode ].LeftNode

		if leftNode != noLocation { routeEdge.LeftNode = leftNode }

		routeEdge.LeftNode = Location{}
		EdgeNodesLL [ routeLeftNode ] = routeEdge

	    } else {

//		Left Node ( L, R ) of Segment was already created
//		when encountered as the Right Node of another Segment
//		Only the Left Node was captured ( see above ) and it
//		needs to be aded to this full Node ( that's all that
//		is missing )

		leftNode := EdgeNodesLL [ routeEdge.LeftNode ].LeftNode
		routeEdge.LeftNode = leftNode

		EdgeNodesLL [ routeLeftNode ] = routeEdge
	    }
/*
*/

	    if RightEdge == emptyEdge {

//		Right Node ( L, R ) of Segment is a new Node, just
//		store the Left Node and complete when R shows up
//		as a Left Node

		newRightEdge := RouteEdge{}
		newRightEdge.LeftNode = routeLeftNode

		EdgeNodesLL [ routeRightNode ] = newRightEdge 

	    } else {

/*		Right node of Segement exists ( it already showed up
		as the Left Node ); its Left Node is in ( L, R )
*/
		leftNode := EdgeNodesLL [ routeLeftNode ]
		leftNode.LeftNode = routeEdge.LeftNode

		EdgeNodesLL [ routeRightNode ] = leftNode
	    }
	}

	fmt.Println ( "Sengement List Created:", len ( EdgeNodesLL ) )
}



