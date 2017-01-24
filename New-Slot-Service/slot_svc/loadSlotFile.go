package slot_svc

import (
	"fmt"
	"math"
	"time"
)

/*	Load Flight Slots from latest Slot Data File into Slot File DB
	This may be New or Existing Slot File ( Add / Update )

	f		Slot File DB 
	slot1		Slot File DB Key ( airport / carrier / season )
	slowRows	Flight Slots for Slot File being loaded
*/
var	slotFileOffset 	map[string]int
var	latestSlotFile	map[SlotKey][]FlightSlot

var	flightSlotCompares 	int
var	slotCompares 		int

func LoadSlotFileJSON(f map[SlotKey][]FlightSlot, slot1 SlotKey, slotRows []SlotJSON, silent bool) map[SlotKey][]FlightSlot {

	fmt.Println ( "------------------------------------------------------" )
	fmt.Println ( "Load Slot File:", slot1, "Flight Slots:", len(slotRows) )

	flightSlice := []FlightSlot{}
	latestSlotFile = make(map[SlotKey][]FlightSlot)

	flightSlotChanges := 0

	latestSlotFile[slot1] = flightSlice

//	Add / Update Dates for this Slot Data File

	sfDetails := SlotDataFileDetails{}

	t := time.Now()
	d := t.Format("2006-Jan-01")
	ymmd := t.Format("2006-01-02")
	d = d[0:9] + ymmd[8:10]

//	New or Existing?  Add or Update

	isNew := false

	if len( f[slot1] ) == 0 {	// No Flight Slots; New Slot Data File

	    fmt.Println("New Slot File:", slot1)
	    fmt.Println("Created:", d)

	    f[slot1] = flightSlice

	    sfDetails.createdDate = d
	    sfDetails.lastUpdatedDate = d 	// Create is Latest Update
	    sfDetails.lastReceivedDate = d 	// Create is Latest File

	    fmt.Println( "New Slot File Details:", sfDetails )

	    isNew = true

	} else {

	    sfDetails = SlotFileInfo[slot1]
	    sfDetails.lastUpdatedDate = d	// Only Change if Updated
	    sfDetails.lastReceivedDate = d 	// Latest File

/*	    Before updating ( or adding ) any Flight Slots, create a hash map 
	    that allows for direct access to Flight Slot in Slot DB 			    ( vs. scanning Flight Slots ( in Slot DB ) to locate the 
	    Flight Slot to update )

	    f[slot1]	Slot File as it exists in Slot DB
	    slotRows 	Flight Slots from Latest Instance of Slot File
*/
	    slotFileOffset = make ( map[string]int )

	    i := 0
	    for _, fs := range f[slot1] {	// Flight Slots from Slot File

		slotFileOffset[ fs.flightNumber ] = i
		i++
	    }
	}

	fmt.Println( "Slot File", slot1, "Details:", sfDetails )

	SlotFileInfo[slot1] = sfDetails

	fmt.Println( "Slot File Info:", SlotFileInfo )
	fmt.Println( "Slot DB:", len ( f ), "Slot Files" )
	fmt.Println ( "------------------------------------------------------" )

/*
	Slot Data File entry initialized, ready to load Flight Slots
	Parameter "slotRows" contains Unmarshalled "SlotJSON" instances
	Input File is converted from Excel by Configuration Service
*/
//	Create the Slot Week for this Flight For each Day...

	FlightSlotsList := []string{}

	flightSlotCompares = 0
	slotCompares = 0

	fmt.Println ( "Processing", len ( slotRows ), "Slot Rows" )

	for _, fs := range slotRows { 	// slotRows are Flight Slots

//	    Load Slot Days for the Slot Week into the Flight Slice
//	    JSON input is extracted into SlotJSON structure

	    flight := fs.FlightNumber

	    ingestSlotWeek := BuildSlotWeek( slot1, fs )

/*	    All Seven Daya of the Week added to Slot Week

	    After a full week (Sun - Sat), create the Flight Slot
	    Flights w/i the Slot Data File will have same Slots / Season
*/
	    f1 := FlightSlot{ flight, ingestSlotWeek }

/*	    Add to Slot File DB, and also create a stand-alone instance
	    for Slot Data File updates ( Same Season can arrive with new
	    Slot information ( changed Current Max Cancels, New Flight Slot
	    are the most common ) )
*/
/*	    If this is a new Slot Data File ( f[slot1] did not exist )
	    then simply append the Flight Slot
*/
	    if isNew { 

		f[slot1] = append( f[slot1], f1 )

	    } else {

/*		This Slot File exists; locate the matching Flight Slot.
		Update as needed.  If a new flight was added to this Season,
		add the Flight Slot to the Slot File ( f[slot1] )
*/
		slotChanged, slotAdded := updateFlightSlot(f, slot1, f1)

		if slotChanged {

		    flightSlotChanges++
		    fmt.Println("Fight Slot Updated with Change(s)")
		    FlightSlotsList = append(FlightSlotsList, fs.FlightNumber)
		}

		if slotAdded {

		    fmt.Println( "Fight", flight, "added to Season schedule" )
		    fmt.Println( "Flight Slot:", f1 )
		}
		debugMatching = false
	    }
	    latestSlotFile[slot1] = append( latestSlotFile[slot1], f1 )
	}

	if isNew {
	    fmt.Println ( "Added Slot File", slot1, "to Slot File DB" )
	} else {

	    fmt.Println ( "Slot File compared to Slot DB:", flightSlotCompares, "Flight Slots and", slotCompares, "Slots" )

	}
	fmt.Println ( "------------------------------------------------------" )

/*	Only update the Slot DB when something changed ( updates to a Slot or
	a totally new Flight Slot ( new Flight for the Season )
*/
	fmt.Println("Slot File", slot1, "had", len( slotRows ), "Slot Rows to Add / Update" )
	fmt.Println("Slot File", slot1, "in DB has", len( f[slot1] ), "Flight Slots")

	if flightSlotChanges == 0 && !isNew {	// No changes, Not New

	    fmt.Println("Slot File", slot1,"had no new / changed Flight Slots" )

	    return f
	}

/*	New Slot File arrived or an existing Slot File was updated ( One or
	more Slots changed ). The in-memory Slot DB needs to be written to the
	Slot DB ( only time it does )
*/

	if isNew { 	// New Slot File ( and Flight Slots ) added to Slot DB

	    fmt.Println("Added Slot File", slot1, "to Slot DB" )

	    fmt.Println("New Flight Slots:", len( latestSlotFile[slot1] ))

	    z := FilterALL( latestSlotFile )	// Get all new Slots
	    fmt.Println("New Slots:", len(z) )

	    fmt.Println("Slot File DB now has", len(f), "Slot Files")

	} else {

	    fmt.Println("Made", flightSlotChanges, "Slot changes to Flight Slots[", FlightSlotsList, "] in Slot File", slot1 )
	}
	fmt.Println ( "------------------------------------------------------" )


/*	Updates to an existing Slot File ( or a new Slot File altogether ) is a
	changes to the Slot DB.  This is the only time the in-memory Slot DB 
	needs to be written to the Slot DB in Blob Store.  When real-time
	cancels are sent to Slot Service, the Slot DB will be saved more often.
*/
	allSlotsInDB := FilterALL(f)

	slotDataFileSaved := denormSlotsFlat( allSlotsInDB, false )

	written := writeSlotDB(slotDataFileSaved)

	if !written { fmt.Println("Error Saving Slot DB to Blob Store" ) }

	return f
}

var noMemory = true

func SetOriginalMax(currentMax, flightsInSeason int) int {

/*	Set Original Max every time a Slot Data File arrives, not
	just for the "true" original ( first file )
*/
	if noMemory {

	    oMax := float64(flightsInSeason)

	    realMax := oMax * 0.2     // Flights in Season * 20%
	    fl := math.Floor(realMax) // Round down

	    originalMaxCancels := int(fl)
	    return originalMaxCancels

	} else { return currentMax } 		// Do use "first File"
}
