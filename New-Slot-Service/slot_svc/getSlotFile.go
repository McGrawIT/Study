package slot_svc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"

	au "github.build.ge.com/AviationRecovery/go-oauth.git"
)

var	addLocal = true		// Hack to include Slots.JSON on Predix

//	Season File ( Details about each Season, indexed by Season Title
//	Example:  Summer S16 ( see sample Filename below )

var	seasonList	map [ string ]seasonDates

type	seasonDates	struct {
	AirportCode	string		// DXB, MCO, LAX, ..
	CarrierCode	string		// FZ, DL, ...
	SeasonName	string		// WINTER, SUMMER, ...
	SlotFileName	string		// DXB Summer S16 (25Mar-20Oct16).xlsx
	Begin		string		// 2016-Jun-17
	End		string		// 2016-Jun-17
}

func getSlotDataFilePredix ( silent bool ) {

	contents := []byte{}

	if fromHome {		// For testing w/o Network Access

	    fmt.Println ( "From Home:, loading Slots.JSON" )
	    contents, _ = readSlotJSON( "Slots.JSON" )

	} else {

            fmt.Println ( "GET Current Slot File:", currentSlotFile )
            fmt.Println ( "Using:", slotFileURL ) 

//	    Request Authorization to Access Config Service resources

	    client := au.RequestAuthorization()

//	    Issue GET to Config Service for Available Slot File

	    response, err := client.Get( slotFileURL )

	    if err != nil {
		fmt.Println ( "Slot File GET Failed using client.GET:" )
		fmt.Println ( "GET Error: ", err.Error() )
		fmt.Println ( "Aborting attempt to get Slot File" )
		return
	    }

	    if response == nil {
		 
		fmt.Println ( "No GET error, but GET response is nil" )
		fmt.Println ( "Aborting attempt to get Slot File" )
		return
	    }

	    if response.StatusCode != http.StatusOK { 

		fmt.Println ( "No GET error, but http.StatusCode NOT 200" )
		fmt.Println("Http.Status:", response.StatusCode )
		fmt.Println ( "Aborting attempt to get Slot File" )
		return
	    }

//	    All OK, pull Slot File Contents from GET Body

	    defer response.Body.Close()

	    contents, err = ioutil.ReadAll(response.Body)
	    if err != nil { 

		fmt.Println ( "Unable to pull Slot Data from GET Body" )
		fmt.Println ( "Aborting attempt to get Slot File" )
		return
	    }

//	    Force Unique Load ... try to create Changes in a Season

	    if addLocal {
		contents = []byte ( slotsJSON )
		addLocal = false
	    }
	}

//	Success:  Slot Data pulled from Config Service endpoint

	fmt.Println ( "Contents:", len ( contents ), "Bytes" )

/*	Update in-memory Slot DB using contents of this Slot File

	If this is a totally new Slot File, create a new Slot File entry
	in Slot DB and add all Flight Slots 

	New Flight Slots ( Flight added to Season ) Slot File in Slot DB
	Check existing Flight Slots for any updates to their Slot Week
*/
	blob := []SlotJSON {}

	err := json.Unmarshal ( contents, &blob )

	if err != nil { 

	    fmt.Println ( "Unmarshal() Error:", err.Error() )
	    fmt.Println ( "Aborting attempt to get Slot File" )
	    return
	}

//	GET Slot File Metadata ( Hit the endpoint for the Season )

	AC, CC, Begin, End, Season, slotFileName := getMetadata( currentSlotFile, silent )

//	Slot DB Key [ Airport, Carrier, Season Dates ] is unique to Slot File

	s := SlotKey{}

	s.Airport = AC		// Airport Code ( DXB, LAX, ... )
	s.Carrier = CC		// Carrier Code ( FZ, DL, B6, ... )
	s.SeasonStart = Begin	// Season Begin Date
	s.SeasonEnd = End	// Season End Date

	fmt.Println ( "------------------------------------------------------" )
	fmt.Println ( "About to Load / Update Slot File:", s )
	fmt.Println ( "------------------------------------------------------" )

//	Save Season Name ( for filtering an entire Seasson )

	seasonRange := seasonDates{}

	seasonRange.AirportCode = AC
	seasonRange.CarrierCode = CC
	seasonRange.SeasonName = currentSlotFile 
	seasonRange.SlotFileName = slotFileName 
	seasonRange.Begin = Begin
	seasonRange.End = End

	seasonList[ Season ] = seasonRange

	fmt.Println ( "Season Detail:", seasonRange )

/*	Load Slot Data File into Slot Data File DB  ( Add / Update )
*/
	SlotDataDB = LoadSlotFileJSON ( SlotDataDB, s, blob, silent ) 

}
