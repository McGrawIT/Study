package slot_svc

import (
	"fmt"
	"strings"
	"strconv"
)


type	indexFilters 	struct {

	symbol		string		// G, L, E, R
	value		float64
	upper		float64		// Upper bound of Range only
}

/*	Filters are in an OR selection, return TRUE on first pass, only
	return false if all filters fail

	To filter by the Slot Index, the ?Index query parameter can be one or more
	of the following ( examples are used for the values )
	Commas allow for more than one, and are used as AND filters

	GT 0.5,  LT 0.3,  EQ 0.75,  RANGE 0.0:0.25
*/

func	slotIndexFiltering ( test string ) {

	fmt.Println ( test )

	filters := parseSlotIndexQuery ( test )
	fmt.Println ( filters )

	slotIndex := 0.4
	keep := keepIndex ( filters, slotIndex )
	if !keep { fmt.Println ( "PASS:", slotIndex, "not Included" ) }

	slotIndex = 0.8
	keep = keepIndex ( filters, slotIndex )
	if keep { fmt.Println ( "PASS:", slotIndex, "Included (GT)" ) }

	slotIndex = 0.3333
	keep = keepIndex ( filters, slotIndex )
	if !keep { fmt.Println ( "PASS:", slotIndex, "not Included" ) }

	slotIndex = 0.1
	keep = keepIndex ( filters, slotIndex )
	if keep { fmt.Println ( "PASS:", slotIndex, "Included (LT, IR)" ) }

	slotIndex = 0.75
	keep = keepIndex ( filters, slotIndex )
	if keep { fmt.Println ( "PASS:", slotIndex, "Included (EQ)" ) }

}

/*
	Determine if Current Slot Index passes at least one of the Index Filters
	(G)reater than, (L)ess than, (E)qual to, (N)ot Equal to, in (R)ange
 */

func	keepIndex ( filters []indexFilters, slotIndex float64 ) ( bool ) {
	
	for _, filter := range filters {

		inValue := filter.value

		switch filter.symbol {
		case "G" : if slotIndex > inValue { return true }
		case "E" : if slotIndex == inValue { return true }
		case "L" : if slotIndex < inValue { return true }
		case "N" : if slotIndex != inValue { return true }
		case "R" : if slotIndex > inValue && slotIndex < filter.upper { return true }
		}
	}
	return false
}


func parseSlotIndexQuery ( test string ) ( []indexFilters ) {

	slotIndexConditions := strings.Split ( test, "," )

	selectIndex :=	[]indexFilters{}

	arg := ""
	sym := ""

	oneFilter := indexFilters{}

	for _, item := range slotIndexConditions {
	
		oneFilter.upper = 0.0	// Only RANGE has non-zero Upper
		
		trimmed := strings.Trim(item, " ")

		arg = strings.Trim(trimmed[2:], " ")
		sym = trimmed[0:1]

		switch trimmed[0:1] {
//		case "E" : fmt.Println ( "Equals", arg )
//		case "G" : fmt.Println ( "Greater than", arg )
//		case "L" : fmt.Println ( "Less than:", arg ) 
//		case "N" : fmt.Println ( "Not Equal:", arg ) 
		case "R" :

			args := strings.Trim(trimmed[5:], " ")
//			fmt.Println ( "Range String:", args ) 

//			args holds range ( upper:lower )

			indexRange := strings.Split ( args, ":" )

//			fmt.Println ( "Range Values:", indexRange )

			for whichBound, value := range indexRange {

				if whichBound == 0 { 

					arg = value
					continue 
				}

				oneFilter.upper = makeFloat ( value )
			}
			
		default :

		}

/*		Convert argument(s) from string to float64
*/
		oneFilter.value = makeFloat ( arg )
		oneFilter.symbol = sym			// Conditional

		selectIndex =	append ( selectIndex, oneFilter )
	}
//	fmt.Println ( "selectIndex:", selectIndex )
	return selectIndex
}

func	makeFloat ( arg string ) ( float64 ) {

	if n, err := strconv.ParseFloat(arg, 64); err == nil { return n }
	return 0.0
}
