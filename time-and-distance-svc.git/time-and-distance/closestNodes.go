package time_and_distance

import (
	"fmt"
)

func collectClosestNodes () {

	closestNode := AssetClosestNode{}
	closeIndex := NodeAsset{}

	for NextAsset, _ := range Assets {

	    Asset := Assets [ NextAsset ]

	    fmt.Println ( "Asset:", NextAsset )

	    closeIndex.Node = Asset.ClosestNode
	    closeIndex.Asset = NextAsset

//	    Get Node / Asset Times for Cloeest Node

	    asset := NodeAssets [ closeIndex ]
	    node := Nodes [ Asset.ClosestNode ]

	    closestNode = AssetClosestNode{}

	    closestNode.Asset			= NextAsset
	    closestNode.ClosestNode		= Asset.ClosestNode
	    closestNode.Destination		= Asset.Destination

	    closestNode.ToClosestNode		= Asset.NodeDistance
	    closestNode.ToDestination		= Asset.RouteLength
	    closestNode.ToLeftEnd		= node.ToLeftEndNode
	    closestNode.ToRightEnd		= node.ToRightEndNode

	    closestNode.TimeToClosest		= asset.AssetToNodeTime
	    closestNode.TimeToDestination	= asset.TimeToDestination
	    closestNode.TimeToLeftEnd		= asset.NodeToLeftEndTime
	    closestNode.TimeToRightEnd		= asset.NodeToRightEndTime

	    AssetFromClosest [ closeIndex ] 	= closestNode
	}

//	Flatten map[] for inclusion in Overall Route result set

	oneAsset := ClosestNodesByAsset{}

	closestAssetNode := NodeAsset{}

	for key, value := range AssetFromClosest {

	    closestAssetNode.Node = key.Node
	    closestAssetNode.Asset = key.Asset

	    oneAsset.ClosestNode = closestAssetNode
	    oneAsset.Detail = value

	    allAssets = append ( allAssets, oneAsset )

	}
}

