package slot_svc

import (
	"fmt"
//	"os"
//	"io"
	"encoding/json"
)

var	ReloadedSlotDB		map[SlotKey][]FlightSlot

func ReloadSlotDB ( slotDBfromBlobStore []byte ) ( bool ) {

//	Read the JSON into a slice of denormalized Flight Slots.
//	There are one to many Slot Files, each with their Flight Slots

        n := len ( slotDBfromBlobStore )

	if n == 0 { 
	    fmt.Println ( "No Slot DB in Blob Store" )
	    return false

	} else { fmt.Println ( "Reloading", n, "bytes of SlotDB from Blob Store" ) }

//	Load entire Slot DB into slice of "flattened" partial Flight Slots

	slotDBreloaded := []slotPartFlat{}
	err := json.Unmarshal ( slotDBfromBlobStore[:n], &slotDBreloaded )

	if err != nil { 
		fmt.Println ( "json.Unmarshal() into Slots Failed:", err )
		return false
	}
	fmt.Println ( "Reloading", len ( slotDBreloaded ), "Slots" )

//	Prepare to load in-Memory Slot DB from SlotDB ( in Blob Store )
//	This will load the entire Slot File DB ( no in-memory DB exists 
//	any time Slot Service starts ( initial boot and restarts )

	f := make(map[SlotKey][]FlightSlot)

	s := SlotKey{}

	s.SeasonStart = ""
	s.SeasonEnd = ""
	s.Airport = ""
	s.Carrier = ""

	slot1 := s

	ingestFlightNumber := ""

/*	Read every row in the flattened Slot DB, loading them into In-Memory 
	Slot DB.  Each time the Key changes, it is a New Slot File in the DB.
*/
	day := 0
	testFlightSlots := 0
	ingestSlotWeek := []SlotRecord{}		// Blank Week

	for _, slots := range slotDBreloaded {

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

	ReloadedSlotDB = f		// Restore / Reload Slot File DB

	return true
}
