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

func TestUpdateFlightSlots ( t *testing.T ) {

//	Intitialize, if needed

/*	Since Go Test does not guarantee order of execution, all _test.go files
	will check this variable to make sure the required initializations are 
	performed before any test (and happens only once / go test execution 
*/
	if !TestingInitialized {

	    InitTestReadSlotDB ()	// Load Test Data
	    TestingInitialized = true
	}

//	Special-case initializations required for this specific _test.go file

	fmt.Println ( "--------------------------------------------------" )
	fmt.Println ( "Starting Testing of UpdateFlightSlots()" )

/*	Call function under test as many times as needed to at least test the
	primary aspects of this function for a thorough Unit Test.  
	Report / log errors with t.Fail()
*/

//	UpdateFlightSlots()


}
