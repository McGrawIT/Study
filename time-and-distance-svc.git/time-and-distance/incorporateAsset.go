package time_and_distance

import (
	"fmt"
)

func incorporateAsset ( previousNode, currentNode, nextNode, NextAsset Location, newNode NodeInstance ) {

	fmt.Println ( "Asset:", NextAsset )
	fmt.Println ( "---------------" )

	CurrentNodeAssetIndex := NodeAsset{}

	CurrentNodeAssetIndex.Node = currentNode
	CurrentNodeAssetIndex.Asset = NextAsset

	CurrentNodeAssetIndex.Asset = NextAsset

	Asset := Assets [ NextAsset ]

	nodeAssetDetails := RouteAssetDetail{}

	nodeAssetDetails.Asset = Asset.AssetLocation 
	nodeAssetDetails.Destination = Asset.Destination
	nodeAssetDetails.Speed = Asset.Speed

//	Result Set ( for Response Body )

//	assetResult.AssetLocation = Asset.AssetLocation 
//	assetResult.AssetDestination = Asset.AssetDestination
//	assetResult.Speed = Asset.Speed

	Speed := Asset.Speed

//	Get Distance to Node ( from Asset )
//	Get Time to Node using Asset Speed

	AssetToNode := Node_AssetLength ( currentNode, NextAsset )

	AssetToNodeTime := 0.0
	if Speed != 0.0 { AssetToNodeTime = AssetToNode / Speed }

	nodeAssetDetails.AssetToNode = AssetToNode
	nodeAssetDetails.AssetToNodeTime = AssetToNodeTime

/*	-------------------------------------------------------------
	This will need to be calculated once the Route is complete
	Walk the Node / Asset records.  Use a Route Segment calcuation
	that leveerages Current Node and End Nodes ( it exists ) to
	get the Asset --> Destination distance and time for all Nodes
	-------------------------------------------------------------
*/
	nodeAssetDetails.TimeToDestination = 0.0

/*	Time to Previous and Next Nodes ( from this Node ) 
	are tied to Asset-by-Asset Speed

	Compute for each Asset in Time and Distance request.  
*/
//	Time for this Asset to go from Current to Left / Right Nodes


	CurrentToLeftTime := 0.0
	CurrentToRightTime := 0.0

	if Speed != 0.0 {
	    CurrentToLeftTime = newNode.ToLeftNode / Speed
	    CurrentToRightTime = newNode.ToRightNode / Speed
	}


	nodeAssetDetails.NodeToLeftNodeTime = CurrentToLeftTime
	nodeAssetDetails.NodeToRightNodeTime = CurrentToRightTime


fmt.Printf ( "%v-->%v: %.2f (%.2f)\n", NextAsset, currentNode, AssetToNode, nodeAssetDetails.AssetToNodeTime ) 

fmt.Println ( previousNode, "<---", CurrentToLeftTime, "---", currentNode, "---", CurrentToRightTime, "--->", nextNode )

	CurrentNodeToLeftEndTime := 0.0
	CurrentNodeToRightEndTime := 0.0

	if Speed != 0.0 {

	    CurrentNodeToLeftEndTime = newNode.ToLeftEndNode / Speed
	    CurrentNodeToRightEndTime = newNode.ToRightEndNode / Speed
	}

	fmt.Println ( "Speed[", Speed, "]" )
	fmt.Println ( "-----------------" )
		
	fmt.Println ( "Times[", "Node to Left End:", CurrentNodeToLeftEndTime , "]" )
	fmt.Println ( "Times[", "Node to Right End:", CurrentNodeToRightEndTime , "]" )
	fmt.Println ( "Times[", "Node to Left:", CurrentToLeftTime , "]" )
	fmt.Println ( "Times[",  "Node to Right:", CurrentToRightTime , "]" )

	nodeAssetDetails.NodeToLeftEndTime = CurrentNodeToLeftEndTime
	nodeAssetDetails.NodeToRightEndTime = CurrentNodeToRightEndTime

	nodeAssetDetails.AssetToLeftEndTime = CurrentNodeToLeftEndTime + AssetToNodeTime
	nodeAssetDetails.AssetToRightEndTime = CurrentNodeToRightEndTime + AssetToNodeTime
	NodeAssets [ CurrentNodeAssetIndex ] = nodeAssetDetails

	fmt.Println ( "Node / Asset", CurrentNodeAssetIndex, "Details:", nodeAssetDetails )
	fmt.Println ( "-----------------------------------" )

/*
	Maintian Closest and Next Closest Nodes for Asset
	This will likely be based on Time, not Distance
	Can we maintain both?
*/
	Closest, ClosestDistance, NextClosest, NextClosestDistance := UpdateAssetClosest ( NextAsset, currentNode, AssetToNode) 

	Asset.ClosestNode = Closest
	Asset.NodeDistance = ClosestDistance

	Asset.NextClosestNode = NextClosest
	Asset.SecondDistance = NextClosestDistance

	Assets [ NextAsset ] = Asset


	cn := Assets [ NextAsset ].ClosestNode
	cd := Assets [ NextAsset ].NodeDistance
	nn := Assets [ NextAsset ].NextClosestNode
	nd := Assets [ NextAsset ].SecondDistance

	fmt.Println ( "---------------" )
	fmt.Println ( "Results in Asset", NextAsset, "Record:" )
	fmt.Println ( "Closest Node:", cn, "Distance:", cd )
	fmt.Println ( "Next Closest:", nn, "Next Distance:", nd )

	fmt.Println ( "===============================================" )

	return
}

