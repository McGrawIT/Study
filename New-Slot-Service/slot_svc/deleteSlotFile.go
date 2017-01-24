package slot_svc

import (
	"fmt"
	"strings"
	"net/http"
)


func DeleteSlotFile ( response http.ResponseWriter, query *http.Request ) {

	fmt.Println ( "Hit Delete Slot File Endpoint, Deleting:" )

	queryParams := query.URL.Query() 	// Read Query String

	slotFilesToDelete := parseQueryString ( queryParams )

//	For now, only Delete "first" Slot File matching first Codes / # / Date

	ac := slotFilesToDelete.airportCodes
	cc := slotFilesToDelete.carrierCodes
	fn := slotFilesToDelete.flightNumbers
	od := slotFilesToDelete.operationDates		// Delete Season

	fmt.Println ( "Airport:", ac )
	fmt.Println ( "Carrier:", cc )
	fmt.Println ( "Flight:", fn )
	fmt.Println ( "Operation Date:", od )

	airports := strings.Split ( ac, "," )
	carriers := strings.Split ( cc, "," )
	flights := strings.Split ( fn, "," )
	opDates := strings.Split ( od, "," )

	deleteAirport := ""
	deleteCarrier := ""
	deleteFlight := ""
	deleteOpDate := ""

//	Just delete the first of each Filter ( for now )

	if len ( airports ) != 0 { deleteAirport = airports[0] }
	if len ( carriers ) != 0 { deleteCarrier = carriers[0] }
	if len ( flights ) != 0 { deleteFlight = flights[0] }
	if len ( opDates ) != 0 { deleteOpDate = opDates[0] }

	deleteFlight = deleteFlight 
	deleteOpDate = deleteOpDate 

//	Find Season for Operation Date

//	==========>>> Need to Get findSeason() working first

//	season := FindSeason ( deleteOpDate )

//	Set Slot File Key ( If any Key part is nil no Slot File will be found )

	deleteKey := SlotKey{}

	deleteKey.Airport = deleteAirport
	deleteKey.Carrier = deleteCarrier

//	==========>> saason.Begin ( and .End ) must work with findSeason()

//	deleteKey.SeasonStart = season.Begin
//	deleteKey.SeasonEnd = season.End

	slotFileDeeleting := SlotDataDB [ deleteKey ]

	if len ( slotFileDeeleting ) == 0 { 

	    fmt.Println ( "No Slot File to Delete:", deleteKey )
	}

//	If Flight Number is included, locate Flight Slot to Delete


	for slotFileKey, _ := range SlotDataDB {

// 	    No Airport, this Slot File is delete cancidate
//	    If Airport Code included, skip if not a match

	    if deleteAirport != "" { 	

		if slotFileKey.Airport != deleteAirport { 
		    fmt.Println ( "Slot File", slotFileKey, "not Deleted" )
		    continue 
		}
	    }

// 	    No Carriers, this Slot File is delete cancidate
//	    If Carrier Code included, skip if not a match

	    if deleteCarrier != "" { 

		if slotFileKey.Carrier != deleteCarrier { 
		    fmt.Println ( "Slot File", slotFileKey, "not Deleted" )
		    continue 
		}
	    }
	    
//	    Slot File still candidate for Deletion
//	    Is Op Date in this Season?

	    fmt.Println ( "No Season Check for;", deleteOpDate  )

	    for _, oneFlightSlot := range SlotDataDB [ slotFileKey ] {

		if oneFlightSlot.flightNumber == deleteFlight {

		    fmt.Println ( "Removing Flight Slot:", deleteFlight )

//		    Make one pass over Flight Slots for this Slot File
//		    Reappend Flight Slots, skipping Flight for Deletion

		    fmt.Println ( "Flight Slot Removed from []flightSlots" )

		}

	    }

	    fmt.Println ( "Deleting Slot File:", slotFileKey )

//	    delete ( SlotDataDB, slotFileKey )

	}


}

