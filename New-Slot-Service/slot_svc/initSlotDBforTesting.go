package slot_svc

import (
	"fmt"
	"os"
	"io"
	"encoding/json"
)

var	TestSlotDataDB	map[SlotKey][]FlightSlot

func InitTestReadSlotDB () ( bool ) {

	slotConfigHostName = defaultConfigHostName   



	homeFlag := os.Getenv("HOME")

	fromHome = false
	if homeFlag == "YES" { fromHome = true }


//	Open Slot File DB for Testing ( all validations are based on this DB )

	slotDB, err := os.Open ( "SlotFile_test.DB" )

	if err != nil { 
		fmt.Println ( "Error opening SlotFile_test.DB:", err )
		return false 
	}

//	Close on exit and check for its returned error

	defer func() {
		if err := slotDB.Close(); err != nil { panic(err) }
	}()

//	Put the input into the Test Version of the in-memoery Slot File DB
//	This starts with "creating" the DB

	sizeSlotDataDB := 20000000
	buf := make( []byte, sizeSlotDataDB )

    fo, err := os.Create("SlotFileReadTest.DB")
    if err != nil { panic(err) }

//	Once created, read the JSON into a slice of Slot Data File Records
//	There are one to many Slot Data Files, each with their Flight Slots


//	Can this be read in all @ once?

        n, err := slotDB.Read ( buf )
        if err != nil && err != io.EOF { return false }

	fmt.Println ( "Read", n, "bytes of SlotDataDB for Unit Testing" )
        if n == 0 { return true }

//	Load entire Slot File DB into Slice of "flattened" partial Flight Slots

	slotDBin := []slotPartFlat{}
	err = json.Unmarshal ( buf[:n], &slotDBin )

	if err != nil { fmt.Println ( "json.Unmarshal() into Slots Failed:", err ) }
	fmt.Println ( len ( slotDBin ), "Slots" )

//	Write a chunk ( for now, test is write to file as is )

        _, err = fo.Write ( buf[:n] )

	if err != nil { return false }

/*	Range over the "flattened" Flight Slot slice, and "unflatten" them
	into the in-memory DB that is keyed by [ Airport, Carrier, Season ]
*/

//	Prepare to load in-Memory Slot DB from SlotDataDB ( Persistent store )
//	This will load the entire Slot File DB (only done on a Service restart

	f := make(map[SlotKey][]FlightSlot)

	TestSlotDataDB = f		// Restore / Reload Slot File DB

	s := SlotKey{}

	s.SeasonStart = ""
	s.SeasonEnd = ""
	s.Airport = ""
	s.Carrier = ""

	slot1 := s

	ingestFlightNumber := ""

/*	Read every row in the flattened Slot DB, converting them into
	the In-Memory Slot DB.  Each time the Key changes, it is a
	New Slot File in the DB
*/
	day := 0
	testFlightSlots := 0
	ingestSlotWeek := []SlotRecord{}		// Blank Week

	for _, slots := range slotDBin {

	    if 	s.SeasonStart != slots.SeasonStart || 
		s.SeasonEnd != slots.SeasonEnd ||
		s.Airport != slots.AirportCode ||
		s.Carrier != slots.CarrierCode {

/*		New Slot File ( [ Airprt / Carrier / Season ] Key changed )

		Do Slot File initializations ( prepare for Rows of
		Flight Slots for the new Airport / Carrier / Season )
*/
//		Set New Slot Key

	 	s.SeasonStart = slots.SeasonStart 
		s.SeasonEnd = slots.SeasonEnd 
		s.Airport = slots.AirportCode 
		s.Carrier = slots.CarrierCode 

		slot1 = s

		flightSlice := []FlightSlot{}	// Empty Flight Slot
		f[slot1] = flightSlice		// Avoid "null map" errors

		fmt.Println ( "New Slot File:", slot1 )
		day = 0
	    }
//	    if k < 7 { fmt.Println ( slots ) }		// Testing Only

/*	    Read the next Seven Rows ( Slot Days ( Sunday to Saturday ) )
	    appending them to the Slot Week for the current Flight to
	    create the Flght Slot is a Flight Number and Slice of Weeks 
	    ( for now, the Slot Data File has a 1:1 relationship between a
	    Flight and a Slot Week, but the Slice supports one Flight to
	    many Slot Weeks for future flexibility )
*/
	    ingestSlotDay := SlotRecord{}		// Another Day

	    if day == 0 {
		ingestSlotWeek = []SlotRecord{}		// Blank Week
		ingestFlightNumber = slots.FlightNumber 
	    }
	    day++

	    ingestSlotDay.FlightsInSeason = slots.FlightsInSeason	
	    ingestSlotDay.Weekday = slots.Weekday	
	    ingestSlotDay.OpDate = slots.OpDate	

	    ingestSlotDay.SlotIndex = slots.SlotIndex

	    ingestSlotDay.OriginalMax = slots.OriginalMaxCancels
	    ingestSlotDay.CurrentMax = slots.CurrentMaxCancels

//	    Day Detail complete, add Day to Week

	    ingestSlotWeek = append ( ingestSlotWeek, ingestSlotDay )

	    if day == 7 { 	// Week Complete

//		Flight Slot is Flight Number and Slot Week

		f1 := FlightSlot{ ingestFlightNumber, ingestSlotWeek }

//		Add Flight Slot to Slot File ( map[slotfile key][]flight slots )

		testFlightSlots++		// Count of All Slots

		f[slot1] = append( f[slot1], f1 ) 

		day = 0
		ingestSlotWeek = []SlotRecord{}		// Blank Week
	    }
	}
	fmt.Println ( "Flight Slots:", testFlightSlots, "in Slot File DB" )

//	ShowLoad ("Read Slot DB" )

	return true
}

