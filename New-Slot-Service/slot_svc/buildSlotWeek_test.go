package slot_svc

import (
	"fmt"
//	"math"
	"encoding/json"
	"testing"
)
var	dayFails = map[string]int{"Sunday":0,"Monday":0,"Tuesday":0,"Wednesday":0,"Thursday":0,"Friday":0,"Saturday":0}

var	slotFails 		int	// Total Slot Failures
var	daysFailed 		int	// Day Failures are reset evey week
var	totalSlotFailures	int	// Total Failures of any type

func TestBuildSlotWeek ( t *testing.T ) {

//	Intitialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

	fmt.Println ( "-------------------------------------" )

	fmt.Println ( "Starting Testing of BuildSlotWeek()" )

//	Simulate GET to Config Service, creating response Body

	contents, readOK := readSlotJSON( "Slots.JSON" )
	if !readOK {
	    fmt.Println ( "JSON read fail, Contents size:", len ( contents ) )
	    t.Fail() 
	} 

	blob := []SlotJSON {}

	err := json.Unmarshal ( contents, &blob )
	if err != nil { 
	    fmt.Println ( "Unmarshll fail:", err )
	    t.Fail() 
	} 

//	Simulate Call to LoadSlotFileJSON ( Slot DB, Key, Slot Rows )

/* 	f 		map[SlotKey][]FlightSlot
	slot1 		SlotKey
	slotRows 	[]SlotJSON 
*/
	
//	f := TestSlotDataDB

	slot1 := SlotKey{}

// 	Arbitrary Slot File Keys for this Test

	slot1.Airport = "DXB"		
	slot1.Carrier = "FZ"
	slot1.SeasonStart = "15-Oct-25"
	slot1.SeasonEnd = "16-Mar-26"

	slotRows := blob	// One Slot Row per Flight in Slot File

//	Intialize Results Counters

	flightSlots := 0	// Total Filght Slots Tested
	slots := 0		// Total Slots Tested

	flightSlotFails := 0	// Total Filght Slot Failures


	maxDayFails  := 0	// Checked each week, higher saved

//	Track Fails for Each Weekday

//	flightSlotPassed := true

	fmt.Println ( "Calling BuildSlotWeek() for Slot File:", slot1 )
	fmt.Println ( "Testing", len ( slotRows ), "Flight Slot ( One Flight, One week )" )

//	Make Repetitive Calls to BuildSlotWeek()

	for _, fs := range slotRows { 

	    flightSlots++	// Another Flight Slot tested

//	    flightSlotPassed = true

//	    Build Flight Slot ( one week )

	    ingestSlotWeek := BuildSlotWeek ( slot1, fs )

/*	    Can we get the Flight Number of the Flight Slot that contains
	    this Slot Week?
*/

//	    Validate expected contents of Slot Week

/*	    Each fs ( range slotRows ( type slotJSON ) ) is the "wide"
	    Flight Slot ( Flights in Season and Current Max Cancels fields
	    for each Week Day ( SunFlightsInSeason, etc. ) and must be
	    compared to the Flight Slot that is a slice of sevem Week Days

	    Both start with Flight Number, the latter is a Flight Slot 
	    struct ( F#, []Week Days ) ... walk the Slice and count for
	    the appropriate Week Day
*/
	    success := true
	    daysFailed  := 0	// Days Failed THIS week

//	    Check the Week, day by day

	    for day, slotDay := range ingestSlotWeek {

		slots++		// Another slot being tested

/*		There will be a slice of Flight Slots for this Slot File
		in the Slot File DB.  Assumes that the slice will be ordered
		as it was created.
*/
		fis := fs.SunFlightsInSeason
		cMax := fs.SunCurrentMax
		DoW := "Sunday"

		switch day {

		case 0:		// Sunday

		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }

		case 1:		// Monday

		    fis = fs.MonFlightsInSeason
		    cMax = fs.MonCurrentMax
		    DoW = "Monday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 2:		// Tuesday

		    fis = fs.TueFlightsInSeason
		    cMax = fs.TueCurrentMax
		    DoW = "Tuesday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 3:		// Wednesday

		    fis = fs.WedFlightsInSeason
		    cMax = fs.WedCurrentMax
		    DoW = "Wednesday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 4:		// Thursday

		    fis = fs.ThuFlightsInSeason
		    cMax = fs.ThuCurrentMax
		    DoW = "Thursday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 5:		// Friday

		    fis = fs.FriFlightsInSeason
		    cMax = fs.FriCurrentMax
		    DoW = "Friday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 6:		// Saturday

		    fis = fs.SatFlightsInSeason
		    cMax = fs.SatCurrentMax
		    DoW = "Saturday"
		    success = validateSlot( slotDay, fs, fis, cMax, DoW, "AAA" ) 
		    if !success { t.Fail() }
		    
		case 7:

		}

		if maxDayFails < daysFailed { maxDayFails = daysFailed } 

	    }


//	    Results for this Week ( Flight Slot )

	    if !success { 

		flightSlotFails++

		fmt.Println ( "FAIL: BuidSlotWeek()" )
		fmt.Println ( "Called w/Flight Slot:", fs )
		fmt.Println ( "Week Retuned / Ingested:", ingestSlotWeek )

	    }
	}
	fmt.Println ( "Finished Testing of BuildSlotWeek()" )

//	Report Results

	fmt.Println ( "Flight Slots Failed:", flightSlotFails, "out of", flightSlots )
	fmt.Println ( "Slots Failed:", slotFails, "out of", slots )
	fmt.Println ( "Total Failes (this Slot File) of any type:", totalSlotFailures )

//	Walk the Map ( twice ) to (1) capture Best and Worst # of Fails, and
//	(2) Add day(s) to appropriate list

	worst := maxDayFails
	best := 0

	worstList := []string{}
	bestList := []string{}

	for _, count := range dayFails {

	    if best > count { best = count }
	    if worst < count { worst = count }
	}
	for day, count := range dayFails {

	    if count == best { bestList = append ( bestList, day ) }
	    if count == worst { worstList = append ( worstList, day ) }
	}

	fmt.Println ( "Day(s) with Most (", worst,") Failures:", worstList )
	fmt.Println ( "Day(s) with Least (", best," ) Failures:", bestList )
}


/*	Validate One Slot ( One Flight, One Day ) for matches of Flights in
	Season, Current Max Cancels, and Day of Week (string).  Also compute
	Slot Index and Original Max Cancels ( by the rules ) for matches.

	Report number of mismatches for the Flight Slot, and record the number
	of Flight Slots that fail ( total failed Flight Slots for Sile File )
*/ 

func validateSlot ( slotDay SlotRecord, fs SlotJSON, flightsInSeason int, currentMax float64, weekDay, prefix string ) ( bool ) {

	success := true
	fails := 0
		    
	if slotDay.Weekday != weekDay { 

	    fmt.Println ( "FAIL:", slotDay.Weekday, "not", weekDay ) 
	    fails++ 
	    success = false
	}

	if slotDay.FlightsInSeason != flightsInSeason  { 
	    fmt.Println ( "FAIL: Flights in Season", slotDay.FlightsInSeason , flightsInSeason  ) 
	    fails++ 
	}

	dayMax := float64 ( slotDay.CurrentMax )
		    
	if dayMax != currentMax { 

	    fmt.Println ( "FAIL: Curent Max Cancles", dayMax , currentMax ) 

	    fails++ 
	    success = false
	}

/*	Calculate the following based on Slot Rules
//	Original Max Cancels:  FLOOR ( 20% of Flights in Season )
//
	slotMax := int ( math.Floor( fs.SunCurrentMax + 0.5 ) )
	oMax := SetOriginalMax ( slotMax, fs.MonFlightsInSeason ) 
	if  slotDay.OriginalMax != oMax { 

	    fmt.Println ( "FAIL: Original Max Cancels", slotDay.OriginalMax , oMax ) 
	    fails++ 
	    success = false
	}

//	Need to Get Carrier Code / Flight Number ( for prefix )

//	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

//	For now, just plug one ( "AAA" )

	slotIndex := ComputeSlotIndex ( slotMax, oMax, "AAA"  )

	if slotDay.SlotIndex != slotIndex { 

	    fmt.Println ( "FAIL:  Slot Index", slotDay.SlotIndex , slotIndex ) 
	    fails++ 
	    success = false
	}
*/
	if success { return success }

//	For one or more reaons, this Slot ( Day ) Failed

//	Update Failure Counters

	dayFails[weekDay] = dayFails[weekDay] + 1	// Total Files by Day

	slotFails++		// # of Days failed for Slot File
	daysFailed++		// # of Days failed this week

//	Maintian # of Failures -- of any type -- for Slot File

	totalSlotFailures = totalSlotFailures + fails

//	Report Slot Failure

	fmt.Println ( weekDay, ": Failed", fails, "times" )
	fmt.Println ( "Slot Day:", slotDay )
	fmt.Println ( "Full Slot:", fs )

	return false
}


