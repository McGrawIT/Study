package slot_svc

import (
	"fmt"
	"testing"
//	"os"
//	"io/ioutil"
//	"strconv"
//	"strings"
//	"net/http"
//	"encoding/json"
)

/*	JSON Body is placed in this structure.  Use returned values to
	match against the relevant fields
*/

/*
type	slotFileJSON struct {

	AirlineCode		string	`json:"airlineName"`
	URI			string	`json:"fileURI"`
	Version			string	`json:"version"`
	DataType		string	`json:"dataType"`
	Format			string	`json:"fileFormat"`
	AssetRef		string	`json:"assetRef"`
	OriginalFileDate	int	`json:"originalFileDate"`
	OriginalFileName	string	`json:"originalFileName"`
	Metadata		string	`json:"metadata"`
}

*/

/*	File actaully returned from the GET /metadata call.  Vary for testing

[20160605-062502]_[FDB-CC46954_9057a54b-16e8-4503-8e35-8c5425f9ed02]_DXB Summer S16 (27Mar-29Oct16)_OBFUSCATED.xlsx

*/

/*	Entire JSON Body of GET

{"airlineName":"FZ","fileURI":"https://bucket-1256b33e-5be8-4168-a7f7-5065f6553c58.store.gecis.io/1208743b-ac49-4ae2-be49-2bc27b68c5dc","version":"1.1","dataType":"event","fileFormat":"xlsx","assetRef":"Flight","originalFileDate":1466635593704,"originalFileName":"[20160605-062502]_[FDB-CC46954_9057a54b-16e8-4503-8e35-8c5425f9ed02]_DXB Summer S16 (27Mar-29Oct16)_OBFUSCATED.xlsx","metadata":null}

*/

func TestSlotNameExtract(t *testing.T) {

//	Intitialize, if needed

	if !TestingInitialized {

		InitTestReadSlotDB() // Load Test Data
		TestingInitialized = true
	}

	fmt.Println ( "--------------------------------------------------" )

	fmt.Println("Starting Testing of slotNameExtract ()")

//	Call Function Under Test 
//	Using test / sample GET Body ( Post-GET call )

//	slotNameExtract ( slotFileName string ) ( string )

	fileNameFromGET := "[20160605-062502]_[FDB-CC46954_9057a54b-16e8-4503-8e35-8c5425f9ed02]_DXB Summer S16 (27Mar-29Oct16)_OBFUSCATED.xlsx"

	fn := slotNameExtract(fileNameFromGET)

//	Validate 

//	Expected Filename based on Get Body ( above )

	expectedFileName := "DXB Summer S16 (27Mar-29Oct16)_OBFUSCATED.xlsx"

	if fn != expectedFileName {

	    not := "FAILED: Unable to Extract"
	    got := "from Slot File Metadata.  Got:"

	    fmt.Println ( not, expectedFileName, got, fn ) 
	    t.Fail() 

	} else { fmt.Println ( "PASSED: Extracted", fn, "from Slot File Metadata" ) }
}

func TestGetMetadata(t *testing.T) {

//	Intitialize, if needed

	if !TestingInitialized {

		InitTestReadSlotDB() // Load Test Data
		TestingInitialized = true
	}


	fmt.Println ( "--------------------------------------------------" )

	fmt.Println("----------Starting Testing of getMetadata()----------")

//	Call Function Under Test ( enough to test each Test Case )

//	getMetadata ( slotFile string ) ( string, string, string, string ) {

	if fromHome { 
		fmt.Println ( "From Home: bypass Test" )
		return
	}
	fmt.Println ( "From Home:", fromHome )
//	return

	file := "DBX_SLOTS_SUMMER"
	airport, carrier, Begin, End := getMetadata( file, false )
	file1 := file

/*	Very little validation is possible ( we have no guaranteed / known
	reponse that can be verified ).  Just check that the four return values
	are non-empty strings.
*/

//	Make sure the Values match the json.Unmarshal () result

	meta := slotFileJSON {}
	for _, meta = range fileMetadataFromGET {}

	fmt.Println ( file, "metadata:", airport, carrier, Begin, End )

	success := true

	if carrier != meta.AirlineCode { 
	    fmt.Println ( "FAIL @ Carrier / Airline:" )
//	    fmt.Println ( "Expected:", carrier, "Got:", airport, meta.AirlineCode )
	    fmt.Println ( "Expected:", carrier, "Got:", meta.AirlineCode )
	    success = false
	    t.Fail() 
	}
/*
	if airport == "" { 
	    fmt.Println ( "FAIL, Blank Airport:", airport )
	    success = false
	    t.Fail() 
	}
	if Begin == "" { 
	    fmt.Println ( "FAIL, Blank Begin Date:", Begin )
	    success = false
	    t.Fail() 
	}
	if End == "" { 
	    fmt.Println ( "FAIL, Blank End Date:", End )
	    success = false
	    t.Fail() 
	}
*/

//	There are two known Seasons ( Summer, Winter ), so at least call again

	file = "DBX_SLOTS_WINTER"
	airport, carrier, Begin, End = getMetadata( file, false )
	file2 := file

	for _, meta = range fileMetadataFromGET {}

	fmt.Println ( "From", file, ":", airport, carrier, Begin, End )
//	fmt.Println ( "Validate using", len (meta), "bytes of File Metadata" )

	if carrier != meta.AirlineCode { 
	    fmt.Println ( "FAIL @ Airline:", airport, meta.AirlineCode )
	    success = false
	    t.Fail() 
	}
/*
	if airport == "" { 
	    fmt.Println ( "FAIL, Blank Airport Code:", airport )
	    success = false
	    t.Fail() 
	}
	if Begin == "" { 
	    fmt.Println ( "FAIL, Blank Begin Date:", Begin )
	    success = false
	    t.Fail() 
	}
	if End == "" { 
	    fmt.Println ( "FAIL, Blank End Date:", End )
	    success = false
	    t.Fail() 
	}
*/
	if success { 
	    fmt.Println ( "PASSED: getMetadata() for", file1, "and", file2 )
	}
}
