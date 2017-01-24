package slot_svc

import (
	"testing"
	"fmt"
//	"os"
//	"io/ioutil"
//	"strconv"
//	"strings"
//	"net/http"
//	"encoding/json"

)
/*	Takes in FileName ( format: "AAA.....(DDMMM....DDMMMYY).....xlsx" )
	"Parses" and Returns:

	airportCode	// ( 3-digit Aiport Code (MCO, SFO, LAX, ... ) )
	carrierCode	// ( 2-digit Airline Code (DL, FZ, ... ) )
	SeasonBegin	// ( YY-MM-DD format )
	SeasonEnd	// ( YY-MM-DD format )
*/

func TestParseSlotFileName ( t *testing.T ) {

//	Intitialize, if needed

/*	Since Go Test does not guarantee order of execution, all _test.go files
	will check this variable to make sure the required initializations are 
	performed before any test (and happens only once / go test execution 
*/
	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

//	Special-case initializations required for this specific _test.go file

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Starting Testing of ParseSlotFileName ()" )

/*	Call function under test as many times as needed to at least test the
	primary aspects of this function for a thorough Unit Test.  
	Report / log errors with 
	    fmt.Println ( Date Failure:", begin )
	    t.Fail()
*/

//	Case One:  Begin Year same as End Year

	file := "DXB Summer S16 (27Mar-29Oct16)   .xlsx"

	fmt.Println ( "Slot File Name:", file )

	airport, carrier, begin, end := parseSlotFileName ( file )

	if airport != "DXB" { 
	    fmt.Println ( "airport Failure:", airport )
	    t.Fail() }
	if carrier != "FZ" { 
	    fmt.Println ( "carrier Failure:", carrier )
	    t.Fail() }

//	Workaround for the 20xx vs. xx Issue ... is it possible that a
//	4-digit year is just wrong??

	if begin != "16-Mar-27" { 
	    fmt.Println ( "Date Failure:", begin )
	    t.Fail() }
	if end != "16-Oct-29" { 
	    fmt.Println ( "Date Failure:", end )
	    t.Fail() }

/*
	if begin != "2016-Mar-27" { t.Fail() }
	if end != "2016-Oct-16" { t.Fail() }
*/

	fmt.Println ( "Result:", airport, carrier, begin, end )

//	Case Two:  Begin Date is Previous Year

	file =  "DXB Summer S16 (28Jun-20May16)   .xlsx"

	fmt.Println ( "Slot File Name:", file )

	airport, carrier, begin, end = parseSlotFileName ( file )

	if airport != "DXB" { 
	    fmt.Println ( "airport Failure:", airport )
	    t.Fail() }
	if carrier != "FZ" { 
	    fmt.Println ( "carrier Failure:", carrier )
	    t.Fail() }

//	Similar adjustment / validation against YY vs. YYYY

	if begin != "15-Jun-28" { 
	    fmt.Println ( "Date Failure:", begin )
	    t.Fail() }
	if end != "16-May-20" { 
	    fmt.Println ( "Date Failure:", end )
	    t.Fail() }

	fmt.Println ( "Result:", airport, carrier, begin, end )
}

//	MCO.....(12Dec....14Oct89).....xlsx

func TestParseFileName ( t *testing.T ) {

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

//	parseFileName ( filename, marker string, offset int ) 

	file := "DXB Summer S16 (27Mar-29Oct16)   .xlsx"

//	First Test Case:  offset after marker ( first date, ignore year )
//	Second Test Case:  negative offset to backup from marker

	bDay, bMonth, _ := parseFileName ( file, "(", 1 ) 
	eDay, eMonth, eYear := parseFileName ( file, ")", -7 ) 

//	Validate

	if bDay != "27" { 
	    fmt.Println ( "Date Failure:", bDay )
	    t.Fail() }
	if bMonth != "Mar" { 
	    fmt.Println ( "Date Failure:", bMonth )
	    t.Fail() }

	if eDay != "29" { 
	    fmt.Println ( "Date Failure:", eDay )
	    t.Fail() }
	if eMonth != "Oct" { 
	    fmt.Println ( "Date Failure:", eMonth )
	    t.Fail() }
	if eYear != "16" { 
	    fmt.Println ( "Date Failure:", eYear )
	    t.Fail() }

//	Repeat Test

	file =  "DXB Summer S16 (28Jun-20May16)   .xlsx"

	bDay, bMonth, _ = parseFileName ( file, "(", 1 ) 
	eDay, eMonth, eYear = parseFileName ( file, ")", -7 ) 

	if bDay != "28" { 
	    fmt.Println ( "Date Failure:", bDay )
	    t.Fail() }
	if bMonth != "Jun" { 
	    fmt.Println ( "Date Failure:", bMonth )
	    t.Fail() }

	if eDay != "20" { 
	    fmt.Println ( "Date Failure:", eDay )
	    t.Fail() }
	if eMonth != "May" { 
	    fmt.Println ( "Date Failure:", eMonth )
	    t.Fail() }
	if eYear != "16" { 
	    fmt.Println ( "Date Failure:", eYear )
	    t.Fail() }
}
