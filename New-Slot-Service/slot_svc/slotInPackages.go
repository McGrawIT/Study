package slot_svc

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"bytes"
	"log"
)

/*
	Ideas / Sammple "Info" Page (from Wayne's Disruption Service)

Received Ping on ... Saturday, 04-Jun-16 02:07:59 UTC 
Disruption Management Microservice - Copyright GE 2015-2016

Instance GUID: 48e54136-cc75-421b-9b7d-bec14548109e
Instance IP: 10.131.17.131:61490
Instance #: 0
Started: Thursday, 19-May-16 19:29:00 UTC
Uptime:  366h38m59.205947775s

System Operation
  Alloc: 2356728  (bytes allocated and not yet freed)
  Total Alloc: 157173272  (bytes allocated (even if freed))
  Sys: 15280376  (bytes obtained from system)
  Free Memory: 12923648  (Sys - Alloc)
  Lookups: 1122  (number of pointer lookups)
  Mallocs: 1551935  (number of mallocs)
  Frees: 1535088  (number of frees)

Configuration Variables
  URI:  disruption-servicegit.run.asv-pr.ice.predix.io
  PORT: 61490
  COST_SERVICE: disruption-cost-servicegit.run.asv-pr.ice.predix.io,aviation-flydubai-eu261-svc.run.asv-pr.ice.predix.io
  PAX_REFLOW_SERVICE: aviation-flydubai-pax-reflow-svc.run.asv-pr.ice.predix.io/paxreflow

*/

func getInfo ( response http.ResponseWriter, query *http.Request ) {

//	Use Maps to Capture # of Unique Airports, Carriers, Seasons, Flights
//	May be able to capture Carries by Airport, Airports by Carrier, etc.o

//	Probably too hard to maintain and better to just calcuate it

//	fmt.Printf ( "Airports:", len (AC) )
//	fmt.Printf ( "Carriers:", len (CC) )
//	fmt.Printf ( "Seasons:", len (DR) )

/*	Ideas for "Info" response:

	Slot File Dates (Created, Updated), all Slot Files
	Slot DB Stats
	Slot File Stats
		# of Airports
		# of Carriers
		# of Seasons
		# of Slots / Season


*/

}


/*	Logic for Read / Write of SF DB

	( also, look @ comments from old convertExcel() code ( below ) )

	WRITE

	Use ssFilter (ALL) to select everything
	Flatten it
	Convert to JSON
	Write JSON

	READ

	Read DB File
	Marshal() Content (need byte slice) into a slice of Flatten Slots
	Create Empty SF DB
	Range over the Flat Slots
		Extract the Key
		Extract the Flight Slot
		If len (SFDB[key]) == 0
			Insert first Flight Slot
		else
			Append next Flight Slot
	End Range

	If saved, compare input DB stats with Created DB stats
	May use a separate file for these "stats"
*/

func rootEndPoint ( response http.ResponseWriter, query *http.Request ) {

	fmt.Fprintf ( response, query.URL.Host )
	fmt.Fprintf ( response, " (HOST) \n" )
	fmt.Fprintf ( response, "Default Endpoint: Crazy Bad Shit\n" )

	fmt.Fprintf ( response, "PATH:" )
	fmt.Fprintf ( response, query.URL.Path )
	fmt.Fprintf ( response, "\n" )

	fmt.Println (  "==> HOST:", query.URL.Host )
	fmt.Println (  "==> PATH:", query.URL.Path )
}


/*
	There will be stuff unique to these, including their param values
	Use the map as a proof of "existence" (is the GET parm valid?)
	This would be part of an init function

type urlParameters struct {
	currentSet	[]string
	otherstuff	[]string
}
	validParams := make ( map[string]urlParameters )

	validParams["carrierCode"] = "OK"
	validParams["airportCode"] = "OK"
	validParams["operationDate"] = "OK"
	validParams["flightNumber"] = "OK"
	
	Better than switch

	Can the parms (and their vaiues) form an index?

	Carriers, Code ... is that too much overhead vs. just doing a
	straight query?  (Walking the structure)

	Look @ the searches and the current storage approach...


	Parm := map[value]	// If not there, return is nil, right?

*/

const chunckSize = 64000

func compareSlotFiles (file1, file2 string) bool {

	//fmt.Println ( "=======>> Comparing:", file1, file2 )

	f1s, err := os.Stat(file1)
	if err != nil { log.Fatal(err) }

	f2s, err := os.Stat(file2)
	if err != nil { log.Fatal(err) }

	if f1s.Size() != f2s.Size() { return false }

	f1, err := os.Open(file1)
	if err != nil { log.Fatal(err) }

	f2, err := os.Open(file2)
	if err != nil { log.Fatal(err) }

	for {
		b1 := make([]byte, chunckSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunckSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF && err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) { return false }
	}
}

func compareSlotFiles2 ( existing, new string ) bool {

	sigExisting, _ := os.Stat ( existing )
	sigNew, _ := os.Stat ( new )

//	fmt.Println ( "Old:", existing )
//	fmt.Println ( "New:", new )

	if sigExisting != sigNew { 
	 //   fmt.Println ( "Files Different" )
	    return false
	}
	return true
}


func Exists( name string ) bool {

    if _, err := os.Stat( name ); err != nil {
    	if os.IsNotExist(err) { return false }
    }
    return true
}

func getSizeSlotDB () (int) {

	f1, err := os.Open( "SlotData.DB" )

	size := chunckSize
	b1 := make( []byte, chunckSize )

	for {
		_, err = f1.Read( b1 )

		size = size + chunckSize

		if err != nil { return size }
	}
	return size
}

