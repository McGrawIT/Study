
package main

import (
	"fmt"

	sr "github.build.ge.com/aviation-intelligent-airport/time-and-distance-svc.git/server"
)

func main() {

	fmt.Println("----------------------------------")
	fmt.Println("Starting Time and Distance Service")
	fmt.Println("----------------------------------")

//  	Start REST Interface here

	sr.SetupRestServer()

	sr.StartRestServer()


}

