package slot_svc

import (
	"fmt"
)


//	Find the Flight Slot ( if it does not exist, just add )

var	debugMatching = false
var	flightSlotsChanged = 0

/*	f	Slot DB
	slot1	Key to Latest Slot File in Slot DB ( Existing, being Updated )
	f1	One Flight Slot from Latest Slot File ( from Config Svc )
*/

func updateFlightSlot ( slotDB map[SlotKey][]FlightSlot, slot1 SlotKey, f1 FlightSlot ) ( bool, bool ) {

	FlightSlotChanged := false	// True if Flight Slot is Updated

	slotFile := slotDB[slot1]	// Slot File being Updated

	currentFlightSlot := slotFileOffset [ f1.flightNumber ] 

	flightSlot := slotFile [ currentFlightSlot ]	

/*	Check to see if this Flight ( from the latest Slot File ) exists in the
	corresponding Slot File in the Slot DB.  ( We already determined the
	latest Slot File is in the Slot DB ( not new ).  We are just handling 
	any Flights added to the Season. ) 
	Slot File 
*/
	if flightSlot.flightNumber != f1.flightNumber {		// New Flight

	    fmt.Println ( "===> Added New Flight:", f1.flightNumber )
	    fmt.Println ( "Flight Slot:", f1.slotWeek )

//	    Add Flight Slot to Slot File
	
	    slotDB[slot1] = append ( slotDB[slot1], f1 )

	    return false, true		// Changed ( No ), Added ( Yes )
	}

//	Matching Flights Slots ( same Flight ) found

/*	Compare Slot Weeks

	flightSlot.flightNumber ( current )	f1.flightNumber	( latest )
	flightSlot.slotWeek ( current )		f1.slotWeek	( latest )

	slot1	( Current Slot File ( Key to File in Slot DB ) )
*/
	flightSlotCompares++

	changed := updateSlots ( &flightSlot, f1, slot1 )

	if changed { 

	    fmt.Println ( "Flight Slot Updated:", flightSlot )
	    flightSlotsChanged++ 
	    FlightSlotChanged = true
	}

	debugMatching = false

// 	Changed ( true / false ), Added ( false )

	return FlightSlotChanged, false	
}

/*	This is called when a Flight Slot matches.  Check each WeekDay.  
	Update any changed information.  Current Max Cancels is typically 
	the change.  If it does, recalculate the Slot Index.  
*/

func updateSlots ( flightSlot *FlightSlot, f1 FlightSlot, slot1 SlotKey ) (bool ) {

	slotChanged := false

	for day, v3 := range flightSlot.slotWeek {

	    slotCompares++

//	    Index into New Slot File to compare Day information

	    fs := f1.slotWeek[day]

/*	    Compare.  Weekday and Original Max Cancels do not change.
	    If Current Max Cancels, changes, calculate New Slot Index
*/
	    if fs.Weekday != v3.Weekday {
		fmt.Println ( "Out of Order:", fs.Weekday, v3.Weekday )
	    }

	    if debugMatching { }

	    updateSlotIndex := false

//	    If Fliights In Season Changed, we Create a New Oringal Max Cancels

	    if fs.FlightsInSeason != v3.FlightsInSeason {

		slotChanged = true

//		Update and Compute new Slot Index

		flightSlot.slotWeek[day].FlightsInSeason = fs.FlightsInSeason
		updateSlotIndex = true
	    }

//	    Most updates with be to Current Max Cancels ( if not all )

	    if fs.CurrentMax != v3.CurrentMax {
		updateSlotIndex = true

fmt.Println ( "" )
fmt.Println ( "----------------------------------------------------------" )
fmt.Println ( "Slot Changed, New:", fs, "Old:", v3 )
fmt.Println ( "----------------------------------------------------------" )
		
//		Update and Compute new Slot Index

		slotChanged = true

		flightSlot.slotWeek[day].CurrentMax = fs.CurrentMax

	    }

	    if updateSlotIndex {

		prefix := slot1.Carrier + f1.flightNumber
		oMax := SetOriginalMax ( fs.CurrentMax, fs.FlightsInSeason ) 

		index := ComputeSlotIndex ( fs.CurrentMax, oMax, prefix  )

		flightSlot.slotWeek[day].OriginalMax = oMax
		flightSlot.slotWeek[day].SlotIndex = index
	    }
	}
	return slotChanged
}
