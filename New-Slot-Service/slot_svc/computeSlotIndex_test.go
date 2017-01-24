package slot_svc

import (
	"testing"
	"fmt"
//	"os"
//	"io/ioutil"
//	"strconv"
//	"strings"
//	"net/http"
//	"encoding/json"

)

//	Total Cancels = Original Max Cancels - Current Max Cancels
//	Cancels = originalMax - currentMax

/*	Slot Index = ( originalMax - currentMax ) / originalMax

	Unless...

	If Original Max Cancels = 0, Slot Index = 0% ( vs. Divide by Zero )
	If Cancels > Original Max Cancels, Slot Index = 100%
	If Carrier = FZ and Flight Number begins with "7", Slot Index = 100%
*/

/*	For all Fly Dubai (FZ) flights from 7000 to 7999, set
	Slot Index to zero (0), regardless of Max Cancels.

	They want to always consider these flight for cancellation
*/

func TestComputeSlotIndex ( t *testing.T ) {


//	Intitialize, if needed

/*	Since Go Test does not guarantee order of execution, all _test.go files
	will check this variable to make sure the required initializations are 
	performed before any test (and happens only once / go test execution 
*/
	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}


	fmt.Println ( "--------------------------------------------------" )

	fmt.Println ( "Starting Testing of ComputeSlotIndex ()" )

	bad := 0
	good := 0

//	ComputeSlotIndex ( current, original int, flight string ) ( float64 )

//	Slot Index = ( originalMax - currentMax ) / originalMax

	slotIndex := ComputeSlotIndex ( 6, 6 , "FZ70" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.00, got", slotIndex )

	} else { good++ }

	slotIndex = ComputeSlotIndex ( 0, 6 , "FZ1745" )
	if slotIndex != 1.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 1.00, got", slotIndex )
	} else { good++ }

	slotIndex = ComputeSlotIndex ( 3, 6 , "FZ70" )
	if slotIndex != 0.5 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.50, got", slotIndex )
	} else { good++ }

	slotIndex = ComputeSlotIndex ( 2, 8 , "FZ284" )
	if slotIndex != 0.75 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.75, got", slotIndex )
	} else { good++ }

//	If Carrier = FZ and Flight Number begins with "7", Slot Index = 100%
//	Caller will concatentate Carrier Code and Flight Number 

	slotIndex = ComputeSlotIndex ( 6, 6 , "FZ7000" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.00, got", slotIndex )
	} else { good++ }

	slotIndex = ComputeSlotIndex ( 6, 6 , "FZ7999" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.00, got", slotIndex )
	} else { good++ }

//	A 7000 Flight, but NOT Fly Dubai

	slotIndex = ComputeSlotIndex ( 6, 6 , "DL7000" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.00, got", slotIndex )
	} else { good++ }

//	If Original Max Cancels = 0, Slot Index = 0% ( vs. Divide by Zero )

	slotIndex = ComputeSlotIndex ( 2, 0 , "FZ7" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 0.00, got", slotIndex )
	} else { good++ }

//	Cancels = originalMax - currentMax
//	If Cancels > Original Max Cancels, Slot Index = 100%

	slotIndex = ComputeSlotIndex ( -2, 8 , "FZ11" )
	if slotIndex != 1.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED: Compute Slot Index() Expecting 1.00, got", slotIndex )
	} else { good++ }

//	If Original Max Cancels < Max Cancels, Slot Index = 0%

	slotIndex = ComputeSlotIndex ( 9, 8 , "FZ11" )
	if slotIndex != 0.0 { 

	    bad++
	    t.Fail()
	    fmt.Println ( "FAILED:  Compute Slot Index() Expecting 0.00, got", slotIndex )
	} else { good++ }

	if bad > 0 {
	    fmt.Println ( "FAILED: (Summary) Compute Slot Index()", good, "PASSED /", bad, "FAILED" )
	} else {
	    fmt.Println ( "PASSED: Compute Slot Index() ", good, "Test Cases PASSED" )
	}
}
