package time_and_distance

import (
	"fmt"
)

func LoadAssets () {

	AssetLocation := Location{}
	asset := AssetInstance{}

	for _, Input := range InputAssets {

		AssetLocation = Input.AssetLocation
	
		asset.AssetLocation = Input.AssetLocation
		asset.Destination = Input.Destination
		asset.Speed = Input.AssetSpeed

		asset.AssetID = Input.AssetID
		asset.AssetType = Input.AssetType
		asset.Weight = Input.Weight
		asset.Height = Input.Height
		asset.Width = Input.Width
		asset.AwayFromNodeMax = Input.AwayFromNodeMax
		asset.ValidZones = Input.ValidZones


		Assets [ AssetLocation ] = asset

		fmt.Println ( "Asset from JSON:", Input )
		fmt.Println ( " Asset in Route:", asset )
		if asset.Speed == 0.0 { fmt.Println ( "No Speed" ) }
		fmt.Println ( "---------------" )

		asset = AssetInstance{}
	}
}


func	LoadPostBodyAssets ( routeParameter Route ) {

	for _, Input := range routeParameter.AssetsInRoute {

		fmt.Println ( "Asset in Route:", Input )

		asset := Input

		InputAssets = append ( InputAssets, asset )
	}

	fmt.Println ( "Loaded", len ( InputAssets ), "Assets" )
}

