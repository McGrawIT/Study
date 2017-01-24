package slot_svc

import (
	"fmt"
	"time"
	"strings"
)


/*	--------------
	Slot Filtering
	--------------
	Handle the Core function of Slot Service:  Respond to Queries for Slot Data
	Apart from Loading Slot Data Files from Config Service, this is the primary
	function of Slot Service.
	--------------------------------------------------------------------
	The API as specified for this service:

	/SlotData?AirportCode=A1,A2,A3&CarrierCode=C1,C2,C3&FlightNumber=F1,F2&OperationDate=O1,O2

		A = 3-character Airport Code ( e.g., MCO, ATL, LAX, ... FXB ( Dubai
		C = Up to 3-character Carrier Code ( e.g., DAL, B6, ... FZ ( Fly Dubai )
		F = Flight Number
		O = Operation Date ( YYYY-MMM-DD ) ( e.g., 2018-Mar-22 )

		Other Operation Date formats are possible; they are normalized on input and will
		handle any combination of the following to convert to the YYYY-MMM-DD format ):

			2-digit Years ( 18 vs. 2018 )
			Numeric Months ( 10 becomes Oct )
			One-digit values ( for either Month or Day; expanded to two-digits )

			The final conversion supports a lookup for Month string and the pattern given
			to the time.Parse() function ( for string to time format )

	Any combination of these four Parameters is valid, including none

	Parameters themselves are connected with AND logic.  ( e.g., AirportCode AND CarrierCode )

	Query Parameters are not required, and any that are not supplied are treated
	as include "ALL" of that Slot Attribute.  No Parameters is equivalent to a SELECT ALL

	All Columns are always returned ( SELECT * FROM SLOT_DATA WHERE ... )
	------------------------------------------------------------------------------------------

	Additional Features / Capabilities
	----------------------------------

	(1) Multiple Parameter Values

	Expanded support for multiple Query Parameter Values for each Query Parameter

	Comma-delimited strings of 0-n values for each Parameter that are used
	with OR logic.  ( Airport LAX or MCO, Airline DL or FZ, ...)

	The combination of "any of these" for each parameter offers greater
	search options.  (Flight Number, not so much.)

	Example Endpoint and Query Parameters / Values:

		/SlotData?AirportCode={...}&CarrierCode={...}&FlightNumber={...}&OperationDate={...}

		AirportCode-MCO,SFO,BOS,ATL,DXB
		CarrierCode=DL,B6,FZ
		FlightNumber=79,3280
		OperationDate=2016-Apr-12,2008-Dec-25

		( No limit on values )

		OperationDate=Range 2015-May-12, 2000-Dec-12

	(2) Slot Index Query Parameter

	Slot Data Queries can also be filtered by Slot Index.  It is in addition to existing
	Query Parameters, so it is added to the AND logic.

	The added Parameter is "SlotIndex"

	Slot Index Query Parameters are a string of Relationships ( e.g., Slot Indexes > 0.5 )
	Commas allow for more than one, and are used as OR filters

	Supported Filtering:  	GT ( Greater Than ),
							LT ( Less Than ),
							EQ ( Equal To ),
							NE ( Not Equal ), and

							RANGE ( Slot Index within specified Range of Slot Indexes )

	Example use of Query Parameter:

	/SlotData?SlotIndex=GT 0.5,  LT 0.3,  EQ 0.75,  RANGE 0.0:0.25, NE 0.0

	This is example would return all  Slots except those with Indexes between 0.3 and 0.5
	because ( like all other Parameter Filter Lists ) the values are connected with OR logic

	(3) Slot Cancellation Calculations, and the effect on a Slot

	Cancellation information is saved to maintain current Slot Indexes.
	Cancellation History per Slot ( Season, Airport, Carrier, Operation Day of Week, and Flight )
	is maintained and persisted in the Slot File DB ( on Blob Store )
*/

type FilterParams struct {

	airportCodes		string		// Include these ( 1 - n ) Airport Codes in Result Set
	carrierCodes		string		// Same, except for Carrier Codes
	flightNumbers		string		// and Flight Numbers

/*	OpData filters support two types of Queries: 

	YYYY-MMM-DD						// Include the ( 1 - n ) Specific Operation Dates
	YYYY-MMM-DD - YYYY-MMM-DD		// Include an Operation Date Range ( can span Seasons )
*/
	operationDates		string

//	Cancellations

	flightOrigin		string		// Origin Airport ( Code )
	flightDestination	string		// Destination Airport ( Code )

	cancelCategory		string		// Weather, Mechanical, ...
	cancelReason		string		// Detail ( Rain, Engine )
	cancelNotes			string		// Free Text @ Source
	cancelSource		string		// Who / What Reported ?

	Season				string		// Season Name

//	Allow Result Sets to be created / affected by the Slot Index

	slotIndex			string		// Slot Index Filter ( GT, LT, EQ, NE, Range )
}

var (

	rangeStart		string			// Beginning of Date Range Filter
	rangeEnd		string			// End of Filter
)

/*	Function to Filter Slots
	------------------------

	Inputs
	------
	slotDataFileDB		in-memory Slot File DB ( SlotDataDB ) map[ SlotKey ][]FlightSlot
						SlotDataDB holds all Flight Slots for all Slot Files

	filters				query parameters as strings ( values ) per query parameter
 */

/* 	----------------
	Results Returned
	----------------
	slots				list of Slots ( []slots )

	Sample Result Set Slot ( one Slot ( Day ) of a Flight Slow Week ) in Result Set
  {
    "carrierCode": "DL",
    "airportCode": "MCO",
    "SeasonStart": "2018-Apr-15",
    "SeasonEnd": "2026-Dec-24",
    "weekday": "Sunday",
    "OpDate": "",
    "flightNumber": "20",
    "flightsInSeason": 20,
    "slotIndex": 0,
    "maxCancels": 4,
    "originalMaxCancels": 4
  },
 */
/*	Filter Function will take SlotData as input and return the filtered Slot(s)
*/
func filterSlots ( slotDataFileDB map[SlotKey][]FlightSlot, filters FilterParams ) ( []slots ) {

	startTime := time.Now()
	fmt.Println ( "Start QUery:", startTime )

/*	Slot Data Files are accessed by their Key ( Airport Code, Carrier Code,
	and Season (a Date Range) ).  Before searching all Slot Records within
	a Slot Data File, Check the Key ( the "Value" side will be the entire
	Slot Data File) and skip unless there is a match.

	Basically, skip Slot Data File if any Key Element does not match Criteria

	There will be One Slot Record (for a Week) for Each Flight in each Slot Data File
*/
	slotResultSet := []slots{}	// Start with an Empty Result Set
	slotDataFiles := 0 			//	Count Slot Data File ( each for a Season ) in the Slot DB
	currentSlot := slots{}		// Current Row ( built along the way, discarded if filters do not pass )

/*	More than one value for each Criteria is supported.  Split out 0-n
	values for each to use in a OR selection

	The first three Query Parameters identify the specific Slot Data File
	that holds the Slot Week for each Flight in the Slot Data File

	Third part of Key is the Season ( a Date Range ) and this Query's Parameter values
	are Operation Dates;  Operation Dates "filter" the Slot Data File's third Key ( Season )
	determines which Season Key is valid ( Operation Date IN Season )
*/
	includedAirports := []string{}		// One of three parts of the Data File Key
	includedCarriers := []string{}		// Second part of Key
	includedDates := []string{}

//	This filters the Flights inside the Slot Data File ( for their Slot Weeks ( one / Flight )

	includedFlights := []string{}

	if len( filters.airportCodes ) != 0 { includedAirports = strings.Split ( filters.airportCodes, "," ) }
	if len( filters.carrierCodes ) != 0 { includedCarriers = strings.Split ( filters.carrierCodes, "," ) }
	if len( filters.flightNumbers ) != 0 { includedFlights = strings.Split ( filters.flightNumbers, "," ) }

	if len( filters.operationDates) != 0 {
	    includedDates = strings.Split ( filters.operationDates, "," )
	    includedDates = opDatesNormalize ( includedDates )
	}

/*	Slot Index Filters are a string of Conditionals with corresponding values.
	Like airportCode, carrierCode, this is another filter for the /SlotData endpoint.

	Valid Conditionals with sample values: slotIndex=GT 0.6,LT 0.1,EQ 0..5,Range 0.4:0.45
*/
	indexFilter := []indexFilters{}

	if len ( filters.slotIndex ) != 0 { indexFilter = parseSlotIndexQuery ( filters.slotIndex ) }

/*	Check every Slot File in the DB, selecting for filtering those Slot Files with
	File Keys that match all of the 0-n ( Airport Code / Carrier Code / Season )

	slotFileKey:	Key of Slot Files ( Airport / Carrier / Season ) in Slot DB
					Each match is an Instance of a Slot Data File

					Each Slot File is made up of 1-n Flight Slots ( one week / Flight )

	flightSlotWeek:	Instances of Flight Slots [ Flight Number, Slot Week ( Days ) ]
*/
	fmt.Println ("Searching", len(slotDataFileDB), "Slot Files to create Result Set" )

	slotFilesSearched := 0
	slotIndexesFailed := 0
/*
	If a Slot File is NOT for any of the "filter-by" Filter Sets, then get the next Slot File
	in the DB.  ( A Slot File's Key must match all Key Elements ( Airport, Carrier, AND Season )
	in their respective Filter Sets )

    If No Slot Data File Key Filters are provided, every Slot Data File in the DB will be
    included in the Flight Slot Week filtering
 */
	oldSlotFiles := 0

	for slotFileKey, _ := range slotDataFileDB {		// Check all Slot Data Files in the DB

/*		Check for "old" Slot Data Files.  Early testing created "invalid" Seasons
		They have 2-digit years ( 16 ) vs. 4-digit ( 2016 ).  Those files remain
		in production Slot DBs ( for Fly Dubai only )
 */
		if 	len ( slotFileKey.SeasonEnd ) == 9 ||
			len ( slotFileKey.SeasonStart ) == 9 {

			oldSlotFiles++				// Count Slot Data File skipped
			continue
		}
		slotDataFiles++					// Count of Slot Data Files in DB ( for Logging only )

//		-----------------
//		AIRPORT FILTERING
//		-----------------
/*		The Slot DB Key Search Filter must be for at least one Airport, or ALL Airports are
		valid and all Slot Files will pass the Airport portion of the Slot File Key test
		( All Slot File DBs pass )
*/
/*	    If at least one Airport Code is in the Airport Code Filter Set, check to see if the Key for this
		Slot Data File [ Airport, Carrier, Season ] includes an Airport in the Airport Filter Set.

		If there are no Airports in the Filter, then ALL Slot Data Files pass and will be checked
		for Carriers ( uses AND filters )
*/
	    if len ( includedAirports ) != 0 {		// Airport Code Filter(s) exist

			foundOne := false

			for _, thisAirport := range includedAirports {

				if thisAirport == slotFileKey.Airport {

					foundOne = true
					break
				}
			}
			if !foundOne { continue }			// No Airport Code match, skip to next Slot Data File
	    }

/*	    This Slot File is for one of the Airports in the Airport Filter Set

		Check the rest of the Key ( Carrier and Season ) before including the Slot Data File
		in Flight Slot Week Filtering ( down to each Slot ( Day ) in each Flight Slot Week )
*/
//		-----------------
//	    CARRIER FILTERING
//		-----------------
/*	    Like the Airport Filter, this Slot Data File Key must include a Carrier from the
		Carrier Filter Set.  If no Carrier filtering no Slot Data file will be excluded
		based on Carrier.
*/
	    if len ( includedCarriers ) != 0 {		// Carrier Code Filter(s) exist

			foundOne := false

			for _, thisCarrier := range includedCarriers {

				if thisCarrier == slotFileKey.Carrier {

					foundOne = true
					break
				}
			}
			if !foundOne { continue }		// No Carrier Code match, continue to next Slot Data File
	    }

//		Begin to load the Current Slot record to add to the Result Set ( if it passes all Filters )

	    currentSlot.airportCode = slotFileKey.Airport
	    currentSlot.carrierCode = slotFileKey.Carrier
	    currentSlot.beginDate = slotFileKey.SeasonStart
	    currentSlot.endDate = slotFileKey.SeasonEnd

//		------------------------------------------------------
//	    OPERATION DATES, FLIGHTS, FLIGHT SLOT WEEKS, and SLOTS
//		------------------------------------------------------

/* 		Slot Data Files are for a Season;  Season Begin / End Dates are part of the
		Slot Data File Key, so at least one of the Operation Dates in the Operation Date
		Filter Set must be for this Slot Data File's Season
 */
/* 		Slots are by Day, seven / week.  A Flight Slot Week is the Week of Slots for a
		specific Flight.  Flight Slots are defined for an entire Season.  Operation Date tests
		find the Season Slot File that includes the date.  Within the Slot File, the Operation
		Date is a Weekday ( Sunday to Saturday ), and matches that Day's slot for all Flights

	    Do a Date Range "inclusive" test using the Season Start / End Dates as the search value(s)

	    If no Operation Date was specified, all Slot Rows of the current Slot Data File
		( based on other Filters ) will be included vs. the specific Operation Day(s)

	    If more than one Operation Date is specified, the corresponding Slot Day(s) ( Sunday to Saturday )
	    will be in the Result Set ( assuming it passes the Flight and Slot Index filters )

	    The Flight Weeks ( one or more days ) included in the Result Set will have the same Slot Days
*/
	    allDays := false					// Default is Operation Date(s) are specified

	    if len ( includedDates ) == 0 { 	// No Operation Date Filters were provided

			allDays = true

//			Make sure you get into the Date Filter range loop.

			includedDates = append (  includedDates, "NONE" )
	    }

	    currentOpDate := ""

	    foundOne := false
	    firstDateInSeason := true

/*	    For each Operation Date Filter in the Season of this Slot Data File, examine the
		Flight Slot Weeks for matching Slots ( for a Day ) to see if they pass the
		Slot Index Filter
*/
/*	    If no Operation Date was specified, the "NONE" date ensures entry into the loop vs.
		separate Flight Slot Week filtering
*/
		showDateRange := true

		for _, thisDate := range includedDates {		// Check each Operation Date Filter

			begin := slotFileKey.SeasonStart
			end := slotFileKey.SeasonEnd

/*			Ignore any Operation Date outside the current Slot Data File Season.
			Also test for the "NONE" "flag" ( that will exist when no Operation Date was given

			The "using Range" flag is part of the "filter by Date Range" feature
*/
			if thisDate != "NONE" || usingRange {

				if usingRange {

					if showDateRange {

						dateRange := strings.TrimSpace( dateRange )

						dateRangeDates := opDatesNormalize ( strings.Split( dateRange, ",") )

						rangeStart = dateRangeDates[0]
						rangeEnd = dateRangeDates[1]

						showDateRange = false
					}

/*					If the Season is within the Date Range Filter, then this Slot Data File
					is included.  ( Operation Date filters are bypassed to process ALL Slot Data
					Files. )  More than one Slot Data File can be value since the Operation Date
					Range can span Seasons.
*/
/*					To handle each case of Date Range / Season overlap, check both conditions:

					Only one end of the Operation Date Range needs to be in the Season, OR
					Only on end of the Season needs to be within the Operation Date Range

					This supports a Season within a Range as well as a Range within a Season as well
					as when a Season and a Range overlap @ Begin -> End or End -> Begin

					If there is no overlap, this Slot Data File ( Airport, Carrier, Season ) will not
					contain any Slots relevant to the Query, so continue to next File

					If this Slot Data File IS within the Range, this is considered "passing the Date Filter"
					and the next step will be Flight Number filtering.  ( Indicate that "allDays" will be
					included in any Flight Slot Week result. )
 */
					f1 := inSeason ( rangeStart, begin, end )
					f2 := inSeason ( rangeEnd, begin, end )
					f3 := inSeason ( begin, rangeStart, rangeEnd )
					f4 := inSeason ( end, rangeStart, rangeEnd )

					if !f1 && !f2 && !f3 && !f4 { continue } else { allDays = true }

				} else {	// Operation Date ( not Range )

					foundOne = inSeason ( thisDate, begin, end )

					if !foundOne { continue }		// Skip, Operation Date not in Season

					currentOpDate = thisDate		// One Operation Date Filter
				}

			} else { currentOpDate = "" }			// No Operation Date Filter

/*			This Operation Date is in Slot Data File Season, OR
			The Season is within the Date Range Filter, OR
			No Date Filtering exists
*/
			if firstDateInSeason { slotFilesSearched++ }
			firstDateInSeason = false

/*			---------------------
			FLIGHT SLOT FILTERING
			---------------------

			Include Any Flight Week in the current Slot Data File that matches a Flight in the
			Flight Number filter ( or include all Flights if no Flight Numbers were specified )
*/
			allFlights := false				// Assume Filtering, otherwise include all Flight Weeks

			if len ( includedFlights ) == 0 { allFlights = true }

			flightSlotWeeks := slotDataFileDB[ slotFileKey ]	// Slot Weeks in current Slot Data File

			flightsFound := 0
			slotsChecked := 0

			for _, slotWeek := range flightSlotWeeks {

				if !allFlights { 		// Flight Filtering Required

					foundOne := false

/*					Since there is only one Slot Week / Flight, once a Slot Week has been found
					for all of the Flights in the Flight Number Filter Set, move to next Slot File
 */
					if flightsFound >= len( includedFlights ) { break }

//					Check Flight Slot Week against Flight Number Filter Set

					for _, thisFlightFilter := range includedFlights {

						if thisFlightFilter == slotWeek.flightNumber {

							foundOne = true
							flightsFound++
							break
						}
					}
					if !foundOne { continue }		// Do not include Slot Week in Result Set
				}

/*			    Flight Slot Week passed the Flight ( Number ) Filters ( or none were given )
				Now one or more Slots ( Days ) in this Slot Week must pass filtering
*/
				currentSlot.flightNumber = slotWeek.flightNumber
				currentSlot.opDate = currentOpDate

				_, opDay := dayOfWeek( currentOpDate )		// Day of Week Number ( 1-7 )

				dayOfWeekIndex := 0							// Use to match Day of Week Number
/*
				Examine the Slot Week for Day of Week information.  One Slot ( Day ) / Result Set row.

				If the Query has Operation Date Filter(s), only one Slot ( Day ) of interest
				and is the only one that will be in the Result Set.

				When no Op Date is supplied, then all 7 Slots are evaluated ( for Slot Index )
 */
				for _, slotDay := range slotWeek.slotWeek {

/*					Add Slot to Result Set based on Slot Index Filtering
					Continue to populate Slot Record
*/
					dayOfWeekIndex++

					currentSlot.weekday = slotDay.Weekday
					currentSlot.flightsInSeason = slotDay.FlightsInSeason
					currentSlot.currentMaxCancels = slotDay.CurrentMax
					currentSlot.originalMaxCancels = slotDay.OriginalMax

					currentSlot.slotIndex = slotDay.SlotIndex

					if !allDays { 		// Operation Date Filtering requested

/*					    There is only one Operation Day, is it this Slot Day?
*/
						if opDay == dayOfWeekIndex {

/*							Found Slot Day, before adding it to the Result Set, perform any
							Slot Index Filtering ( if none, include Slot Record in Result Set
							and check next Slot )
*/
							if len( filters.slotIndex ) == 0 {		// No Slot index filters

								slotResultSet = append( slotResultSet, currentSlot )
								break
							}

//							Does this Slot Index pass all of the Filter(s)?

							dayPassed := keepIndex( indexFilter, slotDay.SlotIndex )

							if dayPassed { slotResultSet = append(slotResultSet, currentSlot)
							} else { slotIndexesFailed++ }

							break				// Only one Slot Day / Date; found it, so skip rest of Week

						} else { continue }		// Check next Slot Day
					}

/* 					No Operation Date, so each Slot ( Day ) that passes Slot Index filtering
 					( if there is any ) will be included in Slot Result Set
 */
					if len( filters.slotIndex ) == 0 {		// No Slot Index filtering

						slotResultSet = append( slotResultSet, currentSlot )
						continue
					}
/*
					Slot Filtering required; entire Slot Week is being considered
 */
					slotIndexTest := slotDay.SlotIndex

/*					The final filter is Slot Index.  Query Parameter ( /SlotIndex ) can
 					be one or more comparisons ( greater or less than, equal, and / or
 					within a range ).  Only one comparison needs to pass.
 */
					keep := keepIndex( indexFilter, slotIndexTest )

//					If it passes Slot Index Filter, add Slot ( Day ) to Result Set

					if !keep { slotIndexesFailed++
					} else { slotResultSet = append( slotResultSet, currentSlot ) }
				}
				slotsChecked++
			}

// 			Don't Print for Stress Testing Slot Data Files

			if len ( slotDataFileDB[slotFileKey] ) > 50 {

				fmt.Println ( "Searched:", slotsChecked, "of", len ( slotDataFileDB[slotFileKey] ), "Flight Slots in", slotFileKey )
				if !allFlights { fmt.Println ( "Found:", flightsFound, "Flight Slots(s) in", slotFileKey ) }
			}

/*			If the full Slot Week was included, no Date filter(s) were supplied
			or a Date Range was the filter, Don't continue to test each day in Slot Week.
*/
			if allDays {

				includedDates = []string{}	// Reset to Empty
				break
			}
	    }
	}

	fmt.Println( "Filtered", slotDataFiles, "Slot Files")
	fmt.Println( "Skipped:", oldSlotFiles, "of those Slot Files ( 2-digit years )" )

	usingRange = false		// Any Slot Range Query is complete; reset for future queries

	return slotResultSet
}

/*	A "filter" that includes everything;  No filters
	Obviously, this is used whe /SlotData has no Query Parameters
	It is also used when reloading the Slot File DB from Blob Store
	It is also used during Unit Test
 */
func FilterALL ( slotDataFileDB map[SlotKey][]FlightSlot ) ( []slots ) {

	filter := FilterParams{}

	filter.airportCodes = ""
	filter.carrierCodes = ""
	filter.flightNumbers = ""
	filter.operationDates = ""

	z := filterSlots ( slotDataFileDB , filter )

	return z
}
