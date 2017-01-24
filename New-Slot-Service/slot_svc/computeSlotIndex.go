package slot_svc

import (
)

//	Calculate Slot Index

func ComputeSlotIndex ( current, original int, flight string ) ( float64 ) {

	originalMax := float64( original )
	currentMax := float64( current )
	slotIndex := 0.0

//	Total Cancels = Original Max Cancels - Current Max Cancels

	Cancels := originalMax - currentMax

/*	Slot Index = ( originalMax - currentMax ) / originalMax

	Unless...

	If Original Max Cancels = 0, Slot Index = 0% ( vs. Divide by Zero )
	If Cancels > Original Max Cancels, Slot Index = 100%
	If Carrier = FZ and Flight Number begins with "7", Slot Index = 100%
*/

/*	( 23-Jun-16, per Leif O. )  Added a check for Current Max Cancels
	greater than Original Max Cancels -- "bonus" alloable cancels?
*/
	if originalMax < currentMax { 
	    slotIndex = float64(0)
	    return slotIndex
	}

/*	For all Fly Dubai (FZ) flights from 7000 to 7999, set
	Slot Index to zero (0), regardless of Max Cancels.

	They want to always consider these flight for cancellation
*/
//	Caller will concatentate Carrier Code and Flight Number 

	if len ( flight ) == 6 {

		flightNumberPrefix := flight[0:3]

		if flightNumberPrefix == "FZ7" { 
			slotIndex = float64(0)
			return slotIndex
		}
	}

	if originalMax == 0 { return slotIndex }

	if Cancels > originalMax { 

		slotIndex = float64(1)
		return slotIndex
	}


//	No special cases, rules, or error checks

	slotIndex = Cancels / originalMax

	return slotIndex
}

