package conversions

import (
	"fmt"
)



/*	Convert Decimal Latitude / Longitude to 
	Degrees, Minutes, Secods, Direction ( DMS )
*/
	
func getDMS ( lat, long float64 ) ( DMS, DMS ) {

	lat_direction := "N"		// Default to North Latitude
	long_direction := "E"		// Default to East Longitude
	
//	Negative Lat and / or Long reverses Direction

	if lat < 0 { 	
	
		lat_direction = "S"
		lat = 0.0 - lat
	}
	if long < 0 { 
	
		long_direction = "W"
		long = 0.0 - long
	}
	
	lat_degrees := int( lat )
	long_degrees := int( long )

//	Minutes are based on the Degrees remainder
	
	lat_minutes_float :=  ( lat - float64 ( lat_degrees ) ) * 60.0
	lat_minutes := int ( lat_minutes_float )

	long_minutes_float :=  ( long - float64 ( long_degrees ) ) * 60.0
	long_minutes := int ( long_minutes_float )

//	Seconds are based on the Minutes remainder
	
	lat_seconds := ( lat_minutes_float - float64 ( lat_minutes ) ) * 60.0
	long_seconds := ( long_minutes_float - float64 ( long_minutes ) ) * 60.0

	dms := DMS{}
	
	dms.Degrees = float64 ( lat_degrees )
	dms.Minutes = float64 ( lat_minutes )
	dms.Seconds = lat_seconds
	dms.Direction = lat_direction
	
	latDMS := dms

	dms.Degrees = float64 ( long_degrees )
	dms.Minutes = float64 ( long_minutes )
	dms.Seconds = long_seconds
	dms.Direction = long_direction
	
	longDMS := dms

	fmt.Println ( "Decimal Convert:", lat, "-->", latDMS )
	fmt.Println ( "Decimal Convert:", long, "-->", longDMS )

	return latDMS, longDMS
}

