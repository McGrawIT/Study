package conversions

import (
	"math"
)

/*************************************************************************
*************************************************************************/

type	ConversionRequest	struct	{

	GeoStart		GeoLocation	`json:"GeoStartLocation"`
	GeoEnd			GeoLocation	`json:"GeoEndLocation"`
	GreatCircleDistance 	float64	 	`json:"GeoSegmentDistance"`
	GeoSegmentETA	 	float64	 	`json:"GeoSegmentETA"`

//	Straight Line Distance will be returned in the "Convert To" Units
//	This "strainght line" Distance is also a Time & Distance capability

	LineStart		Point		`json:"LineStart"`
	LineEnd			Point		`json:"LineEnd"`
	LineDistance 		float64	 	`json:"LineDistance"`
	LineSegmentETA	 	float64	 	`json:"LineSegmentETA"`

	OriginalValue		float64		`json:"OriginalValue"`
	ConvertFrom		string		`json:"ConvertFrom"`
	ConvertTo 		string	 	`json:"ConvertTo"`
	ConvertedValue 		float64	 	`json:"ConvertedValue"`

	ConvertDistanceOnly 	float64	 	`json:"ConvertDistanceOnly"`
	Speed	 		float64	 	`json:"Speed"`

}

/*************************************************************************
*************************************************************************/

type	GeoLocation	struct {

	DMS		GeoPoint 	`json:"DMS"`
	Decimal		Point 		`json:"Decimal"`
}

type	GeoPoint	struct {
	
	Latitude 	DMS	 	`json:"Latitude"`
	Longitude 	DMS	 	`json:"Longitude"`
}

type	DMS		struct {

	Degrees		float64		`json:"Degrees"`
	Minutes		float64		`json:"Minutes"`
	Seconds		float64		`json:"Seconds"`
	Direction	string		`json:"Direction"`
}

type	Point		struct {

	X		float64		`json:"X"`
	Y		float64		`json:"Y"`
}


func ( p Point ) Distance( p2 Point ) float64 {

	first := math.Pow( float64( p2.X - p.X ), 2 )
	second := math.Pow( float64( p2.Y - p.Y ), 2 )

	return math.Sqrt( first + second )
}





//	--------------------------------------------------------------



var	ResponseSet 	[]ConversionRequest
var	Response 	ConversionRequest

	
var	Systems		MeasurementSystems

var	Factor		= map [ MeasurementSystems ] float64 {

	MeasurementSystems{ "feet", "meters" } : 0.3048,
	MeasurementSystems{ "meters", "feet" } : ( 1 / 0.3048 ),
	MeasurementSystems{ "meters", "meters" } : 1.0,
	MeasurementSystems{ "feet", "feet" } : 1.0,
}

type	MeasurementSystems	struct {

	FromSystem	string	
	ToSystem	string
}


type	unitRatio	struct {

	Ratio	float64
	Unit	string
}

var	fromUnits	string
var	toUnits		string

/*************************************************************************
			-----------------
			CONVERSION TABLES
			-----------------
*************************************************************************/



var	fromRatio 	= map [ string ] unitRatio { 

//	Imperial / US

	"inch" : unitRatio{ ( 1.0 / 12.0 ), "feet" },
	"feet" : unitRatio{ 1.0, "feet" }, 
	"yard" : unitRatio{ 3.0, "feet" }, 
	"mile" : unitRatio{ 5280.0, "feet" }, 

//	Imperial / US (	Plural )

	"inches" : unitRatio{ ( 1.0 / 12.0 ), "feet" },
	"yards" : unitRatio{ 3.0, "feet" }, 
	"miles" : unitRatio{ 5280.0, "feet" }, 

//	Imperial / US (	Short From )

	"in" : unitRatio{ ( 1.0 / 12.0 ), "feet" },
	"mi" : unitRatio{ 5280.0, "feet" }, 

//	Imperial Only

	"chain" : unitRatio{ 66.0, "feet" }, 
	"furlong" : unitRatio{ 660.0, "feet" }, 
	"league" : unitRatio{ 15840.0, "feet" }, 
	"link" : unitRatio{ ( 12.0 * 7.92 ), "feet" }, 
	"rod" : unitRatio{ ( 25.0 * ( 12.0 * 7.92 ) ), "feet" }, 
	"imperial cable" : unitRatio{ 100.0 , "fathoms" }, 


//	Metric

	"kilometer" : unitRatio{ 1000.0, "meters" }, 
	"meters" : unitRatio{ 1.0, "meters" }, 
	"dectimeter" : unitRatio{ 0.10, "meters" }, 
	"centimeter" : unitRatio{ 0.01, "meters" }, 
	"milimeter" : unitRatio{ 0.001, "meters" }, 

//	Metric ( Plural )

	"kilometers" : unitRatio{ 1000.0, "meters" }, 
	"centimeters" : unitRatio{ 0.01, "meters" }, 
	"milimeters" : unitRatio{ 0.001, "meters" }, 

//	Metric ( Short From )

	"cm" : unitRatio{ 0.01, "meters" }, 
	"mm" : unitRatio{ 0.001, "meters" }, 
	"km" : unitRatio{ 1000.0, "meters" }, 

//	Nautical

//	"nautical mile" : unitRatio{ ( 5280.0 * 1.15078 ), "feet" }, 
	"nautical mile" : unitRatio{ 1852.0, "meters" }, 

	"fathom" : unitRatio{ 6.08, "feet" }, 
	"cable" : unitRatio{ ( 6.08 * 120.0 ) , "feet" }, 
	"shot" : unitRatio{ ( 6.08 * 15.0 ), "feet" }, 
	"metric cable" : unitRatio{ 200.0, "meters" }, 
}

//	--------------------------------------------------------------

var	toRatio		= map [ string ] unitRatio {

	"inch" : unitRatio{ 12.0, "feet" },	
	"feet" : unitRatio{ ( 1.0 ), "feet" }, 
	"yard" : unitRatio{ ( 1.0 / 3.0 ), "feet" }, 
	"mile" : unitRatio{ ( 1.0 / 5280.0 ), "feet" }, 

	"inches" : unitRatio{ 12.0, "feet" },	
	"yards" : unitRatio{ ( 1.0 / 3.0 ), "feet" }, 
	"miles" : unitRatio{ ( 1.0 / 5280.0 ), "feet" }, 

	"in" : unitRatio{ 12.0, "feet" }, 
	"mi" : unitRatio{ ( 1.0 / 5280.0 ), "feet" }, 

	"bizarre" : unitRatio{ ( 1.0 / 1000.0 ), "meters" }, 
	"kilometer" : unitRatio{ ( 0.001 ), "meters" }, 
	"meters" : unitRatio{ 1.0, "meters" }, 
	"centimeter" : unitRatio{ 100.0, "meters" },
	"millimeter" : unitRatio{ 1000.0, "meters" }, 

	"bizarres" : unitRatio{ ( 1.0 / 1000.0 ), "meters" }, 
	"kilometers" : unitRatio{ 0.001, "meters" },
	"centimeters" : unitRatio{ 100.0, "meters" },
	"millimeters" : unitRatio{ 1000.0, "meters" }, 

	"bm" : unitRatio{ ( 1.0 / 1000.0 ), "meters" }, 
	"km" : unitRatio{ 0.001, "meters" }, 
	"cm" : unitRatio{ 100.0, "meters" }, 
	"mm" : unitRatio{ 1000.0, "meters" }, 

//	Nautical

	"nautical mile" : unitRatio{ ( 1.0 / 1852.0 ), "meters" }, 

	"fathom" : unitRatio{ ( 1.0 / 6.0 ), "feet" }, 
	"cable" : unitRatio{ ( 1.0 / 120.0  ), "fathoms" }, 
	"shot" : unitRatio{ ( 1.0 / 15.0  ), "fathoms" }, 
	"metric cable" : unitRatio{ ( 1.0 / 200.0 ) , "meters" }, 

	"fathoms" : unitRatio{ ( 1.0 / 6.0 ), "feet" }, 
	"cables" : unitRatio{ ( 1.0 / 120.0  ), "fathoms" }, 
	"shots" : unitRatio{ ( 1.0 / 15.0  ), "fathoms" }, 
	"metric cables" : unitRatio{ ( 1.0 / 200.0 ) , "meters" }, 
}



