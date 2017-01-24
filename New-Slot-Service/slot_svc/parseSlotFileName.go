
package slot_svc

import (
	"fmt"
	"strings"
	"strconv"
)


var	mNumber = map[string]int{"Jan":1, "Feb":2, "Mar":3, "Apr":4, "May":5, "Jun":6, "Jul":7, "Aug":8, "Sep":9, "Oct":10, "Nov":11, "Dec":12}

/*	Takes in FileName ( format: "AAA.....(DDMMM....DDMMMYY).....xlsx" )
	"Parses" and Returns:

	airportCode		// ( 3-digit Aiport Code (MCO, SFO, LAX, ... ) )
	carrierCode		// ( 2-digit Airline Code (DL, FZ, ... ) )
	SlotKeyBeginDate,	// ( YY-MM-DD format )
	SlotKeyEndDat		// ( YY-MM-DD format )
*/

func parseSlotFileName ( filename string ) (string, string, string, string ) {

	fmt.Println ( "Now Parsing:", filename )

	airportCode := filename[0:3]
	carrierCode := "FZ"		// Not in Filename, Defaults

//	Begin Date does not have Year, ignore it

	bDay, bMonth, _ := parseFileName ( filename, "(", 1 ) 
	eDay, eMonth, eYear := parseFileName ( filename, ")", -7 ) 

//	Get integer Months, End Year

	start := mNumber[bMonth]
	end := mNumber[eMonth]

	endYear ,_ := strconv.Atoi (eYear)

//	Begin Year (calculated from Month / Month)

	beginYear := endYear 
	if start > end { beginYear = endYear - 1 }

//	Create YY-MM-DD strings (Begin Date has computed Year)

	SlotKeyEndDate := eYear + "-" + eMonth + "-" + eDay
	SlotKeyBeginDate := strconv.Itoa (beginYear) + "-" + bMonth + "-" + bDay

//	fmt.Println ( "Parsing:", filename )

//	Prints for Validations

//	fmt.Println ( "Yields:", airportCode, "Slot Dataa for", SlotKeyBeginDate, "-", SlotKeyEndDate )

	return airportCode, carrierCode, SlotKeyBeginDate, SlotKeyEndDate
}



//	MCO.....(12Dec....14Oct89).....xlsx

func parseFileName ( filename, marker string, offset int ) (string, string, string) {

	ds := strings.Index(filename, marker ) + offset
	ms := ds + 2
	ys := ms + 3
	
	day := filename[ds:ds+2]
	month := filename[ms:ms+3]
	year := filename[ys:ys+2]

	return day, month, year
}

