package slot_svc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

//	Authorization Package

	au "github.build.ge.com/AviationRecovery/go-oauth.git"
)

/*	Retrieve Slot File Metata, pull Carrier ( Airline ) Code and
	File Name.  File Name must be parsed for Airport and Season
*/

var fileMetadataFromGET []slotFileMetadataJSON

func getMetadata( slotFile string, silent bool ) ( string, string, string, string, string, string ) {

//	Assemble the Correct Endpoint Name (special case)
//	Configuration Service uses "FZ" for Metadata

	hostName := slotConfigHostName
	specialCaseEndpoint := "/api/v1/FZ/"

	metadataContent := slotFile + "/metadata"
	metadataEndpoint := hostName + specialCaseEndpoint + metadataContent

	fmt.Println( "Getting Slot File Metadata for", slotFile )
	fmt.Println( "Using:", metadataEndpoint, "in clent.GET() call" )

	if fromHome { 

	    fmt.Println ( "From Home:", slotFile, "fake return values" )
	    if slotFile == "DBX_SLOTS_WINTER" {

		return "DXB", "FZ", "2015-Oct-25", "2016-Mar-26", "Winter 2015", "DXB Winter 2016 (25Oct-26Mar16).xlsx"
	    } 
	    return "DXB", "FZ", "2016-Mar-27", "2016-Oct-29" , "Summer S16", "DXB Summer 2016 (25Oct-26Mar16).xlsx"
	}

//	GET File Metadata from Config Service ( get Authorization first )

	client := au.RequestAuthorization()

	success := true

	response, err := client.Get(metadataEndpoint)

	if err != nil { 

	    fmt.Println( "Slot File:", slotFile )
	    fmt.Println( "Using:", metadataEndpoint, "in clent.GET()" )
	    fmt.Println( "GET Failed:", err ) 
	    success = false

	} else if response == nil { 

	    fmt.Println( "No error on GET, but Reponse is Empty" ) 
	    success = false

	} else { defer response.Body.Close() }

	payload := []byte{}
	contents := []byte{}

	if success { 

	    contents, err = ioutil.ReadAll(response.Body) 
	    if err != nil {
		fmt.Println( "Unable to read reponse.Body:", err )
		success = false
	    }
	}

	if !success { 

	    fmt.Println ( "Metadata:", slotFile, "Not Retrieved, Defaulting " )

	    if slotFile == "DBX_SLOTS_WINTER" {

		return "DXB", "FZ", "2015-Oct-25", "2016-Mar-26", "Winter 2015", "DXB Summer S16 (27Mar-29Oct16).xlsx"
	    } 
	    return "DXB", "FZ", "2016-Mar-27", "2016-Oct-29", "Summer S16", "DXB Summer S16 (27Mar-29Oct16).xlsx"
	}

//	Metadata Retrieved
//	Add brackets to JSON Body to create []byte for unmarshal()

	payload = OneAsSlice(contents, silent)

/*	Load Metadata from GET into veriable visible to slot_svc package
	( This allows for better testing ( like did the funciton put the
	JSON body in the structure correctly? )
*/
	fileMetadataFromGET = []slotFileMetadataJSON{}

	err = json.Unmarshal(payload, &fileMetadataFromGET)

	if err != nil {

	    fmt.Println( "Faied:  json.Unmarshal()", err )
	    success = false

	} else { 

	    file := slotFile 
	    slots := len(fileMetadataFromGET)
	    fmt.Println("Unmarshaled", file, slots, "Metadata Records") 
	}

//	Do not continue, otherwise unexpected results / errors likely

	if !success { 

	    fmt.Println ( "Metadata:", slotFile, "Not Retrieved, Defaulting " )
	    if slotFile == "DBX_SLOTS_WINTER" {

		return "DXB", "FZ", "2015-Oct-25", "2016-Mar-26", "Winter 2015", "DXB Winter 2016 (25Oct-26Mar16).xlsx"
	    } 
	    return "DXB", "FZ", "2016-Mar-27", "2016-Oct-29", "Summer S16", "DXB Summer S16 (27Mar-29Oct16).xlsx"
	}

//	Extract Carrier Code and File Name from JSON Body

	slotFileName := ""
	carrierCode := ""

	for _, fs := range fileMetadataFromGET {

	    carrierCode = fs.AirlineCode
	    slotFileName = fs.OriginalFileName
	}

	if slotFileName == "" { 

	    fmt.Println ( "Slot Filename for", slotFile, "Empty, Use Default" )
	    if slotFile == "DBX_SLOTS_WINTER" {

		return "DXB", "FZ", "2015-Oct-25", "2016-Mar-26", "Winter 2015", "DXB Winter 2016 (25Oct-26Mar16).xlsx"
	    } 
	    return "DXB", "FZ", "2016-Mar-27", "2016-Oct-29", "Summer S16", "DXB Summer S16 (27Mar-29Oct16).xlsx"
	}

//	Slot File Name is only part of the File Name from JSON Body

	slotFileName = slotFileNameExtract( slotFileName )

	fmt.Println ( "Slot File Name:", slotFileName )

	airportCode, seasonName := getSeasonName ( slotFileName )

//	Parse File Slot Name ( for Season )
//	Ignore AirportCode, CarrierCode

//	Come back and simplify parseSlotFileName() to just get Season Dates

	_, _, BeginDate, EndDate := parseSlotFileName( slotFileName )

	return airportCode, carrierCode, BeginDate, EndDate, seasonName, slotFileName
}

//	Extract name of Excel file from Response Body

func slotFileNameExtract(slotFileName string) string {

//	Presumed Unique Marker that precedes actual Slot File Name
//	The Slot File Name will be the String at the End

	fileMetadata := strings.Split(slotFileName, "]_")

	fileName := ""
	for _, value := range fileMetadata { fileName = value }

	return fileName
}

/*	File Name format ( from Config Service )

	"[20160331-045621]_[FDB-CC84881_33633c38-f0b8-42a3-8d54-2ecb87e07caa]_DXB Winter 2015 (25Oct-26Mar16)_OBFUSCATED.xlsx"

*/

func getSeasonName ( slotFile string ) ( string, string ) {

//	Airport Code and Season Name precede the Season Dates
	
	beforeDates := strings.Split ( slotFile, "(" )
	
	first := ""
	for _, part := range beforeDates { 
	
		first = part		// Converts slice to string
		break
	}
	
	airportCode := first[0:3]

	seasonName := first[3:]
	seasonName = strings.TrimSpace( seasonName )
	
	return airportCode, seasonName
}

//	Surround Metadata Reponse Body with brackets []

func OneAsSlice( contents []byte, silent bool ) []byte {

	onlyOne := []byte{}

	bracket := []byte("[")
	endBracket := []byte("]")

	for _, one := range bracket { onlyOne = append(onlyOne, one) }
	for _, one := range contents { onlyOne = append(onlyOne, one) }
	for _, one := range endBracket { onlyOne = append(onlyOne, one) }
//	if !silent { fmt.Println( "Ready for Unmarshal()" ) }

	contentString := string(onlyOne[:])
//	if !silent { fmt.Println("One as Array:", contentString) }
//	if !silent { fmt.Println("") }

	payload := []byte(contentString)

	return payload
}
