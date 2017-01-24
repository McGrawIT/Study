package slot_svc

import (
	"testing"
	"fmt"
)

/*	Used by most of the Date Functions

type	slotDate	struct {

	Year		int
	Month		int
	Day			int
	YMD			string		// YY-MMM-DD
	LongForm	string		// YYYY-MMM-DD
}

*/

func TestInSeason ( t *testing.T ) {

//	Initialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Starting Testing of InSeason ()" )

//	func inSeason ( operationDate, seasonBegin, seasonEnd string ) ( bool )

//	Test Season ( Begin / End Same Year )

	seasonBegin := "2016-Jan-31"
	seasonEnd := "2016-Dec-31"

	operationDate := "2016-May-31"		// Date IN Season

	in := inSeason ( operationDate, seasonBegin, seasonEnd )

//	Validate

	if !in { 
	    fmt.Println ( "FAILED:", operationDate,"NOT InSeason() [", seasonBegin, seasonEnd, "]" )
	    t.Fail() 
	} else {
	    fmt.Println ( "PASSED:", operationDate, "InSeason() [", seasonBegin, seasonEnd, "]" )
	}

//	Test / Validate Again

	operationDate = "2017-May-31" 		// Date NOT in Current Season

	in = inSeason ( operationDate, seasonBegin, seasonEnd )

	if in { 
	    fmt.Println ( "FAILED:", operationDate,"InSeason() [", seasonBegin, seasonEnd, "]" )
	    t.Fail() 
	} else {
	    fmt.Println ( "PASSED:", operationDate, "NOT InSeason() [", seasonBegin, seasonEnd, "]" )
	}

//	Season Spans Years  ( Test / Validate Again )

	seasonBegin = "2015-Oct-31"
	seasonEnd = "2016-Mar-31"

	operationDate = "2015-Dec-01" 	// Last Year Date in Current Season

	in = inSeason ( operationDate, seasonBegin, seasonEnd )

	if !in { 
	    fmt.Println ( "FAILED:", operationDate,"NOT InSeason() [", seasonBegin, seasonEnd, "]" )
	    t.Fail() 
	} else {
	    fmt.Println ( "PASSED:", operationDate, "InSeason() [", seasonBegin, seasonEnd, "]" )
	}
}

func TestDayOfWeek ( t *testing.T ) {

//	Initialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Starting Testing of DayOfWeek ()" )

//	Set Test Dates

	date1 := "2016-Jun-26"
	date2 := "2016-Jun-27"
	date3 := "2016-Jun-28"
	date4 := "2016-Jun-29"
	date5 := "2016-Jun-30"
	date6 := "2016-Jul-01"
	date7 := "2016-Jul-02"

//	func dayOfWeek (date string) ( string, int )

//	Call / Validate

	tests := 0
	fails := 0

	weekDay, dayNumber := dayOfWeek ( date1 )
	tests++
	if weekDay != "Sunday" && dayNumber != 1 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date2 )
	tests++
	if weekDay != "Monday" && dayNumber != 2 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date3 )
	tests++
	if weekDay != "Tuesday" && dayNumber != 3 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date4 )
	tests++
	if weekDay != "Wednesday" && dayNumber != 4 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date5 )
	tests++
	if weekDay != "Thursday" && dayNumber != 5 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date6 )
	tests++
	if weekDay != "Friday" && dayNumber != 6 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}

	weekDay, dayNumber = dayOfWeek ( date7 )
	tests++
	if weekDay != "Saturday" && dayNumber != 7 { 
	    fmt.Println ( "Unexpected Day Number (1-7) / Day of Week:", dayNumber, weekDay )
	    fails++
	    t.Fail() 
	}
	fmt.Println ( "Completed", tests, "Day of Week Tests:", fails, "Failed" )
}
