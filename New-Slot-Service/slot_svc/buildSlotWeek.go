package slot_svc

import (
//	"fmt"
	"math"
)

func BuildSlotWeek ( slot1 SlotKey, fs SlotJSON ) ( []SlotRecord ) {

//	fmt.Println ( "buildSlotWeek() SlotKey:", slot1 )
	ingestSlotDay := SlotRecord{}
	ingestSlotWeek := []SlotRecord{}

	oMax := 0

//	A Slot Index rule:  Flights starting with FZ7 get a Slot Index of 1
//	Create a Combined Flight Number ( include Carrier ) for Prefix test

	prefix := slot1.Carrier + fs.FlightNumber

//	Sunday

	slotMax := int ( math.Floor( fs.SunCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Sun
	ingestSlotDay.FlightsInSeason = fs.SunFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax

/*	For now, we are NOT saving an "Original Max Cancels" but rather
	letting it be defined as FLOOR ( 20% * Flights In Season ), so it
	will change for every new Slot Data File.

	Initially, it was set / saved as the FIRST Current Max Cancels

		ingestSlotDay.OriginalMax = slotMax

	( This is easily changed back as the SetOrinalMax() has a flag )
*/
	oMax = SetOriginalMax ( slotMax, fs.SunFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex := ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Monday

	slotMax = int ( math.Floor( fs.MonCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Mon
	ingestSlotDay.FlightsInSeason = fs.MonFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.MonFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Tuesday

	slotMax = int ( math.Floor( fs.TueCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Tue
//	ingestSlotDay.FlightsInSeason, _ = strconv.Atoi ( fs.TueFlightsInSeason )
	ingestSlotDay.FlightsInSeason = fs.TueFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.TueFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Wednesday

	slotMax = int ( math.Floor( fs.WedCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Wed
	ingestSlotDay.FlightsInSeason = fs.WedFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.WedFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Thursday

	slotMax = int ( math.Floor( fs.ThuCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Thu
	ingestSlotDay.FlightsInSeason = fs.ThuFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.ThuFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Friday

	slotMax = int ( math.Floor( fs.FriCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Fri
	ingestSlotDay.FlightsInSeason = fs.FriFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.FriFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

//	Saturday

	slotMax = int ( math.Floor( fs.SatCurrentMax + 0.5 ) )

	ingestSlotDay.Weekday = Sat
	ingestSlotDay.FlightsInSeason = fs.SatFlightsInSeason
	ingestSlotDay.CurrentMax = slotMax
	ingestSlotDay.OriginalMax = slotMax

	oMax = SetOriginalMax ( slotMax, fs.SatFlightsInSeason ) 
	ingestSlotDay.OriginalMax = oMax
	slotIndex = ComputeSlotIndex ( slotMax, oMax, prefix  )

	ingestSlotDay.SlotIndex = slotIndex

//	Append Slot Day to Slot Week

	ingestSlotWeek = append( ingestSlotWeek, ingestSlotDay)

	return ingestSlotWeek
}
