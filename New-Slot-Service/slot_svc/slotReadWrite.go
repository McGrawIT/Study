package slot_svc

import (
	"fmt"
	"io"
	"bytes"
	"os"

	blobstore "github.build.ge.com/AviationRecovery/blobstore-support-go"

)



//	Loads JSON Slot File ( typically used when running Local )

func readSlotJSON ( localSlotFile string ) ( []byte, bool ) {

//	Open Slot File (JSON)

	nothing := []byte{}

	slotDB, err := os.Open ( localSlotFile )
	if err != nil { return nothing, false }

//	Close on exit and check for its returned error

	defer func() {
		if err := slotDB.Close(); err != nil { panic(err) }
	}()

	bufSize := 20000000 
	buf := make( []byte, bufSize )

//	Read to Simulate the JSON Body ( for Home Testing )

        n, err := slotDB.Read ( buf )
//	fmt.Println ( "Read", n, "byes from Local JSON File ( readJSON() )" )

        if err != nil && err != io.EOF { return buf, false }

	bufReal := buf[:n]

	return bufReal, true		// Return with simulated GET Body
}

/*********************************************************************
	Currently not used ( no one calls it ) ... keep ?
**********************************************************************/

func writeSlotContents ( contents []byte, slotFile string ) ( bool ) {

	fo, err := os.Create( slotFile )
	if err != nil { 
		fmt.Println("Create failed:", slotFile ) 
		return false 
	}

	n := len ( contents ) 
        if _, err := fo.Write(contents[:n]); err != nil { 

		fmt.Println ( "Write failed:", err ) 
		return false
	}
//	fmt.Println ( "Slot Contents ( JSON ) Written:", slotFile ) 

	if err := fo.Close(); err != nil { panic(err) }
	return true
}

/*	Used to Save Slot Data DB ( done each time a new / existing
	Slot File is pulled from Config Service.  ( The DB is read-only
	except when a Slot File arrives. )

	This could change if Cancellations are sent ( save more often )
*/
func writeSlotDB ( slotDB string ) ( bool ) {

/*	In-memory version of Slot DB is the DB of record.  Saving to Blob Store
	is for restart persistence.  ( Only loaded on boot ).

	All Slot Files are kept in the DB (for analytics)
*/
	fo, err := os.Create("SlotFile.DB")

	if err != nil { return false }

// 	Close fo on exit and check for its returned error

	defer func() { if err := fo.Close(); err != nil { panic(err) } }()

/*	Write the JSON input ( the JSON could be denormalized in this function
	from the Slot Data File(s) vs. having callers do it (they may already
	be doing it) )
*/

	DB := []byte(slotDB)
	n := len ( DB ) 

/*	The next Read will always read the entire Slot File DB, so save
	the size for the make () size
*/
	sizeSlotDataDB = n + 2048	// "Cushion", just in case ;)

        if _, err := fo.Write(DB[:n]); err != nil { 

	    fmt.Println ( "Write failed:", err ) 
	    return false

	} else { fmt.Println ( "DB Written ( Local Copy )" ) }

	if fromHome { return true }		// Skip Blob Store write

//	Write Slot DB to Blob Store

	b := bytes.NewBufferString ( slotDB )

	storedUrl := "Not on Predix"
	if port != "2525" && port != "3535" {
		storedUrl, err = blobstore.PutObject ( "SlotDB", b )
	}

	if err != nil { 
	    fmt.Println ( "Error on Blob Store PutObject:", err ) 
	    return false
	}

	fmt.Printf("Slot DB Blob Store URL: %s\n", storedUrl)

	DB2 = storedUrl		// Permanent DB in Blob Store location

	fmt.Printf("Slot DB saved to Blob Store\n")

	return true
}

