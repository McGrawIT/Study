package slot_svc

import (
	"testing"
	"fmt"
	"strings"
	"strconv"
)

/*
type FilterParams struct {

	airportCodes		string
	carrierCodes		string
	flightNumbers		string
	operationDates		string	
}
*/

func TestParseQueryString( t *testing.T ) {

//	Intitialize, if needed

	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

//	Special-case initializations required for this specific _test.go file


	fmt.Println ( "--------------------------------------------------" )

	fmt.Println ( "Starting Testing of ParseQueryString()" )

//	Test Data

	airports := []string{ "DXB,MCO,LAX,JFK", }
	airportsString := "DXB,MCO,LAX,JFK"

	carriers := []string{ "FZ,DL,B6", }
	carriersString := "FZ,DL,B6"

	opDates := []string{ "2016-May-15", }
	opDatesString := "2016-May-15"

	flights := []string{ "662,15,266,7000", }
	flightsString := "662,15,266,7000"


//	Make sure thie is the right structure to pass in......

	queryParams := make ( map[string][]string )

	queryParams["airportCode"] = airports
	queryParams["carrierCode"] = carriers
	queryParams["operationDate"] = opDates
	queryParams["flightNumber"] = flights

//	parseQueryString ( queryParams map[string][]string ) ( FilterParams )

	filters := parseQueryString ( queryParams )

//	Validation

	queryStringFailures := 0

//	airports := []string{ "DXB,MCO,LAX,JFK", }

	includedAirports := strings.Split ( filters.airportCodes, "," ) 

	failMsg := "FAIL: parseQueryString()"

	for _, value := range includedAirports {

	    if value == "DXB" { continue }
	    if value == "MCO" { continue }
	    if value == "LAX" { continue }
	    if value == "JFK" { continue }

	    fmt.Println ( failMsg, "Airport Code =", value )
	    t.Fail()
	    queryStringFailures++
	    continue		// Check for more Airport Code failures
	}
	air := strings.Split( airportsString, "," )
	car := strings.Split( carriersString, "," )
	op := strings.Split( opDatesString, "," )
	flt := strings.Split( flightsString, "," )

	if len ( includedAirports ) != len ( air ) {

	    fmt.Println ( failMsg, "Airports:", air, "Got:", includedAirports )
	    t.Fail()
	    queryStringFailures++
	}

//	carriers := []string{ "FZ,DL,B6", }

	includedCarriers := strings.Split ( filters.carrierCodes, "," ) 

	for _, value := range includedCarriers {

	    if value == "FZ" { continue }
	    if value == "DL" { continue }
	    if value == "B6" { continue }

	    fmt.Println ( "FAIL: parseQueryString() Carrier Code =", value )
	    t.Fail()
	    queryStringFailures++
	    continue		// Check for more Carrier Code failures
	}

	if len ( includedCarriers ) != len ( car ) {

	    fmt.Println ( failMsg, "Carriers:", car, "Got:", includedCarriers )
	    t.Fail()
	    queryStringFailures++
	}

//	flights := []string{ "662,15,266,7000", }

	includedFlights := strings.Split ( filters.flightNumbers, "," )

	for _, value := range includedFlights {

	    if value == "662" { continue }
	    if value == "15" { continue }
	    if value == "266" { continue }
	    if value == "7000" { continue }

	    fmt.Println ( "FAIL: parseQueryString() Flight Number =", value )
	    t.Fail()
	    queryStringFailures++
	    continue		// Check for more Filght Number failures
	}

	if len ( includedFlights ) != len ( flt ) {

	    fmt.Println ( failMsg, "Flights:", flt, "Got:", includedFlights )
	    t.Fail()
	    queryStringFailures++
	}

//	opDates := []string{ "2016-May-15", }

	includedDates := strings.Split ( filters.operationDates, "," )
	includedDates = opDatesNormalize ( includedDates )

	for _, value := range includedDates {

	    if value == "2016-May-15" { continue }

	    fmt.Println ( "FAIL: parseQueryString() Op Date =", value )
	    t.Fail()
	    queryStringFailures++
	    continue		// Check for more Op Date failures
	}

	if len ( includedDates ) != len ( op ) {

	    fmt.Println ( failMsg, "Op Dates:", op, "Got:", includedDates )
	    t.Fail()
	    queryStringFailures++
	}

	if queryStringFailures > 0 { 
	    fmt.Println ( "FAILED ", strconv.Itoa( queryStringFailures ), "time(s)" )
	} else { fmt.Println ( "PASSED: All Filters" ) }
	
}
