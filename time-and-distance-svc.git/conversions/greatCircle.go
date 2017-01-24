package conversions

import (
	"fmt"
	"math"
)

const 	kmtomiles = float64(0.621371192)
const 	earthRadius = float64(6371)

type	gcLocation struct {

	Name	string
	Lat	float64
	Long	float64
}

func GreatCircle ( O, D gcLocation ) ( d float64 ) {

	lonO := O.Long
	latO := O.Lat
	lonD := D.Long
	latD := D.Lat

	d = Haversine ( lonO, latO, lonD, latD )

	inMiles := d * kmtomiles

	fmt.Printf ( "%s to %s: %.02f miles\n", O.Name, D.Name, inMiles )

	return
}


/*	The haversine formula will calculate the spherical distance as the 
	crow flies between Lat and Long for two given points in kilometers
*/
func Haversine ( lonFrom, latFrom, lonTo, latTo float64 ) ( distance float64 ) {

	pi180 := ( math.Pi / 180 )
	
	var 	deltaLat = ( latTo - latFrom ) * pi180
	var 	deltaLon = ( lonTo - lonFrom ) * pi180
	
	var a = math.Sin ( deltaLat / 2 ) * math.Sin ( deltaLat / 2 ) + 
		math.Cos ( latFrom * pi180 ) * math.Cos ( latTo * pi180 ) *
		math.Sin ( deltaLon / 2 ) * math.Sin ( deltaLon / 2 )

	var c = 2 * math.Atan2 ( math.Sqrt(a), math.Sqrt(1-a) )
	
	distance = earthRadius * c
	
	return
}

