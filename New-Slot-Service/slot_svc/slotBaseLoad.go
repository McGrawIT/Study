package slot_svc

import (
	"fmt"
	"os"
)

/*	Function creates, initializes and populates Slot Data File records
	to create a varied DB for testing the Filter criteria.

	This always creates an iniitial Slot File DB @ Slot Service start
	When Slot DB is loaded @ start, this may be removed, but it does
	help to keep around for remote / from home developemnt / testing
*/
/*=======================================================================*
		BASE LOAD ( Support >1 Value for Query Parameters )
=========================================================================*/

//func loadStubData() ([]slots, map[SlotKey][]FlightSlot) {
func loadStubData( f map[SlotKey][]FlightSlot ) {

	initialLoad = true

//	f := make(map[SlotKey][]FlightSlot)

//	Create a Slot Data File and its Slot Rows

//	Key for each stubbed Slot Data File ( Loop for >1 Slot Data File )

	airportCodes := []string{}
	carrierCodes := []string{}
	flightNumbers := []string{}

//	Default to Very Small Test Set

	veryLarge, large, medium, small := false, false, false, false
	verySmall := false

//	Select Test Data, if Environment is not set, use

	stubSize := os.Getenv ( "SLOT_STUB" )

	if stubSize == "" || ( port != "2525" && port != "3535" ) {	// On Predix

	    airportCodes = []string{"MCO","LAX",}
	    carrierCodes = []string{"DL",}
	    flightNumbers = []string{"20",}
	
	} else {

	    if stubSize == "VL" { veryLarge = true }
	    if stubSize == "L" { large = true }
	    if stubSize == "M" { medium = true }
	    if stubSize == "S" { small = true }
	    if stubSize == "VS" { verySmall = true }
	}

	medium = true		// Test @ Scale -- On Predix

	if veryLarge {

	    airportCodes  = []string{ "LAX", "SFO", "DFW", "MCO", "ATL",}
	    more := []string{ "LHR", "ORD", "CDG", "HND", "PEK",}
	    most := []string{ "AUS", "TNL","BOS","SEA","YYC",}
	    
	    for _, addOne := range more {
		airportCodes = append( airportCodes, addOne )
	    }
	    for _, addOne := range most {
		airportCodes = append( airportCodes, addOne )
	    }

	    carrierCodes = []string{ "DL", "AF", "BA", "AA", "B6", "VJ", "ZZ",}

	    flightNumbers = []string{"020", "120", "721", "070", "10", "3420",}
	    more = []string{"520", "59", "521", "3", "66", "643", "88", "744",}
	    alot := []string{"67", "00", "4", "7100", "12", "3", "7", "18",}
	    most = []string{"1816", "4880", "7758", "888", "98", "9", "1111", }

	    for _, addOne := range more {
		flightNumbers = append(flightNumbers, addOne)
	    }
	    for _, addOne := range alot {
		flightNumbers = append(flightNumbers, addOne)
	    }
	    for _, addOne := range most {
		flightNumbers = append(flightNumbers, addOne)
	    }

	} else if large {

	    airportCodes = []string{ "LAX", "SFO", "DFW", "MCO", "ATL"}
	    carrierCodes = []string{ "DL", "AF", "BA", "B6"}
	    flightNumbers = []string{"020", "120", "721", "10", "520", "521"}

	    more := []string{ "93", "0", "116", "42", "7750", "7100", "12", "7"}
	    for _, addOne := range more {
		flightNumbers = append(flightNumbers, addOne)
	    }

	} else if medium {

	    airportCodes = []string{"DXB", "LAX", "SFO", "MCO"}
	    carrierCodes = []string{"FZ", "DL", "B6"}
	    flightNumbers = []string{"20", "12",  "1", "52", "521", "930", "7"}

	} else if small {

	   airportCodes = []string{"MCO", "TNL"}
	   carrierCodes = []string{"DL", "FZ"}
	   flightNumbers = []string{"020", "018", "120"}

	} else if verySmall {

	    airportCodes = []string{"MCO"}

	    carrierCodes = []string{"DL"}
	    flightNumbers = []string{"20", "18", "120"}
	}

	weekDays := []string{Sun, Mon, Tue, Wed, Thu, Fri, Sat}

	s := SlotKey{}

	flightSlice := []FlightSlot{}

	seasons := 0

	keyBeginDate := "2014-Apr-05"
	keyEndDate := "2015-Jun-04"

//	Number of Seasons of Test Data

	for i := 1; i <= 6; i++ {

	    switch i {
	    case 1:

		keyBeginDate = "2012-Jan-11"
		keyEndDate = "2013-Aug-06"

	    case 2:

		keyBeginDate = "2013-Sep-30"
		keyEndDate = "2014-Jul-04"

	    case 3:

		keyBeginDate = "2014-Sep-11"
		keyEndDate = "2015-Feb-01"

	    case 4:

		keyBeginDate = "2017-Jan-25"
		keyEndDate = "2017-Nov-30"

	    case 5:

		keyBeginDate = "2018-Apr-15"
		keyEndDate = "2026-Dec-24"

	    case 6:

		keyBeginDate = "2028-Jan-11"
		keyEndDate = "2050-Aug-06"
	    }
	    seasons++

//	    Slot Data File Key [ Airport, Carrier, Season ]

//	    Add Season Dates to Key

	    s.SeasonStart = keyBeginDate
	    s.SeasonEnd = keyEndDate

	    for a, airportCode := range airportCodes {

//		Add Airport Code to Key

		s.Airport = airportCode

		for c, carrierCode := range carrierCodes {

//		    Add Carrier Code ( Key complete )
//		    Create Slot File Index

		    s.Carrier = carrierCode

		    slot1 := s 			// Slot File Index

		    f[slot1] = flightSlice 	// Initialize (no Flight Slots)

//		    Do the following for each Flight ( and the Slots )

		    for fnum, flightNumber := range flightNumbers {

			flight := flightNumber

			ingestFlightNumber := flight

//			Create the Slot Week for this Flight For each Day

			ingestSlotDay := SlotRecord{}
			SlotWeek := []SlotRecord{}

			varySlotInex := 0
			fis := 0
			cMax := 0

			for _, dayName := range weekDays {

			    ingestSlotDay.Weekday = dayName

			    varySlotInex++

			    if varySlotInex == 1 { fis = 20 }
			    if varySlotInex == 2 { fis = 25 }
			    if varySlotInex == 3 { fis = 31 }

			    ingestSlotDay.FlightsInSeason = fis

			    if varySlotInex == 1 { cMax = 4 }
			    if varySlotInex == 2 { cMax = 3 }
			    if varySlotInex == 3 { 
				cMax = 2 
				varySlotInex = 0
			    }

			    ingestSlotDay.CurrentMax = cMax

			    fNum := ingestFlightNumber
			    fNum = carrierCode + fNum

			    oMax := SetOriginalMax( cMax, fis )

			    ingestSlotDay.OriginalMax = oMax

			    slotIndex := ComputeSlotIndex( cMax, oMax, fNum )

			    ingestSlotDay.SlotIndex = slotIndex

//			    Append each Day to the Week

			    SlotWeek = append( SlotWeek, ingestSlotDay )
			}

/*			After a full week, create the Flight Slot record for
			the Slot Data File ( Flights w/i the Slot Data File 
			will have same Slots / Season
*/
//			Flight Slot

			f1 := FlightSlot{ingestFlightNumber, SlotWeek}

//			Add Flight Slot to Slot File ( slot1 )
//			This adds to the Slot File in the Slot File DB ( f )

			if seasons == 2 && a == 3 && c == 4 { continue }
			if seasons == 1 && c == 3 { continue }
			if seasons == 5 && a == 2 { continue }
			if seasons == 3 && fnum == 6 { continue }

			nogo := fnum + c + a

			if nogo > 12 { continue }

			f[slot1] = append( f[slot1], f1 )
		    }
		}
	    }
	}

	fmt.Println("Base Data:", len(f), "Slot Files (", seasons, "Seasons )")

	fmt.Println("Airports:", len(airportCodes))
	fmt.Println("Carriers:", len(carrierCodes))
	fmt.Println("Flight Numbers:", len(flightNumbers))

//	Use ssFilter() to return the slots{} rercords by a "Get All" search

	z := FilterALL(f)

	totalFlightSlots := len(z) / 7
	fmt.Println( "Flight Slots (Weeks) :", totalFlightSlots )
	fmt.Println( "Slots (Days): (", len(z))

	initialLoad = false

	return
}
