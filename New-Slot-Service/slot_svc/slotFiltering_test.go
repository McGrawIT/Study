package slot_svc

import (
	"fmt"
//	"strings"
	"testing"
)


/*	Filter (including ALL) responses will most likely be
	returned as traditional Rows (one Weekday / row)
*/
	
/*	Comma-delimited strings of 0-n values for each Paramter that are used 
	with OR logic.  (Aiport LAX or MCO, Airline DL or FZ, ...)

	Parameters themselves are connected with AND logic.  (LAX and DL)

	The combination of "any of these" for each parameter offers greater
	search options.  (Flight Number, not so much.)
*/

var	TestingInitialized = false

func TestSlotFiltering ( t *testing.T ) {

/*	All Test Functions will check this, and the first one will load 
	the Slot File DB and other variablee in support of Unit Tests
	The data is known, which creates expected results for each
	function under test.
*/
	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Starting Testing of filterSlots()" )

//	Init function must return a test Slot File DB
//	Fake it for initial testing

//	Input: 		map[SlotKey][]FlightSlot
//	Input: 		filters filterParams
//	Returns:	[]slots

	goodResultSets := 0
	badResultSets := 0

//	Slots Returned

	testFilters := setSlotFilters ( "DXB", "FZ", "", "" ) 
	testSlots := filterSlots ( TestSlotDataDB, testFilters )

	if len( testSlots ) != 7994 {

	    fmt.Println ( "# Slots Returned (FAIL)" )
	    t.Fail()
	    badResultSets++

	} else { goodResultSets++ }


	testFilters = setSlotFilters ( "", "", "8254", "" ) // 7
	testSlots = filterSlots ( TestSlotDataDB, testFilters )

	if len( testSlots ) != 7 {

	    fmt.Println ( "# Slots Returned (FAIL)" )
	    t.Fail()
	    badResultSets++

	} else { goodResultSets++ }

	testFilters = setSlotFilters ( "", "", "", "2016-06-30" ) // 792
	testSlots = filterSlots ( TestSlotDataDB, testFilters )

	if len( testSlots ) != 792 {

	    fmt.Println ( "# Slots Returned (FAIL)" )
	    t.Fail()
	    badResultSets++

	} else { goodResultSets++ }

//	Operation Dates

//	One Flight on a Date:  1 / Thursday / 2016-06-30

	testFilters = setSlotFilters ( "", "", "662", "2016-06-30" ) 
	testSlots = filterSlots ( TestSlotDataDB, testFilters )

	good := true
	if len( testSlots ) != 1 { good = false }

	for _, s := range testSlots {

	    if s.opDate != "2016-Jun-30" { good = false } 
	    if s.weekday != "Thursday" { good = false } 
	    if s.flightNumber != "662" { good = false } 

	    if !good {

			fmt.Println ( "FAILED:", "Single Op Date" )
			fmt.Println ( "Expeceted: 662 on 2016-06-30 (Thursday)" )
			fmt.Println ( "Got:", s )
			fmt.Println ( "# of Slots:", len( testSlots ) )
			t.Fail()
			badResultSets++

		} else { goodResultSets++ }
	}

//	All Flights on an Op Date:  792 / Wednesday / 2016-06-29

	testFilters = setSlotFilters ( "", "", "", "2016-06-29" ) 
	testSlots = filterSlots ( TestSlotDataDB, testFilters )

	if len( testSlots ) != 792 { 

	    fmt.Println ( "FAILED:", "All Flights on Op Date" )
	    fmt.Println ( "Expected: 792, Got:", len( testSlots ) )
	    t.Fail() 
	    badResultSets++

	} else { goodResultSets++ }

	for _, s := range testSlots { 

	    if s.opDate != "2016-Jun-29" {

			fmt.Println ( "FAILED:", "792 Flights One Op Date" )
			t.Fail()
			badResultSets++
			break

	    } else { goodResultSets++ }

	    if s.weekday != "Wednesday" {

					fmt.Println ( "FAILED:", "792 Flights One Weekday" )
			t.Fail()
			badResultSets++
			break

		} else { goodResultSets++ }
	}

/*	Six (6) Flights: 662, 8254, 434, 7380, 8241, 881
	On two (2) Op Dates:  2016-09-30 (Friday), 2015-12-31 (Thursday)
	Result:  Nine Slots ( Six on Friday / Three on Thursday )
	Found in two Seasons ( Summer 2016, Winter 2015 )
*/

	flightList := "662,8254,434,7380,8241,881"
	opDateList := "2016-09-30,2015-12-31"
	testFilters = setSlotFilters ( "", "", flightList, opDateList )
	testSlots = filterSlots ( TestSlotDataDB, testFilters )

	if len( testSlots ) != 9 {

		fmt.Println ( "FAILED:", "Complex (not 9)" )
		t.Fail()
		badResultSets++

	} else { goodResultSets++ }

	for _, s := range testSlots { 

	    if s.opDate != "2016-Sep-30" && s.opDate != "2015-Dec-31" {

			fmt.Println ( "FAILED:", "Complex (not Op Date)" )
			fmt.Println ( "Slot:", s )
			t.Fail()
			badResultSets++

	    } else { goodResultSets++ }

	    if s.weekday != "Friday" && s.weekday != "Thursday" {

			fmt.Println ( "FAILED:", "Complex (not Weekday)" )
			t.Fail()
			badResultSets++

		} else { goodResultSets++ }

	    break
	}

//	No Filters, should return ENTIRE Slot File DB ( 8.022 Slots )

	testSlots = FilterALL ( TestSlotDataDB )

//	Added exclude filter for Base Load's test data, need to decrement
//	the expected results to match ( valid change to Test Slot DB )

	adjustedLen := len( testSlots ) + 28	// Adjust for Valid change

	if adjustedLen != 8022 {

	    fmt.Println ( "FAILED (No Filters) Expected 8022 slots got:", len( testSlots ), "(expecting difference of 28)" )
	    t.Fail() 
	    badResultSets++

	} else { goodResultSets++ }
	

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Ending Testing of filterSlots()" )

	if badResultSets == 0 { fmt.Println ( "filterSlots() Passed ALL", goodResultSets, "Test Cases" )
	} else { fmt.Println ( "filterSlots() Passed", goodResultSets, "and Failed", badResultSets ) }
}

func setSlotFilters ( ac, cc, fn, od string ) ( FilterParams ) {

	testFilters := FilterParams{}

	testFilters.airportCodes = ac
	testFilters.carrierCodes = cc
	testFilters.flightNumbers = fn
	testFilters.operationDates = od

	return testFilters
}
