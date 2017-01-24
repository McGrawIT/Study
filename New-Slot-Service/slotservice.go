package main

import (

	"fmt"
	"net/http"

//	"github.build.ge.com/AviationRecovery/slot-service.git/slot_svc"
	"github.build.ge.com/502612370/New-Slot-Service/slot_svc"
)

const (

//	DB = "/bucket-c2cd2eb6-69e3-41c5-8ac7-276fe99ffbfc/?location= HTTP/1.1"
)

func main() {

//	----------------------------------------------------------------
//	Primary SlotData endpoint ( Queries ) ( Avoid Case Sensitivity )
//	----------------------------------------------------------------

	http.HandleFunc ("/api/v1/SlotData",	slot_svc.FullSearch )
	http.HandleFunc ("/api/v1/Slotdata",	slot_svc.FullSearch )
	http.HandleFunc ("/api/v1/slotData",	slot_svc.FullSearch )
	http.HandleFunc ("/api/v1/slotdata",	slot_svc.FullSearch )	// Official Slot Service ( Filtering )  endpoint

	http.HandleFunc ("/SlotData",			slot_svc.FullSearch )

//	---------------------------------
//	Provide Detailed Slot Information
//	---------------------------------

	http.HandleFunc ("/SlotInfo", 	slot_svc.SlotInfo )
	http.HandleFunc ("/", 			slot_svc.SlotInfo )

	http.HandleFunc ("/info", 		slot_svc.SlotInfo )				// Official "Info" endpoint
	http.HandleFunc ("/Info", 		slot_svc.SlotInfo )
/*
	----------------------------------------------------------------------------------
	CANCELLATIONS
	----------------------------------------------------------------------------------
	Added support for Cancellations; Slot Index is only valuable with remaining
	maximum Cancels.  Today, that is only provided in an Excel spreadsheet that is
	manually created / maintained by Fly Dubai.

	If Disruption Detection Service can hit the REST endpoint /Cancellation, then
	 Current Max Cancels can be maintained in real time and Slot Index updated as well
	----------------------------------------------------------------------------------

 */
	http.HandleFunc ( "/Cancellation", 	slot_svc.Cancellation )
	http.HandleFunc ( "/Delete", 		slot_svc.DeleteSlotFile )
/*
	----------------------------------------------------------------------------------
	Added Filtering by Slot Index.  It is a supported Query Parameter to /SlotData
	----------------------------------------------------------------------------------
 */
	http.HandleFunc ("/SlotIndex", 		slot_svc.FilterSlotIndex )

/*	----------------------
	Slot Data File polling
	----------------------
	Pull in most recent Slot Files, update Slot File DB ( new / changed )

	This endpoint is hit from a REST call internal to Slot Service
	We could consider a modification to Config Service to hit this endpoint when it
	receives a new Slot Data File.

	It's real-time data flow vs. polling or checking a file queue.

	Config Service already has logic to load Rabbit MQ, why not change it to have an environment variable
	with the endpoint vs. the queue.
 */

	http.HandleFunc ( "/api/v1/LoadSlotFile", 	slot_svc.SlotFileLoad )
	http.HandleFunc ( "/LoadSlotFile", 			slot_svc.SlotFileLoad )



	fmt.Println ( "----------------------------" )
	fmt.Println ( "Slot Service Server Starting" )
	fmt.Println ( "----------------------------" )

	slot_svc.SlotServiceServer ()
}
