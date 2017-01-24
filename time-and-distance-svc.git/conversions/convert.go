package conversions

import (
	"fmt"
//	"os"
	"encoding/json"
	"net/http"

//	au "github.build.ge.com/AviationRecovery/go-oauth.git"
)

func 	Convert ( rw http.ResponseWriter, req *http.Request ) {

	var	Request 	ConversionRequest

	Request = ConversionRequest{}

	decoder := json.NewDecoder ( req.Body )
	err := decoder.Decode ( &Request )
	if err != nil { }

//	Response Body uses Request Body

	Response = ConversionRequest{}
	ResponseSet = []ConversionRequest{}

//	Some of the Request content is returned "as-is"
//	Try and assign all in one fell swoop...

	Response = Request

	fmt.Println ( "------------------------------------------------" )
	fmt.Println ( "In:", Request )
	fmt.Println ( "------------------------------------------------" )

/*	From "Convert From" Units and "Convert To" Units are provided 
	in the Body of the /Convert Request ( e.g., Inches --> Centimeters )
*/
	fromUnits := Request.ConvertFrom
	toUnits := Request.ConvertTo

	input := Request.OriginalValue		// Value to Convert ( Basic )

	Response.ConvertedValue = doConvert ( fromUnits, toUnits, input ) 

	fmt.Println ( "Converted:", input, fromUnits, "to", Response.ConvertedValue, toUnits )

	Response.Speed = doConvert ( fromUnits, toUnits, Request.Speed  )

	fmt.Println ( "Speed:", Request.Speed, fromUnits, "/ hour converted to", Response.Speed, toUnits, "/ hour" )
	Request.Speed = Response.Speed
	fmt.Println ( "------------------------------------------------" )

//	----------------------------------------------------------
/*	In addition to a simple value conversion, there may be a Geo Segment
	( Two Latitude / Longitude pairs, each in Degrees, Minutes, Seconds, 
	and Direction ).  The segmemnt can also be represented by a pair of
	Decimal ( X, Y ) coordinates.  Whichever is present ( Lat / Long or 
	Decimal ) will be used to convert to the other.
	----------------------------------------------------------
*/
	noGeoLocation := GeoLocation{}
	noDMS := GeoPoint{}
	noDecimal := Point{}

/*	The Geo Segment must have a Starting Geo Position.  Ending is not
	required as this my simply be a conversion between DMS and Decimal.
*/
//	If No Starting Geo Segment, do not process any Ending Segment

	if Request.GeoStart != noGeoLocation { 

//	    Do Start Location ( or start of Geo Segment )

	    startDMS := ( Request.GeoStart.DMS != noDMS )
	    startDecimal := ( Request.GeoStart.Decimal != noDecimal )

	    if startDMS && startDecimal {	// Nothing to convert

	    fmt.Println ( "(Start) DMS and Decimal (no conver)" )

	    } else {

		if startDMS {		// Do DMS --> Decimal

		    Response.GeoStart = GeoConvert ( Request.GeoStart )

		    fmt.Println ( "(Start) DMS:", Request.GeoStart.DMS, "-->", Response.GeoStart.Decimal )

		} else {		// Do Decimal --> DMS

		    geoLocation := GeoLocation{}

		    lat := Request.GeoStart.Decimal.X
		    long := Request.GeoStart.Decimal.Y

		    latDMS, longDMS := getDMS ( lat, long )

		    geoLocation.DMS.Latitude = latDMS
		    geoLocation.DMS.Longitude = longDMS

		    geoLocation.Decimal.X = lat
		    geoLocation.Decimal.Y = long

		    fmt.Println ( "(Start) Decimal:", Request.GeoStart.Decimal, "-->", geoLocation.DMS )

		    Response.GeoStart = geoLocation
		}


	    }

	    endExists := ( Request.GeoEnd != noGeoLocation )

/*	    If End Geo Location exists, convert ( if needed ) 
	    Also, a Segment exists ( compute Distance and ETA )
*/
	    if endExists {

		endDMS := ( Request.GeoEnd.DMS != noDMS )
		endDecimal := ( Request.GeoEnd.Decimal != noDecimal )

		if endDMS && endDecimal {	// Nothing to convert

		    fmt.Println ( "(End) DMS and Decimal (no covert)" ) 

		} else {


		    if endDMS {		// Do DMS --> Decimal

			Response.GeoEnd = GeoConvert ( Request.GeoEnd )

			fmt.Println ( "(End) DMS:", Request.GeoEnd.DMS, "-->", Response.GeoEnd.Decimal )

		    } else {		// Do Decimal --> DMS

			geoLocation := GeoLocation{}

			lat := Request.GeoEnd.Decimal.X
			long := Request.GeoEnd.Decimal.Y

			latDMS, longDMS := getDMS ( lat, long )

			geoLocation.DMS.Latitude = latDMS
			geoLocation.DMS.Longitude = longDMS

			geoLocation.Decimal.X = lat
			geoLocation.Decimal.Y = long

			fmt.Println ( "(End) Decimal:", Request.GeoEnd.Decimal, "-->", geoLocation.DMS )

			Response.GeoEnd = geoLocation
		    }
		}

	fmt.Println ( "------------------------------------------------" )

		startX := Response.GeoStart.Decimal.X
		startY := Response.GeoStart.Decimal.Y

		endX := Response.GeoEnd.Decimal.X
		endY := Response.GeoEnd.Decimal.Y

		gcd := DriveGreatCircle ( startX, startY, endX, endY )
		gcdUnits := "kilometers"

		Response.GreatCircleDistance = gcd

		converted := doConvert ( gcdUnits, toUnits, gcd )

		fmt.Println ( "Converted into:", converted, toUnits )

		ETA := converted / Request.Speed
		Response.GeoSegmentETA = ETA
		fmt.Println ( "ETA:", ETA, "hours" )
	    }
	}


//	----------------------------------------------------------
	fmt.Println ( "------------------------------------------------" )

	noLine := Point{}

//	Only calculate / convert if there is an ( X,Y ) Line Segment

	if Request.LineStart != noLine {

	    startDecimal := Request.LineStart
	    endDecimal := Request.LineEnd

	    fmt.Println ( "Start:", startDecimal, "End:", endDecimal )

	    Response.LineDistance = startDecimal.Distance ( endDecimal )

	    converted := doConvert ( fromUnits, toUnits, Response.LineDistance )

	    fmt.Println ( "Line Distance:", Response.LineDistance, fromUnits )
	    fmt.Println ( "Converted into:", converted, toUnits )

	    Response.LineSegmentETA = converted / Request.Speed
	    fmt.Println ( "ETA:", Response.LineSegmentETA, "hours" )
	    fmt.Println ( "------------------------------------------------" )
	}

/*	Conversion supports a "Convert Distance Only" field.  It exists 
	( in addition to the "Original Value" ) to support a field that gets
	used for an ETA.  The converted Distance overwrites the original
	distance and the ETA field is computed / populated.
*/

//	*** Not implemented ***

//	----------------------------------------------------------------
//	----------------------------------------------------------------

	fmt.Println ( "------------------------------------------------" )
	fmt.Println ( "Out:", Response )
	fmt.Println ( "------------------------------------------------" )

//	Load JSON Response
	
	ResponseSet = append ( ResponseSet, Response )

	result, err := json.MarshalIndent ( ResponseSet , "", "  " )

	if err != nil { fmt.Println ( "JSON Marshall Failed:", err ) }

	fmt.Fprintf ( rw, string ( result ) ) 
}

/*	Use Great Circle calculation for Geo Segment distnaces.  It takes a
	Decimal Lat / Long, so when we only have Degrees, Minutes, Seconds
	they are converted to Decimal Lat / Long.  Response always contains
	both Deciaml and DMS values.  If either is missing, the other is
	used to perform the conversion.
*/

func DriveGreatCircle ( X1, Y1, X2, Y2 float64 ) ( float64 ) {

	Origin := gcLocation{}
	Destination := gcLocation{}

	Origin.Name = "Austin, TX"
	Destination.Name = "West Palm Beach, FL"

	Origin.Lat = X1
	Origin.Long = Y1

	Destination.Lat = X2
	Destination.Long = Y2
	
//	Use haversine to get the resulting diatance between the two values

	d := GreatCircle ( Origin, Destination )

	fmt.Println ( "Distance:", d, "kilometers" )

	return d
}
