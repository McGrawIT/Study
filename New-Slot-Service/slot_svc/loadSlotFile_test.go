package slot_svc

import (
	"testing"
	"fmt"
//	"os"
//	"io/ioutil"
//	"strconv"
//	"strings"
//	"net/http"
	"encoding/json"

)

//	Use following Flights for Validation of Slot File Load

var	validateFlights		= map [string]int{"18":1,"340":1,"25":1,"573":1,"8244":1,"22":1,"4512":1}

func TestLoadSlotFileJSON ( t *testing.T ) {

//	Intitialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}
	fmt.Println ( "--------------------------------------------------" )

	fmt.Println ( "Starting Testing of LoadSlotFileJSON()" )

/*
	LoadSlotFileJSON (f, slot1, slotRows, slient) (map[SlotKey][]FlightSlot)

	f 		map[SlotKey][]FlightSlot	// Slot File DB
	slot1 		SlotKey				// Slot File Key
	slotRows 	[]SlotJSON			// Slots to Load
	silent 		bool				// Logging?
*/

//	Required initializations / test setup

//	Simulate GET from Config Service

	contents := []byte{}
	contents, _ = readSlotJSON( "Slots.JSON" )

//	Place the JSON body ( simulated GET ) into the internal 
//	Slot structure

	blob := []SlotJSON {}

	err := json.Unmarshal ( contents, &blob )

//	If the Unmarshal fails, no ohter testing / valiation is possible

	if err != nil { 
	    fmt.Println ( "Unmarshal() error:", err ) 
	    t.Fail()
	    return
	}

//	Validate expected number of Flight Slots

	fmt.Println ( len(blob), "Filght Slots in test JSON" )

/*	Simulate call to GetMetadata() that retrieves the airportCode, 
	carrierCode, Season BeginDate, and Season EndDate 
*/
	airportCode := "MCO"
	carrierCode := "DL"
	keyBeginDate := "2016-May-31"
	keyEndDate := "2016-Dec-31"

//	Create Slot File DB Key ( Airport, Carrier, Season )

	s := SlotKey{}

	s.Airport = airportCode
	s.Carrier = carrierCode
	s.SeasonStart = keyBeginDate
	s.SeasonEnd = keyEndDate

	fmt.Println ( "Ready to Call LoadSlotFileJSON ( SlotDB, s, blob )" ) 
	fmt.Println ( "Slot File Key:", s )

//	Load Slot Data File into Slot Data File DB  

	fmt.Println ( "-------------loadSlotFileJSON () ------------------" )

	silent := false
	TestSlotDataDB = LoadSlotFileJSON ( TestSlotDataDB, s, blob, silent ) 

	fmt.Println ( "-----loadSlotFileJSON () [ Validation ]------------" )

//	Validation

//	Use the Slot File Key created above to locate the newly created
//	Slot File in the DB.  Examine the contents in the DB to validate

//	Total Flight Slots ( for now, same as # of Flights )

	fmt.Println ( "Slot File:", len ( TestSlotDataDB[s] ), "Flight Slots" )

	v2 := TestSlotDataDB[s]

	fail := "FAILED: Flight"
	day := "Slot Day:"

	failed := false
	testCases := 0
	fails := 0

	slotCount := len ( TestSlotDataDB[s] )
	if slotCount != 792 { 

	    fmt.Println ( "FAILED:  Slot File Load of", slotCount, "vs 792" )
	    fails++
	    t.Fail() 
	}

	for _, flightSlot := range v2 {

	    j := validateFlights[ flightSlot.flightNumber ]

	    if j != 1 { continue }

//	    This Flight is used for Validation

	    failed = false			// Reset for Each Flight

	    testCases++ 

	    switch  flightSlot.flightNumber {
	    case "18" : fmt.Println ( "Slot Day", flightSlot.slotWeek[0] )
	    case "340" : fmt.Println ( "Slot Day", flightSlot.slotWeek[1] )
	    case "25" : fmt.Println ( "Slot Day", flightSlot.slotWeek[2] )
	    case "573" : fmt.Println ( "Slot Day", flightSlot.slotWeek[3] )
	    case "8244" : fmt.Println ( "Slot Day", flightSlot.slotWeek[4] )
	    case "22" : fmt.Println ( "Slot Day", flightSlot.slotWeek[5] )
	    case "4512" :

		fmt.Println ( "Testing Slot Week:", flightSlot.slotWeek )
		if flightSlot.slotWeek[6].Weekday != "Saturday" { failed = true }
		if flightSlot.slotWeek[6].FlightsInSeason != 11 { failed = true }
		if flightSlot.slotWeek[6].OriginalMax != 2 { failed = true }
		if flightSlot.slotWeek[6].CurrentMax != 1 { failed = true }
		if flightSlot.slotWeek[6].SlotIndex != 0.5 { failed = true }
	    }

	    if failed {

		expected := "{Saturday  11 2 1 0.5}"
		fmt.Println ( fail, flightSlot.flightNumber, day, flightSlot.slotWeek[6] )
		fmt.Println ( "Expected Slot Day:", expected )

		fails++
		t.Fail()
continue
	    }
	}

	if fails == 0 {		// PASSED

	    fmt.Println ( "PASSED:", s, "Tests:", testCases )

	} else {		// FAILED

	    fmt.Println ( "FAILED:", s, "Fails:", fails, "Tests:", testCases )
	}

	fmt.Println ( "---------loadSlotFileJSON () [ END ]---------------" )
}


func TestSetOriginalMax ( t *testing.T ) {

//	Intitialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}
	fmt.Println ( "--------------------------------------------------" )

	fmt.Println ( "Starting Testing of SetOriginalMax()" )

// 	SetOriginalMax(currentMax, flightsInSeason int) int

//	Decision was made to NOT save the first Current Max as the Original 
//	Max Cancels, so first arugment ( for now ) is not relevant

//	Each response is simply the floor of 20% of Flights in Season

	fails := 0

	original := SetOriginalMax ( 0, 22 )
	if original != 4 { 
	    fails++
	    t.Fail() 
	}

	original = SetOriginalMax ( 0, 25 )
	if original != 5 { 
	    fails++
	    t.Fail() 
	}

	original = SetOriginalMax ( 0, 19 )
	if original != 3 { 
	    fails++
	    t.Fail() 
	}

	if fails ==  0 {
	    fmt.Println ( "PASSED: 3 SetOriginalMax() calculations" )
	    return
	}
	fmt.Println ( "FAILED: SetOriginalMax() calculation", fails, "times" )
}
