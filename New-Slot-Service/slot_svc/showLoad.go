package slot_svc
//package main

import (
	"fmt"
)

func ShowLoad (source string) {

	fmt.Println ( "[[[[ ShowLoad(", source, ")]]]]" )

	
	fmt.Println ( "=======[[[[ Slot DB:", len(TestSlotDataDB ), "Slot Files ]]]]" )

	for k1, v1 := range TestSlotDataDB { 
		fmt.Println ( "Key:", k1 )
		slotsToShow := 1
		for _, v2 := range v1 {
			fmt.Println ( "Flight:", v2.flightNumber )
			fmt.Println ( "Slot Week:", v2.slotWeek )
			slotsToShow++
			if slotsToShow > 2 { break }
		}
	}
}

