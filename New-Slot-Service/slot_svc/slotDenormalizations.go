
package slot_svc

import (
	"fmt"
	"net/http"
	"encoding/json"
)



/*
*/
func denormSlotsFlat ( z []slots, stub bool ) (string) {

//	Denormalizing (and Flattening for single Slot Data File) Result Set

	mySlot := slotPartFlat{}

	emptySlotPart := []slotPartFlat{}

	fmt.Println ( "About to Denorm", len(z), "Slots" )

	for _, slotRow := range z {

		mySlot.AirportCode 	= slotRow.airportCode	
		mySlot.CarrierCode 	= slotRow.carrierCode

		mySlot.SeasonStart	= slotRow.beginDate
		mySlot.SeasonEnd 	= slotRow.endDate

		mySlot.Weekday 		= slotRow.weekday 
		mySlot.OpDate 		= slotRow.opDate 

		mySlot.FlightNumber 	= slotRow.flightNumber
		mySlot.FlightsInSeason	= slotRow.flightsInSeason	

/*		Right now, the Slot Data File (in Excel) ingested has
		the Current Max Cancels initialized to 20% of Original
*/
		mySlot.CurrentMaxCancels = slotRow.currentMaxCancels
		mySlot.OriginalMaxCancels = slotRow.originalMaxCancels

		mySlot.SlotIndex	= slotRow.slotIndex

		emptySlotPart = append ( emptySlotPart, mySlot )
	}

	fmt.Println ( "About to Marsal to JSON" )

	result, err := json.MarshalIndent(emptySlotPart, "", "  ")

	fmt.Println ( "Marsal complete" )
	if (err != nil)  { fmt.Println ( "JSON Marshal Failed:", err ) }

	return string(result)

}

/*	This function does NOT remove the Slot Data File key (AC, CC, Season)
	and should be retained (for when more than just DXB / FZ are part
	of the production solution
*/

func denormSlots ( z []slots, response http.ResponseWriter, stub bool ) {

//	Denormalizing Result Set

	f := make(map[SlotKey][]slotPart) 

	key := SlotKey{}
	mySlot := slotPart{}

	emptySlotPart := []slotPart{}

	for _, slotRow := range z {

		key.Airport = slotRow.airportCode	
		key.Carrier = slotRow.carrierCode
		key.SeasonStart	= slotRow.beginDate
		key.SeasonEnd = slotRow.endDate

/*		If the map[key] does not exist, create an empty slice
		All records would then simply append
*/
		_, exists := f[key]

		if !exists {	// New map[key], create empty Map

			f[key] = emptySlotPart
		}

//		mySlot.Self = "Removed"
		mySlot.AirportCode = slotRow.airportCode	
		mySlot.CarrierCode = slotRow.carrierCode

		mySlot.FlightNumber = slotRow.flightNumber
		mySlot.FlightsInSeason	= slotRow.flightsInSeason	
		mySlot.OpDate		= slotRow.opDate	

/*		Right now, the Slot Data File (in Excel) ingested has
		the Current Max Cancels initialized to 20% of Original
*/
		mySlot.CurrentMaxCancels = slotRow.currentMaxCancels
		mySlot.OriginalMaxCancels = slotRow.originalMaxCancels

		mySlot.Weekday 		= slotRow.weekday 
		mySlot.SlotIndex	= slotRow.slotIndex

		f[key] = append ( f[key], mySlot )
	}

//	fmt.Println ( f )

/*	This is producing a bizarro error... complains that I can't use a
	map[SlotKey][]slotPart as type map[SlotKey][]slotPart in the following
	funtion.  Bizarre.

	UnpackSlots := unPackRealSlotforJSON ( f )

	Work around:  just put the function in-line (it's only called here)
*/

	unpackSlots := []ForPartJSON{}
	OneSlot := ForPartJSON{}

	for key, value := range f {

//		fmt.Println ( "key:", key, "value:", value )
		OneSlot.Key = key
		OneSlot.Slots = value

		unpackSlots = append ( unpackSlots, OneSlot )
	}
//	fmt.Println ( unpackSlots )

	UnpackSlots := unpackSlots

	result, _ := json.MarshalIndent(UnpackSlots, "", "  ")

	fmt.Fprintf ( response, string(result) )

//	fmt.Fprintf ( response, "Denormalizing Slots....\n" )

//	....Slots Denormalized.
}

