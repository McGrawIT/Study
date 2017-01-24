package conversions

import (
	"fmt"
)

func	GeoConvert ( Location GeoLocation ) ( GeoLocation ) {

	GeoCoordinate := Location.DMS.Latitude

	Decimal := Point{}

	degrees := GeoCoordinate.Degrees
	minutes := GeoCoordinate.Minutes
	seconds := GeoCoordinate.Seconds
	direction := GeoCoordinate.Direction

	X := degrees + ( minutes / 60.0 ) + ( seconds / 3600.0 )

	if direction == "S" { X = X * -1.0 }

	Decimal.X = X

	fmt.Println ( "Geo Convert:", GeoCoordinate, "-->", X )

	GeoCoordinate = Location.DMS.Longitude

	degrees = GeoCoordinate.Degrees
	minutes = GeoCoordinate.Minutes
	seconds = GeoCoordinate.Seconds
	direction = GeoCoordinate.Direction

	Y := degrees + ( minutes / 60.0 ) + ( seconds / 3600.0 )

	if direction == "W" { Y = Y * -1.0 }

	Decimal.Y = Y

	fmt.Println ( "Geo Convert:", GeoCoordinate, "-->", Y )

	Location.Decimal = Decimal

	return Location
}

