package slot_svc

import (

	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"bytes"
	"net/http"

	blobstore "github.build.ge.com/AviationRecovery/blobstore-support-go"

	cl "github.build.ge.com/AviationRecovery/slot-service.git/clients"
)

var (

	SlotFileSummer	string
	SlotFileFall 	string
	SlotFileWinter 	string
	SlotFileSpring 	string

	testReadDB		[]byte

	port			string

	defaultConfigHostName = "https://aviation-dca-config-svc.run.asv-pr.ice.predix.io"

	fromHome		bool
	allSeasons		= false
)


/*=========================================================================
			Drivers (MAINs)
=========================================================================*/
/*
	https://aviation-dca-slot-svc.run.asv-pr.ice.predix.io
	/api/v1
	/SlotData

	    carrierCode=FZ,DL,B6,...
	    &airportCode=DXB,MCO,LAX,JFK,...
	    &flightNumber=7000,11,442,...
	    &operationDate=2016-05-11,2016-05-12,2016-07-31,...
*/


func SlotServiceServer() {

	fmt.Println("-----------------------")
	fmt.Println("Starting Slot Service!!")
	fmt.Println("-----------------------")

//	Slot DB does not exist in memory.  Instantiate the In-Memory DB

	SlotDataDB = make ( map[SlotKey][]FlightSlot )
	seasonList = make ( map[string]seasonDates )

	localAdd := os.Getenv("LOCAL_ADD")
	if localAdd == "" { addLocal = false }

	port = os.Getenv("PORT")

	if port == "" { port = "2525" }

	fmt.Println ( "Listening on Port:", port )

	homeFlag := os.Getenv("HOME")

	fromHome = false
	if homeFlag == "YES" { fromHome = true }

/*	Test Write / Read to / from Blob Store.  This code will be relocated
	and used to persist the Slot File DB.  Read @ Start-up ( Slot History )
	and pull all Slot Data Files for latest new / updated Slot Data

	Right now, the Configuration Service overwrites the seasonal Slot Data Files
	The persistence of the Slot information allows for any history / analytics
	This will also help if the /Cancellation endpoint is every used

	Only use this section when running on Predix.  There will be no BlobStore access
	when running locally.  Also, VCAP Services will NOT work locally.
*/
	if port != "2525" && port != "3535" {

		fmt.Println ( "-------------------------------" )
	    fmt.Println ( "Loading Slot DB from Blob Store" )
		fmt.Println ( "-------------------------------" )

//	    See if what was written earlier is still there....

	    storedURL := "https://bucket-c2cd2eb6-69e3-41c5-8ac7-276fe99ffbfc.store.gecis.io/SlotDB"

	    readBlobStore, err := blobstore.GetObject( storedURL )
	    if err != nil { fmt.Println ( "Error on Blob Store GetObject:", err ) }

		if len ( readBlobStore ) == 0 {	// Slot DB does not Exist

			fmt.Println ( "Slot DB does not Exist, Creating." )

			SlotDataDB = make ( map[ SlotKey ][]FlightSlot )

	    } else {	// SlotDB exists, Reload from Blob Store

			loaded := ReloadSlotDB ( readBlobStore )

			if !loaded { 		// Load failure.  Create.

				fmt.Println ( "Slot DB Load Failed,  Creating" )

				SlotDataDB = make ( map[SlotKey][]FlightSlot )

			} else {

				SlotDataDB = ReloadedSlotDB	// Reload Slot File DB

//			    Retrieve entire Slot DB as denormalized Flight Slots
//			    For validation ( write to Slot Service Logs )

				z := FilterALL( SlotDataDB )
				fmt.Println ( "Reloaded", len(z), "Slots" )
			}

	    }
	}

	if fromHome { fmt.Println ( "Start-up: @ Home" ) }

//	If Variable Init function works ( replaces following code ), then
//	remove this code.  ( Function IS this code; check variable Scope )

//	Configuration Service Slot File: Host Name + Endpoint + Slot File Name

//	Uses CUPS Service to get configurable Host Name


	hn := cl.GetServiceName("CONFIG_URL")
	host := cl.GetServiceHostName(hn)

	fmt.Println ( "-------------------" )
	fmt.Println ( "Config Service Host:", host )

	slotConfigHostName = host

/*	Pull Slot Service Host Name.  If it does not exist, assume running
	Local and use the default testing Host Name:  slot-svc-test
*/
	hn = cl.GetServiceName("SLOT_URL")
	host = cl.GetServiceHostName(hn)

	PredixHost = host

	if PredixHost == "" { 

	    fmt.Println ( "Slot Service Host: Blank ( Defaulting )" )
	    PredixHost = "https://slot-svc-test.run.asv-pr.ice.predix.io" 
	}

	fmt.Println ( "Slot Service Host:", PredixHost )
	fmt.Println ( "-----------------" )


	slotConfigEndpoint = os.Getenv( "SERVICE_ENDPOINTS")
	currentSlotFile = os.Getenv( "SEASON_SUMMER")

//	If an Environment Variable is NOT set, use Defaults

	if slotConfigHostName == "" { 

	    fmt.Println ( "Environment Variable NOT Set:","CONFIG_URL")

	    slotConfigHostName = os.Getenv("CONFIG_SERVICE")

	    if slotConfigHostName == "" {

			fmt.Println ( "Environment Variable NOT Set:","CONFIG_SERVICE")
			slotConfigHostName = defaultConfigHostName
		}
	}

	if slotConfigEndpoint == "" {

	    fmt.Println ( "Environment Variable NOT Set:","SERVICE_ENDPOINTS")
	    slotConfigEndpoint =  "/api/v1/fz/" 
	}

	if currentSlotFile == "" {

	    currentSlotFile = "DBX_SLOTS_SUMMER"
	    fmt.Println ( "Environment Variable NOT Set:","SEASON_SUMMER")
	}

	slotFileURL = slotConfigHostName + slotConfigEndpoint + currentSlotFile	

	fmt.Println("Config Service Slot File:", slotFileURL )

	fmt.Println("Port:", port )

	SlotFileSummer = os.Getenv( "SEASON_SUMMER")
	SlotFileFall = os.Getenv( "SEASON_FALL")
	SlotFileWinter = os.Getenv( "SEASON_WINTER")
	SlotFileSpring = os.Getenv( "SEASON_SPRING")

	if SlotFileSummer == "" { SlotFileSummer = "DBX_SLOTS_SUMMER" }
	if SlotFileFall == "" { SlotFileFall = "DBX_SLOTS_FALL" }
	if SlotFileWinter == "" { SlotFileWinter = "DBX_SLOTS_WINTER" }
	if SlotFileSpring == "" { SlotFileSpring = "DBX_SLOTS_SPRING" }

	fmt.Println( "" )
	fmt.Println( "Config Service Host:", slotConfigHostName )
	fmt.Println( "Config Service Endpoint:", slotConfigEndpoint )
	fmt.Println( "Config Service Slot File:", currentSlotFile )
	fmt.Println( "Summer Slot File:", SlotFileSummer )
	fmt.Println( "" )

	summer := SlotFileSummer 
	fall := SlotFileFall
	winter := SlotFileWinter 
	spring := SlotFileSpring

	seasonFiles = []string{ summer, winter }

	if allSeasons { seasonFiles = []string{ summer, winter, fall, spring } }

// 	Initialize Slot Data DB with Test Data ( keep until Slot DB is loaded )

	fmt.Println ( "-----------------------------" )
	fmt.Println ( "Loading Base Flight Slot Data" )
	fmt.Println ( "-----------------------------" )

/*	Rather than loading Stub / Test Data, load the Slot DB into memory
	and then pull all Slot Data Files ( one / Season ).  Use them to
	bring the Slot DB current.  ( This logic only runs @ Slot Service
	start-up.  In normal mode, the In-Memory DB is used for all Search
	Result Sets.  Basically, Slot Service is read-only, updated when a
	Slot Data File is available.

	Logic exists to have a more real-time Slot Index.  Cancellations can
	be sent to the /Cancellation endpoint for immediate update.
*/

/*	Only issue the following if NO Slot File DB exists in Blob Store
	Commented out ( for now ) until Stub Logic no longer creates one.

	This should also be changed to only do this on a Slot Restart
	Otherwise, the in-memory DB is "erased"
*/

//	Capturing Information ( Slot File Created ( and Updated ) Dates )

//	** This may be overwriting the Created Date for at least one
//	** Slot File in the Slot DB ( This is DB load, so all already
//	** have a created date OR there are none )

	SlotFileInfo = make ( map[ SlotKey ]SlotDataFileDetails )

	if port == "3535" || port == "2525" {

/*	    Use code developed for Go Test routines to load the in-memory 
	    version of Slot File DB ( this test logic will be added to
	    the function that loads the existing Slot File DB from the 
	    Blob Store ( used for Slot Service restarts ) )
*/
	    InitTestReadSlotDB ()

	    SlotDataDB = TestSlotDataDB		// Reload Slot File DB

	    loadStubData ( SlotDataDB )		// Add Stress Size DB
	}


	fmt.Println ( "------------------------------" )
	fmt.Println ( "Loading Base Cancellation Data" )
	fmt.Println ( "------------------------------" )

	initCancelTracking ()
	makeCancellations ()

	fmt.Println ( "Unique Cancels:", len ( UniqueCancels ) )
	fmt.Println ( "Carriers with Cancellations", len ( CarrierCancels ) )
	fmt.Println ( "Flights with Cancellations:", len ( FlightCancels ) )
	fmt.Println ( "Opersation Dates with Cancellations:", len ( OpDate ) )
	fmt.Println ( "Years with Cancellations:", len ( Year 	) )
	fmt.Println ( "Months with Cancellations:", len ( Month ) )
	fmt.Println ( "Weeks of Year with Cancellations:", len ( Week ) )
	fmt.Println ( "Days of the Weeek with Cancellations:", len ( WeekDay ) )
	fmt.Println ( "Unique Types of Cancellation:", len ( Types ) )
	fmt.Println ( "Sources reporting Cancellations:", len ( Source ) )

	fmt.Println ( "-----------------------------" )
	fmt.Println ( "Loaded Base Slot Service Data" )
	fmt.Println ( "-----------------------------" )

	slotsLoaded = false
	dumpAllSlots = false
	showEndPoints = false
	doDebug = false
	initialLoad = false

/*	Load Slot File(s) from Config Service as an update to current
	Slot DB in Blob Store.  This is a restart, and likely the Slot DB
	is still current, but do an update check just in case.
*/
	fmt.Println( "Predix Service Host:", slotConfigHostName )
	fmt.Println( "Predix Service Endpoint:", slotConfigEndpoint )
	fmt.Println( "Predix Service Slot File:", currentSlotFile )

//	Load Summer and Winter once before launching Monitoring Agent

	if port != "3535" && port != "2525" {

		LoadSlotFile ( "Slot Server", seasonFiles )
	}

	cancellationSimulation := false			// Enable when Cancellation Simulation ( aid in solution definition )

/*	Create a Monitoring agent that requests Slot Files for all Season
	from Config Service.  These fies change so infrequently, the makes
	little sense to set up and use a Message Queue
 */

/*	Simulate Cancellations?
 */
	if cancellationSimulation {

		fmt.Println ( "-----------------" )
		fmt.Println ( "Test Cancellation" )
		fmt.Println ( "-----------------" )
		cancelSlot ( "2016-May-31", "768" )
		cancelSlot ( "2016-May-31", "41" )
		cancelSlot ( "2016-May-31", "752" )

		go pollSlotFiles ( "Starting: Slot File Polling" )

		go generateCancels ( "Starting: Cancellation Generation" )
	}



	fmt.Println ( "--------------------------------------------" )
	fmt.Println ( "SlotService Handlers Set, Listening on Port:", port )
	fmt.Println ( "--------------------------------------------" )

	http.ListenAndServe(":"+port, nil)
}

/*
	Slot File Monitoring Agent
 */

//	Only two Season Files are created / available ( according to FlyDubai )

var (

	seasonFiles  		[]string

	PredixHost 	string

	devHost = "https://slot-dev.run.asv-pr.ice.predix.io"
	localHost = "http://localhost:2525"

	loadEndpoint 	= "/LoadSlotFile"
	cancelEndpoint 	= "/Cancellation"

)

/*	Function for cancellation simulation / testing
 */

func generateCancels ( msg string ) {

	fmt.Println ( "Start:", msg )

//	Hit Cancellation Endpoint @ regular intervals
//	Use Port == 2525 to ID running locally vs. on Predix
//	Only run Local, disable for Predix

	if port != "2525" && port != "3535" { return }

	cancel := localHost + cancelEndpoint 	// Local

//	Create Cancels during Local Test, up to Max Cancels ( Env Variable )

	maximumCancels := os.Getenv ( "CANCELS" )
	maxCancels, _ := strconv.Atoi ( maximumCancels )

	if maxCancels == 0 { return }		// No Cancels

	for i := 1; i <= maxCancels; i++ {

	    time.Sleep( 1 * time.Second )	// One Second

	    fmt.Println ( "Creating Cancellation:", i )

	    _, err := http.Get ( cancel )
	    if err != nil { fmt.Println ( "Cancellation Failed:", err ) }
	}
	fmt.Println ( "End:", msg )
}


func FilterSlotIndex ( response http.ResponseWriter, query *http.Request ) {

	slotIndexFilters := query.URL.Query()

	filterString := ""
	for _, values := range slotIndexFilters {
		for _, filter := range values { filterString = filter }
	}

	siFilters := strings.Split ( filterString, "," )

	fmt.Println ( "Slot Index Filters:", siFilters )

	ifWhat := "If: Conditial Missing"

	lastOne := len ( siFilters ) - 1
	for ifThis, filter := range siFilters {

		offset := 2

		if ifThis == 0 { fmt.Println ( "Will Filter into Result Set" ) }

//		fmt.Println ( "Filter:", ifThis, "Instance:", filter )
//		fmt.Println ( "Filter[0]:", filter[0:1] )

		switch filter[0:1] {
		case "G" : ifWhat = "If: Greater than"
		case "L" : ifWhat = "If: Less than"
		case "E" : ifWhat = "If: Equal to"
		case "R" : ifWhat = "If: within Range"
			offset = 5
		}
		if ifThis < lastOne {
			fmt.Println ( ifWhat, filter[offset:], "OR" )
		} else { fmt.Println ( ifWhat, filter[offset:] ) }

//		fmt.Println ( "Condition;", filter[0:1], "Length of all Filters:",len ( siFilters ) )
	}
	fmt.Println ( "include in Result Set" )
}

/*	------------------------------------------------------------------------------
	Added Support for Cancellations, this will parse the SlotIndex Query Parameter
	------------------------------------------------------------------------------

	/Cancellation		Cancellation Endpoint

	Cancellations endpoint will have a Query String

	Parameters, Example Values
	--------------------------
	Short Form Parameters / Examples ( long form ( e.g., AirportCode ) parameter names supported )

	AC=MCO					// Airport Code ( same as Flight Origin )
	FD=MCO					// Flight Origin
	FO=MCO					// Flight Destination
	CC=DL					// Carrier Code
	OD=2016-05-31			// Operation / Cancel Date
	FN=300					// flight Number
	CT=Weather				// Cancel Category ( Type )
	CR=Snow					// Cancel Reason
	CN=free text			// Cancel Notes
	CS=Tower				// Cancel Source
*/

func Cancellation ( response http.ResponseWriter, query *http.Request ) {

	fmt.Println ( "Hit Cancellation Endpoint, Cancelling:" )

	queryParams := query.URL.Query() 	// Read Query String

	cancelledFlight := parseQueryString ( queryParams )
	
	ac := cancelledFlight.airportCodes
	cc := cancelledFlight.carrierCodes
	fn := cancelledFlight.flightNumbers

//	For now, assume / support ONE Operation Date

	od := cancelledFlight.operationDates
	fo := cancelledFlight.flightOrigin
	fd := cancelledFlight.flightDestination

	fmt.Println ( "Carrier:", cc )
	fmt.Println ( "Airport:", ac )
	fmt.Println ( "Origin:", fo )
	fmt.Println ( "Destination:", fd )
	fmt.Println ( "Flight:", fn )
	fmt.Println ( "Operation Date:", od )

	cancelSlot ( cancelledFlight, response )

//	Add Cancellation to Cancellation DB

	ct := cancelledFlight.cancelCategory
	cr := cancelledFlight.cancelReason
	cn := cancelledFlight.cancelNotes
	cs := cancelledFlight.cancelSource

	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )
}

var	slotCancels	= map[ cancelKey ]int{{"MCO","MCO","LAX","DL","20","2012-May-31"}:0}

func cancelSlot ( cancelFilter FilterParams, response http.ResponseWriter ) {

//	Handle a Cancellation

/*	Will need to Parse the Query String ( very similar to /SlotData,
	but there is no search, and only One Key )

	Get Key ( Parse Query String )
*/
	fmt.Println ( "----------------------" )
	fmt.Println ( "Cancellation Requested" )
	fmt.Println ( "----------------------" )

//	AirportCode := "DXB"
//	CarrierCode := "FZ"
//	operationDate := cancelDate

	AirportCode := cancelFilter.airportCodes
	CarrierCode := cancelFilter.carrierCodes

//	WTH HACK

	AirportCode = "DXB"
	CarrierCode = "FZ"

	operationDate := cancelFilter.operationDates

	flightNumber := cancelFilter.flightNumbers

	cancelReason := cancelFilter.cancelReason
	cancelNotes := cancelFilter.cancelNotes

	cancelledSlot := cancelKey{}

	cancelledSlot.airport = AirportCode 
	cancelledSlot.carrier = CarrierCode 
	cancelledSlot.flight = flightNumber 
	cancelledSlot.opDate = operationDate 

/*	Update Slot w/i correct Flight Slot

	flightSlot := SlotDataDB[key]
	Position to Op Date w/i Slot Week
	Update Cancels by decrementing Max Cancels
*/

/*	Search Slot File DB for Slot File that contains the Flight Slot
	that experienced the Cancellation ( one of the Slot Days )
*/
	for k1, _ := range SlotDataDB {

	    if k1.Airport != AirportCode { continue }
	    if k1.Carrier != CarrierCode { continue }

/*	    Find the Season the contains the OpDate
	    This will use the Season Begin / End Dates
*/
	    begin := k1.SeasonStart
	    end := k1.SeasonEnd

	    foundOne := inSeason ( operationDate, begin, end )
	    if !foundOne { continue } 

//	    Found the Season for the Operation Date

	    s := SlotKey{}

	    s.Airport = AirportCode 
	    s.Carrier = CarrierCode 
	    s.SeasonStart = begin
	    s.SeasonEnd = end

//	    Now, find the Flight Slot

	    flightSlots := SlotDataDB[s]

	    for _, v2 := range flightSlots {

			//		Once Flight Slot is located, Get Slot Day

			//fmt.Println ( "WTH In Slot Check:", v2.flightNumber, "=", flightNumber )
			//fmt.Printf ( "[%s:%s]", v2.flightNumber, flightNumber )

			if v2.flightNumber == flightNumber {

				cancelDay, _ := dayOfWeek ( operationDate )

				fmt.Println ( "WTH Found:", v2.flightNumber, "=", flightNumber )
				fmt.Println ( "WTH Cancel Day:", cancelDay, "Op Date:",  operationDate )


				fmt.Println ( "Cancel", cancelDay, "in Flight Slot:", v2 )

				slotWeekIndex := 0
				for _, slot := range v2.slotWeek {

					if slot.Weekday == cancelDay {

						fmt.Println ( "Postioned to Cancel:" )
						fmt.Println ( "WeekDay:", slot.Weekday )
						fmt.Println ( "Current Max:", slot.CurrentMax )
						fmt.Println ( "Slot Index:", slot.SlotIndex )

						count := slotCancels[cancelledSlot]

						count++
						slotCancels[cancelledSlot] = count

						fmt.Println ( "Cancel Reason:", cancelReason )
						fmt.Println ( "Cancel Notes:", cancelNotes )

						if count != 1 {

							fmt.Println ( "Cancel", count, "for this Slot" )

							break		// Not a new cancel
						}

						oMax := slot.OriginalMax
						flt := CarrierCode + flightNumber

						newMax := slot.CurrentMax

						if slot.CurrentMax > 0 {

							newMax = newMax - 1
						}
						newIndex := ComputeSlotIndex ( newMax, oMax, flt )

						fmt.Println ( "New Current Max:", newMax )
						fmt.Println ( "New Slot Index:", newIndex )

						v2.slotWeek[slotWeekIndex].CurrentMax = newMax
						v2.slotWeek[slotWeekIndex].SlotIndex = newIndex

						break
					}
					slotWeekIndex++

				}


			} else { continue }

	    }

	}

//	For now, reuse Search Filter to get updated Slot into desired format
//	Likely a better way vs. searching for was we already have...

	fmt.Println ( "Start Filtering ( /Cancellation )" )

	z := filterSlots ( SlotDataDB, cancelFilter )

	fmt.Println ( "End Filtering ( /Cancellation )" )

	fmt.Println ( "Found", len(z), "Slots" )

//	Filtering is Complete.  Creating JSON from Result Set...

	flatJSON := denormSlotsFlat ( z, true )

	status := http.StatusOK		// Default ( Result Set non-Empty )

//	Check for an empty Result Set ( Query returned nothing )

	if len( z ) == 0 { status = http.StatusNoContent }  // Empty

	response.WriteHeader( status )

//	Add JSON Header

	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

//	Write Response ( Result Set )

	fmt.Fprintf ( response, flatJSON ) 
}


/*
	Slot Data File polling ( from Config Service )
	Use an increasing poll delay...
 */
var		sleepTime = []time.Duration{ 10, 20, 30, 60, 120, 300, 600, }

func pollSlotFiles ( msg string ) {

//	fmt.Println ( "-----------------------" )
//	fmt.Println ( "Start Slot File Loader:", msg )
//	fmt.Println ( "-----------------------" )

//	Hit SlotFileReload Endpoint @ regular intervals

//	Use Port == 2525 to run locally
//	Use Port == 3535 to simulate running on Predix
// 	Use "hack" endpoint ( while testing on Predix )

//	Switch between "real" name of Slot Service and Development Host

	reload := PredixHost + loadEndpoint	// On Predix ( real Host )
//	reload = devHost + loadEndpoint	// On Predix ( fake Host )

//	Testing Endpoints ( Fully Local is default )

//	if port == "2525" { reload = localHost + loadEndpoint }	// Local
	if port == "3535" { reload = devHost + loadEndpoint }	// Predix

/*	Ever increasing delays @ start-up / restart until last entry in
	Delay List ( sleepTime ( slice of Durations ) ) and then use that
	Duration to loop forever ( until Slot Service is stopped )
*/
//	Right Now, Max Delay is 10 minutes ... ( probably should be longer )

	foreverDelay := sleepTime[0]

	for _, delay := range sleepTime {

//	    where := "Locally"
//	    if port == "2525" { where = "Slot File"  }

	    foreverDelay = delay		// Save "final" delay amount

	    sleeping := delay * time.Second
	    time.Sleep( sleeping ) 

//	    fmt.Println ( "Poll [", where, "] after", sleeping, "seconds" )

	    if fromHome { continue }

	    _, err := http.Get ( reload )

	    if err != nil { 

//		fmt.Println ( "Load File Poll Failed:", err ) 

	    } else { 

//		fmt.Println ( "Slot File Polled / Loaded from", reload )
	    }
	}
	fmt.Println ( "----------------" )
	fmt.Println ( "Slot Monitoring:", foreverDelay, "Seconds" )
	fmt.Println ( "----------------" )

	for i := 1; i <= 20; i++ {

	    if fromHome { break }

	    sleeping := foreverDelay * time.Second
	    time.Sleep( sleeping ) 

//	    where := "Forever"
//	    fmt.Println ( "Poll [", where, "] after", sleeping, "seconds" )

	    _, err := http.Get ( reload )

	    if err != nil { 

//		fmt.Println ( "Load File Poll Failed:", err ) 

	    } else { 

//		fmt.Println ( "Slot File Polled / Loaded" ) 
//		fmt.Println ( "Source of GET:", reload )
	    }

	    i = 5
	    continue		// Pull forever
	}
	fmt.Println ( "-----------------------------------" )
	if fromHome { fmt.Println ( "Skipped Infinite Polling ( at Home )" ) 
	} else { fmt.Println ( "Ending Infinite Polling is an ERROR" ) }
	fmt.Println ( "-----------------------------------" )
}

func SlotFileLoad ( response http.ResponseWriter, query *http.Request ) {

	fmt.Println ( "------------------" )
	fmt.Println ( "Loading Slot Files" )
	fmt.Println ( "------------------" )

	queryParams := query.URL.Query() 	// Read Query String

//	If only one:  ?season=SUMMER ( or WINTER )

	seasons := []string{}
	queryParameters := ""
	for _, value := range queryParams {	// No Key validation

	    for _, queryParameters = range value { }
	}

	if queryParameters != "" { 	// One Season requested, default if invalid

	    seasons = append ( seasons, seasonFiles[0] )	// Default

	    if queryParameters == "SUMMER" { seasons = append ( seasons, seasonFiles[0] )}
	    if queryParameters == "WINTER" { seasons = append ( seasons, seasonFiles[1] )}

	} else { seasons = seasonFiles }

	fmt.Println ( " Issuing LoadSlotFile (", seasons, ")" )

	LoadSlotFile ( "Slot Data Pull", seasons )

	action := "] Pulled and Loaded"

	fmt.Println ( "----------" )
	fmt.Println ( "Slot Files [", seasons, action )
	fmt.Println ( "----------" )

//	Build string for Response Body ( doesn't do as expected otherwise )

	slotFiles := ""
	for _, file := range seasonFiles { slotFiles = slotFiles + file + " " }

	loadResult := "Slot File(s) [ " + slotFiles + "] Pulled and Loaded"

	fmt.Fprintf ( response, loadResult )
}

//	Load Requested Seasons

func LoadSlotFile ( requestedBy string, seasons []string ) {

	slotFilesPulledTime := time.Now()

	fmt.Println ( "--------------------------------------------" )
	fmt.Println ( "Loading Seasonal Slot Files:" )
	fmt.Println ( "Start Time:", slotFilesPulledTime )
	fmt.Println ( "Requested By:", requestedBy )
	fmt.Println ( "Slot Files: [", seasons, "]" )
	fmt.Println ( "--------------------------------------------" )

//	Pull each requested Season's Slot File, compare with current Slot DB
//	Update as needed ( new Flight Slots, Slot Week changes )

//	if Slot Server, this is the Initial Load, logging is ON
//	if User, this is an On-Demand ( File-specific ( or All ) request

	loadType := "First Update Check:"	// Slot Server

	if requestedBy == "User" { loadType = "User Query:" }
	if requestedBy == "Agent" { loadType = "Agent Loading:" }

	silent := false

	for _, slotFile := range seasons {

	    slotFileURL = slotConfigHostName + slotConfigEndpoint + slotFile

	    currentSlotFile = slotFile

	    if requestedBy == "Agent" { silent = true }

	    fmt.Println( loadType, slotFile ) 

	    getSlotDataFilePredix ( silent )

	    fmt.Println( "Loaded:", currentSlotFile ) 
	    fmt.Println( "from:", slotFileURL )
	}


	slotFilesPulledTime = time.Now()

	fmt.Println ( "-------------------------" )
	fmt.Println ( "Latest Slot Files Pulled," )
	fmt.Println ( "... Checked for Updates," )
	fmt.Println ( "... and Slot DB Updated." )
	fmt.Println ( "End Time:", slotFilesPulledTime )
	fmt.Println ( "-------------------------" )
}


/*	Initialize Primary Slot Service variables for Configuration Service
	Interaction ( for Seasonal Slot Files )
*/

func loadSlotServiceEnviroment () {

	//	Configuration Service Slot File: Host Name + Endpoint + Slot File Name

	//	Uses CUPS Service to get configurable Host Name

	hn := cl.GetServiceName("CONFIG_URL")
	host := cl.GetServiceHostName(hn)

	fmt.Println("Config Service Host:", host)

	slotConfigHostName = host

	//	Can I get the Slot Name?

	hn = cl.GetServiceName("SLOT_URL")
	host = cl.GetServiceHostName(hn)

	fmt.Println("Slot Service Host:", host)

	PredixHost = host

	//	Enable when Testing on Predix ( slot-svc-test )

	//	PredixHost = "https://slot-svc-test.run.asv-pr.ice.predix.io"

	slotConfigEndpoint = os.Getenv("SERVICE_ENDPOINTS")
	currentSlotFile = os.Getenv("SEASON_SUMMER")

	//	If an Environment Variable is NOT set, use Defaults

	if slotConfigHostName == "" {

		fmt.Println("Environment Variable NOT Set:", "CONFIG_URL")

		slotConfigHostName = os.Getenv("CONFIG_SERVICE")

		if slotConfigHostName == "" {
			fmt.Println("Environment Variable NOT Set:", "CONFIG_SERVICE")
			slotConfigHostName = defaultConfigHostName
		}
	}

	//	HACK ( Temporarily force Host when running on Predix )
	//	This may be fixed from updating UAA package ( per Wayne )

	if port != "2525" && port != "3535" {
		slotConfigHostName = defaultConfigHostName
	}

	//	Default Config Service endpoint for Slot Data Files

	if slotConfigEndpoint == "" {

		fmt.Println("Environment Variable NOT Set:", "SERVICE_ENDPOINTS")
		slotConfigEndpoint = "/api/v1/fz/"
	}

	//	Defaults to Summer; subsequent code varies as Seasons change

	if currentSlotFile == "" {

		currentSlotFile = "DBX_SLOTS_SUMMER"
		fmt.Println("Environment Variable NOT Set:", "SEASON_SUMMER")
	}

	slotFileURL = slotConfigHostName + slotConfigEndpoint + currentSlotFile

	fmt.Println("Slot File from Config Service:", slotFileURL)

	fmt.Println("Port:", port)

//	Slot Data Files are Seasonal, presently, there are only two Seasons:
//	Summer and Winter.  Support all four.

	SlotFileSummer = os.Getenv("SEASON_SUMMER")
	SlotFileFall = os.Getenv("SEASON_FALL")
	SlotFileWinter = os.Getenv("SEASON_WINTER")
	SlotFileSpring = os.Getenv("SEASON_SPRING")

//	Defaults

	if SlotFileSummer == "" { SlotFileSummer = "DBX_SLOTS_SUMMER" }
	if SlotFileFall == "" { SlotFileFall = "DBX_SLOTS_FALL" }
	if SlotFileWinter == "" { SlotFileWinter = "DBX_SLOTS_WINTER" }
	if SlotFileSpring == "" { SlotFileSpring = "DBX_SLOTS_SPRING" }

	fmt.Println("Config Service Host:", slotConfigHostName)
	fmt.Println("Config Service Endpoint:", slotConfigEndpoint)
	fmt.Println("Config Service Slot File:", currentSlotFile)
	fmt.Println("Summer Slot File:", SlotFileSummer)

//	Setup Seasons Lists ( current and future )

	summer := SlotFileSummer
	winter := SlotFileWinter


/*	For now, only Summer and Winter are created; in the future, likely all Seasons will be created
 */
	seasonFiles = []string{summer, winter }

	if allSeasons {

		fall := SlotFileFall
		spring := SlotFileSpring
		seasonFiles = []string{summer, winter, fall, spring }
	}

}



//	-------------------------------------
//	Test Code ( Blob Store Read / Write )
//	-------------------------------------

//	Interact with Blob Store ( SlotFile DB is stored in Blob Store )
//	Function likely unused.  ( Delete? )

func	BlobStoreReadWrite ( fakeDB string ) ( bool ) {

//	This tests Slot Data DB write.  ( One string )

	fmt.Println ( "Blob Store Write:", fakeDB )

	b := bytes.NewBufferString ( fakeDB )

	storedUrl, err := blobstore.PutObject ( "SimpleTest", b )

	if err != nil { 
	    fmt.Println ( "Error on Blob Store PutObject:", err ) 
	    return false
	}

	fmt.Printf("URL of stored object is: %s\n", storedUrl)

	DB2 = storedUrl		// Save ( Permanent DB in Blob Store location )

	fmt.Printf("Done saving Slot DB\n")

//	Test Blob Store Read ( GetObject )

	testReadDB, err = blobstore.GetObject( storedUrl )
	if err != nil { fmt.Println ( "Error on Blob Store GetObject:", err ) }

/*	For now, just a string was stored / read.  ( See Slot Read / Write
	logic for loading JSON into the memory-resident Slot File DB )
*/
	DBasString := string(testReadDB[:])

	fmt.Printf("Blob Store Read: %s\n", DBasString )

	DBwritten = DBasString

	return true
}

/*
	DBasString := string(readBlobStore[:])
	fmt.Println("Length of Blob Store Restore Content", len ( DBasString ) )
	if len ( DBasString ) < 200 { fmt.Printf("Blob Store Restore: %s\n", DBasString ) }
*/

/*	    This code is in the function ( above )

	    b := bytes.NewBufferString( "Slot DB Content" )
	    storedUrl, err := BlobStore.PutObject("SlotDB", b )
	    if err != nil { 
		fmt.Println ( "Error on Blob Store PutObject:", err )
	    }
	    fmt.Printf("URL of stored object is: %s\n", storedUrl)

	    DB2 = storedUrl	// Save ( Permanent DB in Blob Store location )

	    fmt.Printf("Done saving Slot DB\n")

//	    Test Slot Data DB read

	    testReadDB, err = BlobStore.GetObject( storedUrl )
	    if err != nil { 
		fmt.Println ( "Error on Blob Store GetObject:", err )
	    }

//	    For now, just a string was stored / read.  ( See Slot Read / Write
//	    logic for loading JSON into the memory-resident Slot File DB )
//
	    DBasString := string(testReadDB[:])

	    fmt.Printf("Content Written: %s\n", DBasString )

	    DBWritten = DBasString
*/




/***********************************
//	Delete a specific Slot File from Slot DB
//	Requires complete Key ( AirportCode, CarrierCode, Season )

func DeleteSlotFile ( response http.ResponseWriter, query *http.Request ) {

//	Get Delete Parameters ( e.g., ?AC=DXB&CC=FZ&SB=2016-May-31&SE=2016-Dec-31 )


//	Make sure parameters are non-null ?

//	Make sure Slot File exists in DB ( serves as validation? )

//	if Not Found, 
		fmt.Println ( "Slot File", key, "not found" )
		return

	fmt.Println ( "Deleting:", len ( DB[key] ), "Slots from", key )

//	If found, Delete

	delete ( slotDB, key )	// Put in correct DB and Key names

//	Just like new Slot Files that change ( add, update ) the Slot DB,
//	Write updated Slot DB

	writeDB( slotDB )

	return

}

*******************************/

	
//var	airportList = []string{"MCO", "DXB", "LAX", "LAX", "LAX", "JFK", "MCO", "DXB", }

type	seasonDateKey 	struct {

	begin		string
	end			string
}

func SlotInfo ( response http.ResponseWriter, query *http.Request ) {
	
	ac := make ( map[ string ]int )
	cc := make ( map[ string ]int )
	s := make ( map[ seasonDateKey ]int )

	seasonKey := seasonDateKey{}

	out := ""

	slotsInFile := 0

	for key, _ := range SlotDataDB {

	    fmt.Println ( "Slot File:", key, "Flight Slots:", len ( SlotDataDB[key] ) )
	    out = fmt.Sprintln ( "Slot File:", key, "Flight Slots:", len ( SlotDataDB[key] ) )
	    fmt.Fprintf ( response, out )

//	    Airports

	    slotsInFile = len ( SlotDataDB[ key ] )

	    currentSlots := ac[key.Airport]

	    currentSlots = currentSlots + slotsInFile
	    ac[key.Airport] = currentSlots

//	    Carriers

	    currentSlots = cc[key.Carrier] 

	    currentSlots = currentSlots + slotsInFile
	    cc[key.Carrier] = currentSlots

//	    Seasons

	    seasonKey.begin = key.SeasonStart
	    seasonKey.end = key.SeasonEnd

	    currentSlots = s[seasonKey]

	    currentSlots = currentSlots + slotsInFile
	    s[seasonKey] = currentSlots 
	}

	fmt.Println ( "" )
	fmt.Fprintf ( response, "\n" )

	flightList := []string{}

	slotIndexAll := 0
	slotIndexOne := 0
	slotIndexZero := 0
	slotIndexOther := 0
	lowest := 1.0
	highest := 0.0	

	for key, flightSlots := range SlotDataDB {

	    slotIndexAll = 0
	
	    slotIndexOne = 0
	    slotIndexZero = 0
	    slotIndexOther = 0
	    lowest = 1.0
	    highest = 0.0

	    for _, flightSlot := range flightSlots {

			for _, slot := range flightSlot.slotWeek {

				slotIndexAll++
				if slot.SlotIndex == 1.0 { slotIndexOne++
				} else if slot.SlotIndex == 0.0 { slotIndexZero++
				} else {

					slotIndexOther++
					flightList = append ( flightList, flightSlot.flightNumber )
				}
				if slot.SlotIndex != 0 && slot.SlotIndex < lowest { lowest = slot.SlotIndex }
				if slot.SlotIndex != 1 && slot.SlotIndex > highest { highest = slot.SlotIndex }
			}

		}

//	    Set Precision to two decimals

	    low := fmt.Sprintf("%.2f", lowest )
	    printLow,_ := strconv.ParseFloat(low, 2)

	    high := fmt.Sprintf("%.2f", highest )
	    printHigh,_ := strconv.ParseFloat(high, 2)

//	    Print to Log and as REST Response

	    fmt.Println ( "Slot File:", key, "Slots[", slotIndexAll, "]  Indexes[ 0:", slotIndexZero, "] [ 1:", slotIndexOne, "] [ From", printLow, "to", printHigh, ":", slotIndexOther, " ]" )

	    out = fmt.Sprintln ( "Slot File:", key, "Slots[", slotIndexAll, "]  Indexes[ 0:", slotIndexZero, "] [ 1:", slotIndexOne, "] [ From", printLow, "to", printHigh, ":", slotIndexOther, " ]" )

	    fmt.Fprintf ( response, out )

//	    Log / Return all Flights with Slot Index other than 0 and 1

/*	    Unused ( for now )

	    fmt.Println ( "Slot File:", key, "Flights with Slot Index Other:", flightList )
	    out = fmt.Sprintln ( "Slot File:", key, "Flights with Slot Index Other:", flightList )
	    fmt.Fprintf ( response, out )
	    fmt.Fprintf ( response, "\n" )
*/
	}

//	Log / Report Number of Airports, Carriers, and Seasons
//	Do we want Slots / Airport, Slots / Carrier, Slots / Season?

	acList := []string{}
	ccList := []string{}
	seasonList := []seasonDateKey{}
	
	for aCode, _ := range ac { acList = append ( acList, aCode ) }
	for cCode, _ := range cc { ccList = append ( ccList, cCode )  }
	for sCode, _ := range s { seasonList = append ( seasonList, sCode ) }
	
	fmt.Println ( "Airports:", len ( ac ) , acList )
	out = fmt.Sprintln ( "Airports:", len ( ac ) , acList )
	fmt.Fprintf ( response, "\n" )
	fmt.Fprintf ( response, out )

	for _, airport := range acList {

	    fmt.Println ( airport, ":", ac[ airport ], "Flight Slots" )
	    out = fmt.Sprintln ( airport, ":", ac[ airport ], "Flight Slots" )
	    fmt.Fprintf ( response, out )
	}

	fmt.Println ( "Carriers:", len ( cc ) , ccList )
	out = fmt.Sprintln ( "Carriers:", len ( cc ), ccList )
	fmt.Fprintf ( response, "\n" )
	fmt.Fprintf ( response, out )

	for _, carrier := range ccList {

	    fmt.Println ( carrier, ":", cc[ carrier ], "Flight Slots" )
	    out = fmt.Sprintln ( carrier, ":", cc[ carrier ], "Flight Slots" )
	    fmt.Fprintf ( response, out )
	}

	fmt.Println ( "Seasons:", len ( s ) , seasonList )
	out = fmt.Sprintln ( "Seasons:", len ( s ) , seasonList )
	fmt.Fprintf ( response, "\n" )
	fmt.Fprintf ( response, out )

	for _, season := range seasonList {

	    fmt.Println ( season, ":", s[ season ], "Flight Slots" )
	    out = fmt.Sprintln ( season, ":", s[ season ], "Flight Slots" )
	    fmt.Fprintf ( response, out )
	}
}
