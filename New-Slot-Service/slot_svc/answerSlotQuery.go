package slot_svc

import (
	"fmt"
	"net/http"
	"time"

	au "github.build.ge.com/AviationRecovery/go-oauth.git"

	"strconv"
)

/*
	/SlotData endpoint; any authentication before "Answering" the Query
 */
func FullSearch ( response http.ResponseWriter, query *http.Request ) {

/*
	When running from Home, the In-Memory DB will be loaded with Unit Test data.
	In this case, no Authorization needed
 */
	fmt.Println ( "Query:", query )

	if fromHome {

		fmt.Println ( "Skipped Authorization" )

	} else { 			//	Authenticate Search Request

		statusCode := au.CheckAuthentication( query )

		if statusCode != http.StatusOK {

			response.WriteHeader(statusCode)
			authFailure := "Authentication Status Code: " + strconv.Itoa(statusCode)
			fmt.Fprintf( response, authFailure )
			return
		}
		fmt.Println ( "Authorization Passed" )
	}

//	Authenticated Request, perform Query

	Stub = false
	answerSlotQuery ( response, query, false )
}

func answerSlotQuery ( response http.ResponseWriter, query *http.Request, stub bool ) {


	fmt.Println ( "-------------------------------------------" )

	fmt.Println ( "Blob URL:", DB2 )
	fmt.Println ( "Blob Content:", DBwritten )

/*
	Primary Slot Service logic

	1) Read Query String
	2) Parse Query String
	3) Filter Slots
	4) Return Search Results in JSON
 */
	queryParams := query.URL.Query() 								// Read Query String

	searchFilters := parseQueryString ( queryParams )				// Parse Query String

	slotResultSet := filterSlots ( SlotDataDB, searchFilters )		// Search ( Filter ) Slots
	fmt.Println ( "Found", len( slotResultSet ), "Slots" )

//	Filtering is Complete.  Creating JSON from Result Set...

	flatJSON := denormSlotsFlat ( slotResultSet, stub )

	fmt.Println ( "Loaded", len ( flatJSON ), "characters into JSON Ressponse ( ~300 / Slot )" )

	status := http.StatusOK											// Default ( Result Set non-Empty )

	if len( slotResultSet ) == 0 { status = http.StatusNoContent }  // Empty Result Set ( No Slots )

	response.WriteHeader( status )
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

	fmt.Fprintf ( response, flatJSON ) 								//	Write Response ( Result Set )

	writeTime := time.Now()
	fmt.Println ( "Query Response Sent:", writeTime )
	fmt.Println ( "-------------------------------------------" )
}