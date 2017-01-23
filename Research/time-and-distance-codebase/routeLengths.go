package time_and_distance

import (
	"fmt"
)
//	Start --> End Times fore Each Asset

var	routeTimes	map [ Location ] float64

func	computeRouteTimes ( ) {

	routeTimes = make ( map [ Location ] float64 )

	index := NodeAsset{}

	currentNode := LeftEndNode
	nextNode := Nodes [ LeftEndNode ].RightNode

	for {

	    index.Node = currentNode

	    for asset, _ := range Assets {

		index.Asset = asset

		timeToNode := NodeAssets [ index ].NodeToRightNodeTime

		routeTimes[ asset ] = routeTimes[ asset ] + timeToNode
	    }
//	    Break @ Last Node

	    if currentNode == RightEndNode { break }

	    currentNode = nextNode
	    nextNode = Nodes [ currentNode ].RightNode
	}

	fmt.Println ( LeftEndNode, "-->",  RightEndNode, "Times:", routeTimes )
}

