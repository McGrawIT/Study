package slot_svc

import (
	"strings"
	"time"
	"strconv"
)

type	slotDate	struct {

	Year		int
	Month		int
	Day			int
	YMD			string		// YY-MMM-DD
	LongForm	string		// YYYY-MMM-DD
}

/*	Simple check; Is the Operation Date in Season?  ( Can be used for any "is date func init() {
 	range? comparison )  Added to allow simple Date String compares
 */
func inSeason ( operationDate, seasonBegin, seasonEnd string ) ( bool ) {

	pattern := "2006-Jan-02"

	begin, _ := time.Parse ( pattern, seasonBegin )
	end, _ := time.Parse ( pattern, seasonEnd )
	opDate, _ := time.Parse ( pattern, operationDate )

	return opDate.Before( end ) && opDate.After ( begin )
}


/*	Date Normalization() will convert several date strings into YYYY-MM-DD
	for consistent use in the Slot DB Key and for the expected format of
	all Slot Service functions ( e.g., Search )
*/

var	dayString = map[ string ]string{ "01":"Jan", "02":"Feb", "03":"Mar", "04":"Apr", "05":"May", "06":"Jun", "07":"Jul", "08":"Aug", "09":"Sep", "10":"Oct", "11":"Nov", "12":"Dec" }

/*	While API officially calls for YYYY-MMM-DD, add support to "correct" likely valid
	dates ( e.g., 2016-10-23 should be 2016-Oct-23 )

	Another case would be single digit days:  2016-Oct-1
	This one causes the time.Parse() to fail ( it expects 2-digit days )

	Check for YY ( 16 vs 2016 ) assume 21st Century, use time.Now() to avoid Y2K-type error
 */

func opDatesNormalize ( includedDates []string ) ( []string ) {

	datesStringMonth := []string{}		// No dates

	for _, dateMMM := range includedDates {

//		fmt.Println ( "Filter Date:", dateMMM )

		dateParts := strings.Split ( dateMMM, "-" )

		year, _, _ := time.Now().Date()		// Do not hard-code Century ( remember Y2k )
		century := int ( year / 100 )

		if len ( dateParts[0] ) == 2 { 										// Year
			dateParts[0] = strconv.Itoa ( century ) + dateParts[0]
		}
		if len ( dateParts[1] ) == 1 { dateParts[1] = "0" + dateParts[1] }	// Month
		if len ( dateParts[2] ) == 1 { dateParts[2] = "0" + dateParts[2] }	// Day

/*		Convert any dates of the form YYYY-MM-DD to YYYY-MMM-DD
		Example: 2018-10-31 to 2018-Oct-31 ( if already alpha, use it )
 */
		dayNumber := dayString[ dateParts[1] ]
		newDate := ""

		if dayNumber == "" { newDate = dateParts[0] + "-" + dateParts[1] + "-" + dateParts[2]
		} else { newDate = dateParts[0] + "-" + dayNumber + "-" + dateParts[2] }

		datesStringMonth = append (  datesStringMonth, newDate )
	}
	return datesStringMonth 
}


func dayOfWeek ( date string ) ( string, int ) {

	if date == "" { return "" , 0 }

	date1 := date + "T00:00:00.000Z"

	layout := "2006-Jan-02T15:04:05.000Z"

	t, _ := time.Parse( layout, date1 )
	x := t.Weekday()

	dayNum := int ( x )

	if dayNum == 7 { dayNum = 1 } else { dayNum++ }

	day1 := x.String()

//	fmt.Println ( "dayOfWeek() is", day1, "Input date:", date )

	return day1, dayNum
}


