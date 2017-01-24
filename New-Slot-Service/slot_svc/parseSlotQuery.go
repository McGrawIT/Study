package slot_svc

import (
	"fmt"
//	"time"
//	"strconv"
//	"net/http"

//	au "github.build.ge.com/AviationRecovery/go-oauth.git"
)

/*
	Primary Search Function, called by FullSearch()

	Query String Parameters / Values

	airportCode=DXB,LAX,MCO
	carrierCode=FZ,DL,B6
	flightNumber=105,7000,14,6,88
	operationDate=2016-May-14,2015-Oct-31

	Added Filtering ( on Slot Index )

	slotIndex=Gx,Lx,Ex,Rx:y

	x = Slot Index value (e.g., 1.0, 0.0, 0.75, 0.6, ...)
	G = Grafter than
	L = Less than
	E = Equal to
	R = From x to y
*/

var (
	DB2				string			// Test URL, assigned in Main
	DBwritten		string			// Data Redd ( after Write )
	usingRange 		= false
	dateRange		string
)


func parseQueryString ( queryParams map[string][]string ) ( FilterParams ) {

//	debug := false

/*	Capture each Filter String (0-n values for each)

	Create a blank Search Filter (when that occurs,
	ALL data is selected vs. "filtering in"
*/
	searchFilters := FilterParams{}

//	queryFilters := []string{}

	fmt.Println ( "Handling", len(queryParams), "search filters" )
	fmt.Println( " Using Range:", usingRange )

	for key, value := range queryParams {

	    parms := ""

//	    Extract the one slice element (to "convert" to a string)

	    for _, parms = range value { }
		switch key {

//		Official Query Parameters

		case "airportCode" : 	searchFilters.airportCodes = parms
		case "carrierCode" : 	searchFilters.carrierCodes = parms

/*		In addition to a 1-n Operation Dates ( an Operation Date Filter List ),
 		the Operation Date Query Parameter supports a filter by a Range that will

 		(a) Cross Seasons ( a Range can overlap more than one, whereas the Operation
 		Date is always in one Season ), and
 		(b) Assume that a Slot Week is included ( it does not attempt to cover the
 		edge case of a Date range being within the first or last 6 days of a Season

 		At the point of parameter detection ( this is the "Operation Date" parameter,
 		check for "Range" vs. a date string ( e.g., 2018-Dac-17 )

 		The format will be OperationDate=Range 2017-Mar-12, 2018-Jan-11
  */
		case "operationDate" :

			if parms[0:5] == "Range" {

				usingRange = true
				dateRange = parms[6:]
				fmt.Println("Date Range Parameter:", dateRange )

				searchFilters.operationDates = ""

			} else {

				usingRange = false
				searchFilters.operationDates = parms
			}

		case "flightNumber" : 	searchFilters.flightNumbers = parms

//		Short-Names ( used primarily for testing ( much easier to type )

		case "AC" : 			searchFilters.airportCodes = parms
		case "CC" : 			searchFilters.carrierCodes = parms
		case "OD" : 			searchFilters.operationDates = parms

			if parms[0:5] == "Range" {

				usingRange = true
				dateRange = parms[6:]
				fmt.Println("Date Range Parameter:", dateRange )

				searchFilters.operationDates = ""

			} else {

				usingRange = false
				searchFilters.operationDates = parms
			}

		case "FN" : 			searchFilters.flightNumbers = parms

/*
		Added these Query Parameters in support of the /Cancellation endpoint
		Cancellation() was added to allow real-time cancellations ( would need
		to be called from another service ( disruption detection? )
 */
		case "flightOrigin" : 		searchFilters.flightOrigin = parms
		case "flightDestination" : 	searchFilters.flightDestination = parms
		case "cancelCategory" : 	searchFilters.cancelCategory = parms
		case "cancelReason" : 		searchFilters.cancelReason = parms
		case "cancelNotes" : 		searchFilters.cancelNotes = parms
		case "cancelSource" : 		searchFilters.cancelSource = parms

		case "FD" : 			searchFilters.flightDestination = parms
		case "CT" : 			searchFilters.cancelCategory = parms
		case "CR" : 			searchFilters.cancelReason = parms
		case "CN" : 			searchFilters.cancelNotes = parms
		case "CS" : 			searchFilters.cancelSource = parms

		case "Season" : 		searchFilters.Season = parms
		case "SN" : 			searchFilters.Season = parms
/*
		A Slot Index query parameter was added to allow filtering ( GT, LT, etc. )
		by Slot Indexes ( in addition to existing filter parameters ( above )
 */
		case "slotIndex" : 		searchFilters.slotIndex = parms
		case "SI" : 			searchFilters.slotIndex = parms

		default : fmt.Println ( "Invalid Query Parameter:", key, "Ignoring" )

		}

	}

	if usingRange {			// Log Search Filters ( Date Range was reported elsewhere )

		if len ( queryParams ) > 1 { fmt.Println ( "Also using Slot Filters:", searchFilters ) }

	} else { fmt.Println ( "Search Using Filters:", searchFilters ) }

/*	Filter Function will take SlotData as input and return JSON
	or the SlotData slice for JSON conversion (by makeJSON())
*/

//	SlotDataDB holds all Slot Data for all Slot Data Files
//	The Result Set (z) is a slice of Slot Records (all fields in API)

	fmt.Println ( "Slot Filters:",searchFilters )
	
	reportFilters ( searchFilters )	// Detailed Printing for Log

	return searchFilters
}
	
func reportFilters ( filters FilterParams ) {

	if filters.airportCodes != "" { fmt.Println ( "Airports:",filters.airportCodes ) }
	if filters.carrierCodes != "" { fmt.Println ( "Carriers:",filters.carrierCodes ) }
	if filters.flightNumbers != "" { fmt.Println ( "Flights:",filters.flightNumbers ) }
	if filters.operationDates != "" { fmt.Println ( "Operation Dates:",filters.operationDates ) }
	if filters.slotIndex != "" { fmt.Println ( "Slot Index Conditions:",filters.slotIndex ) }
	if filters.Season != "" { fmt.Println ( "Season Name:",filters.Season ) }
}
