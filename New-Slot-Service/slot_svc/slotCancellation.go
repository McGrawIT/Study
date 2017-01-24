package slot_svc

import (
	"fmt"
	"time"
	"strings"
)

/*
//	Load Stub Data relevant to Stub ( Test ) Airports / Carriers, etc.

func slotCancellationStubLoad () {

	initCancelTracking ()

	makeCancellations ()

	fmt.Println ( "Unique Cancels:", len ( UniqueCancels ) )
	for key, _ := range UniqueCancels { 
		i := len ( UniqueCancels[key] )
		if i > 1 { 
			fmt.Println ( key, ":", len ( UniqueCancels[key] ) )  
		}
	}

	fmt.Println ( "" )
	fmt.Println ( "Carriers with Cancellations", len ( CarrierCancels ) )
	for key, _ := range CarrierCancels { fmt.Println ( key, ":", len ( CarrierCancels[key] ) ) }
	fmt.Println ( "" )
	fmt.Println ( "Flights with Cancellations:", len ( FlightCancels ) )
	for key, cancelSet := range FlightCancels { 

		fmt.Println ( "" )
		fmt.Println ( key, ":", len ( FlightCancels[key] ) ) 
		if key == "20" || key == "88888" {
			fmt.Println ( key, ":", cancelSet )
		}

		for count, value := range FlightCancels[key] { 
			fmt.Println ( key, ":", value )
			if count == 3 { break }
		}
	}


	fmt.Println ( "" )
	fmt.Println ( "Opersation Dates with Cancellations:", len ( OpDate ) )
	for key, cancelSet := range OpDate { 

		fmt.Println ( "" )
		fmt.Println ( key, ":", len ( OpDate[key] ) ) 
		if key == "2016-08-14" || key == "2016-08-13" {
			fmt.Println ( key, ":", cancelSet )
		}

		for count, value := range OpDate[key] { 
			fmt.Println ( key, ":", value )
			if count == 3 { break }
		}
	}

	fmt.Println ( "" )
	fmt.Println ( "Years with Cancellations:", len ( Year 	) )
	for key, _ := range Year { fmt.Println ( key, ":", len ( Year[key] ) ) }

	fmt.Println ( "" )
	fmt.Println ( "Months with Cancellations:", len ( Month 	) )
	for key, _ := range Month { 

		month := "January"

		switch key {
		case "02": month = "February"
		case "03": month = "March"
		case "04": month = "April"
		case "05": month = "May"
		case "06": month = "June"
		case "07": month = "July"
		case "08": month = "August"
		case "09": month = "September"
		case "10": month = "October"
		case "11": month = "November"
		case "12": month = "December"
		}
		fmt.Println ( month, ":", len ( Month[key] ) ) 
	}

	fmt.Println ( "" )
	fmt.Println ( "Weeks of Year with Cancellations:", len ( Week ) )

	for key, cancelSet := range Week { 
		fmt.Println ( "" )
		fmt.Println ( key, ":", len ( Week[key] ) )

		if key == 32 { fmt.Println ( key, ":", cancelSet ) }

		for count, value := range Week[key] { 
			fmt.Println ( key, ":", value )
			if count == 3 { break }
		}
	}

	fmt.Println ( "" )
	fmt.Println ( "Days of the Weeek with Cancellations:", len ( WeekDay ) )
	for key, cancelSet := range WeekDay { 

		fmt.Println ( "" )
		fmt.Println ( key, ":", len ( WeekDay[key] ) ) 

		if key == "Saturday" || key == "Friday" {
			fmt.Println ( key, ":", cancelSet )
		}

		for count, value := range WeekDay[key] { 
			fmt.Println ( "Cancellations on", key, ":", value )
			if count == 3 { break }
		}
	}

	fmt.Println ( "" )
	fmt.Println ( "Unique Types of Cancellation:", len ( Types ) )
	for key, _ := range Types { 
		fmt.Println ( "Cancellation due do", key, ":", len ( Types[key] ) ) 
	}

	fmt.Println ( "" )
	fmt.Println ( "Sources reporting Cancellations:", len ( Source ) )
	for key, cancelSet := range Source { 

		fmt.Println ( "Cancellations reported by", key, ":", len ( Source[key] ) ) 
		if key == "ATC" { fmt.Println ( key, ":", cancelSet ) }
	}
}
*/

/*	Functions to drive Cancellation generation
*/
func	makeOneCancellations () {

	ac := "MCO"
	cc := "DL"
	fn := "20"
	od := "2016-05-31"

	ct := "Weather"
	cr := "High Winds"
	cn := "Over 75 mph"
	cs := "Tower"

	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )
}

func	makeCancellations () {

	ac := "MCO"
	cc := "DL"
	fn := "20"
	od := "2016-05-31"

	ct := "Weather"
	cr := "High Winds"
	cn := "Over 75 mph"
	cs := "Tower"

	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

	ct = "Weather"
	cr = "Vary High Winds"
	cn = "Over 200 mph"
	cs = "Gate Agent"

	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

	cs = "A770 Choice"

	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

	fn = "70"
	cs = "Gate Agent"
	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

	cs = "ATC"
	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

	ct = "Mechanical"
	cr = "Air Conditioning"
	cn = "Over 90 in Cabin"
	cs = "Airline"
	ingestCancel ( ac, cc, fn, od, ct, cr, cn, cs )

//	carriers := []string{"DL","FZ","B6",}
//	flights := []string{"21","22","23","24","25","26",}
//	dates := []string{"2014-05-31","2015-07-31","2017-02-12","2017-02-11",}

	airports := []string{"MCO","DXB",}
	carriers := []string{"DL","FZ",}
	flights := []string{"20","120","18","7146","723","7027",}
	dates := []string{"2016-05-31","2016-09-01","2016-05-30","2016-10-11","2012-10-11","2012-10-12",}

	types := []string{"Weather","Mechanical","Crew","Passenger",}
	sources := []string{"Tower","Gate","Airline",}

//	Cancellation Creation

	total := 1

	for j, fn := range flights {

	    fn2 := fn
	    switch j {
	    	case 1: 
			cr = "Hurricane Charley (2005)"
	    	case 2: 
			fn2 = "888"
			cr = "Hurricane ( Donna ) leveled MCO"
	    	case 3: 
			ct = "Mechanical"
			cr = "Wing fell off @ 35k feet"
	    	case 4: 
			fn2 = "7000"
			ct = "Passenger"
			cr = "Seizure ( Grand Mal )"
	    	case 5: 
			fn2 = "555"
			cr = "Tornados"
	    	case 6: 
			ct = "Crew"
			cr = "Missed Connection @ DFW"
	    }

	    for _, ac := range airports {
	    for _, cc := range carriers {
	    	for k, od := range dates {

		    od2 := od
		    ingestCancel ( ac, cc, fn2, od2, ct, cr, cn, cs )
		    total++

		    od2 = "2016-08-14"
		    ingestCancel ( ac, cc, fn2, od2, ct, cr, cn, cs )
		    total++
			

		    if k == 3 { continue }

		    od2 = "2016-08-12"
		    ingestCancel ( ac, cc, fn2, od2, ct, cr, cn, cs )
		    total++

		    od2 = "2016-08-13"
		    ingestCancel ( ac, cc, fn2, od2, ct, cr, cn, cs )
		    total++

		    for _, ct := range types {
		    	for k, cs := range sources {

			    cs2 := cs
			    ingestCancel ( ac, cc, fn2, od, ct, cr, cn, cs2 )
			    total++

			    cs2 = "Gate Agent"
			    ingestCancel ( ac, cc, fn2, od, ct, cr, cn, cs2 )
			    total++

			    cs2 = "Air Traffic Control"
			    ingestCancel ( ac, cc, fn2, od, ct, cr, cn, cs2 )
			    total++

			    if k == 2 { continue }

			    cs2 = "Airline"
			    ingestCancel ( ac, cc, fn2, od, ct, cr, cn, cs2 )
			    total++

		    	}
		    }
		}
	    }
	    }
	}
	fmt.Println ( "Total:", total )
}

/*	Endpoints / Parameters

	/CancelStats?TotalCancels=Flights, WeekDays, Weeks, ...
	/CancelStats?MaxCancels=
	/CancelStats?MinCancels=
	/CancelStats?Categories=
	/CancelStats?Reasons=

	<...>
*/

var	initCancels	= false

var	newEvent	[]cancelRecord

//	All Cancellations, indexed by Cancel Key

var	UniqueCancels 	map [ cancelKey ] []cancelRecord

//	Cancellation Tracking by Key Categories

var	AirportCancels 	map [ string ] []cancelRecord
var	CarrierCancels 	map [ string ] []cancelRecord
var	FlightCancels 	map [ string ] []cancelRecord

var	OpDate 		map [ string ] []cancelRecord

var	Year 		map [ int ] []cancelRecord
var	Month 		map [ string ] []cancelRecord
var	Week 		map [ int ] []cancelRecord
var	WeekDay 	map [ string ] []cancelRecord

var	Types 		map [ string ] []cancelRecord
var	Source 		map [ string ] []cancelRecord

//	Initializations, run once upon first REST Endpoint access

/*	Cancellations will be collected / grouped for Stats / Summary
	a wide variety of Slot / Cancellation events

	Day of Week		( Sunday, ... )
	Operation Date		( 2016-06-12 / 2016-Jun-12 / 16-06-12 )
	Flight Number		( Numeric String )
	Carrier Code		( FZ, DL, B6, ... )

	Airport Code		( LAX, MCO, JFK, DXB, ... )

//	These two are needed as part of Unique ID

	FlightOrigin
	FlightDestination

	Week of Year		( Week Number ( 1 - 52 )
	Month of Year		( Month Number ( 1 - 12 )
	Year			( Year Number ( 2016, ... )
	Cancel Category / Type 	( Weather, Mechanical, Crew, ... )
	Cancel Source	 	( Gate Agent, Tower, ATC, Carrier, ... )
*/

func initCancelTracking () {

	initCancels = true

	newEvent = []cancelRecord{}	// Empty Map for new Cancel Events

	UniqueCancels = make ( map [ cancelKey ] []cancelRecord )

	AirportCancels = make ( map [ string ] []cancelRecord )
	CarrierCancels = make ( map [ string ] []cancelRecord )
	FlightCancels = make ( map [ string ] []cancelRecord )

	OpDate = make ( map [ string ] []cancelRecord )
	Year = make ( map [ int ] []cancelRecord )
	Month = make ( map [ string ] []cancelRecord )
	Week = make ( map [ int ] []cancelRecord )
	WeekDay = make ( map [ string ] []cancelRecord )

	Types = make ( map [ string ] []cancelRecord )
	Source = make ( map [ string ] []cancelRecord )
}

//	Detailed information about a Slot Cancellation

//	This record format will be used for each Result Set instance

type	cancelRecord 	struct {

	Airport	    string	`json:"airport"`// 3-digit Code (MCO, LAX, ...)
	Carrier	    string	`json:"carrier"`// 2-digit Code (DL, FZ, ...)
	Flight	    string	`json:"flight"`	// Flight Number
	OpDate	    string	`json:"opdate"`	// Operation Date (2016-May-31)

//	Airport Codes for Flight Origin and Destination

	Origin	    string	`json:"origin"`	
	Destination string	`json:"destination"`

	DayOfWeek   string	`json:"weekday"`// Day of Week (Saturday, ..)
	Category    string	`json:"type"`	// Weather, Mechanical

	Reason	    string	`json:"reason"`	// Hurricane, Flooding, ....
	Notes	    string	`json:"notes"`	// Freeform Text
	Source	    string	`json:"source"`	// Gate Agent, Tower, Airline
}

//	Key to unique Slot Cancellation that may be reported multiple times
//	from different sources ( 2.5 Sources / unique Cancel on average )

type	cancelKey	struct {

	airport		string		// Airport Code
	origin		string		// Flight Origin
	destination	string		// Flight Destination
	carrier		string		// Carrier Code
	flight		string		// Flight Number
	opDate		string		// Operation Date
}

func ingestCancel ( ac, cc, fn, od, category, reason, notes, source string ) {

	canceEventKey := cancelKey{}		// Access Key for New Cancel 
	cancelEvent := cancelRecord{}		// New Cancellation Record

/*	Index / Key to access latest Cancellation

	This is not always a new one.  More than one Source can report the 
	same Cancellation.  Typically 2.5 sources / cancellation
*/
	canceEventKey.airport = ac
	canceEventKey.carrier = cc
	canceEventKey.flight = fn
	canceEventKey.opDate = od

/*	Populate detail for the latest Cancellation ( additional instances
	of a Cancellation will be attached to the Cancellation for unique
	Source, Notes, Reason ( can be different for same Category / Type )
*/
//	Some collections of interest are derived from Month

//	od = dateNormalize ( od )	// Assume "2016-07-26" format for now

	year, month, week, day := YearMonthWeekDay ( od )

	cancelEvent.Airport = ac
	cancelEvent.Carrier = cc
	cancelEvent.Flight = fn
	cancelEvent.OpDate = od
	cancelEvent.DayOfWeek = day
	cancelEvent.Category = category
	cancelEvent.Reason = reason
	cancelEvent.Notes = notes
	cancelEvent.Source = source

/*	All Slot Cancellations are captured, by Slot Cancel Key.  If this is
	the first Cancellation for the Slot, initialize it's record
*/
	if len ( UniqueCancels [ canceEventKey ] ) == 0 {

		UniqueCancels [ canceEventKey ] = newEvent
	}
	UniqueCancels [ canceEventKey ] = append ( UniqueCancels [ canceEventKey ], cancelEvent )

	if len ( AirportCancels [ ac ] ) == 0 { AirportCancels [ ac ] = newEvent }
	AirportCancels [ ac ] = append ( AirportCancels [ ac ], cancelEvent )

	if len ( WeekDay [ day ] ) == 0 { WeekDay [ day ] = newEvent }
	WeekDay [ day ] = append ( WeekDay [ day ], cancelEvent )

	if len ( OpDate [ od ] ) == 0 { OpDate [ od ] = newEvent }
	OpDate [ od ] = append ( OpDate [ od ], cancelEvent )

	if len ( FlightCancels [ fn ] ) == 0 { FlightCancels [ fn ] = newEvent }
	FlightCancels [ fn ] = append ( FlightCancels [ fn ], cancelEvent )

	if len ( CarrierCancels [ cc ] ) == 0 { CarrierCancels [ cc ] = newEvent }
	CarrierCancels [ cc ] = append ( CarrierCancels [ cc ], cancelEvent )


	if len ( Week [ week ] ) == 0 { Week [ week ] = newEvent }
	Week [ week ] = append ( Week [ week ], cancelEvent )

	if len ( Month [ month ] ) == 0 { Month [ month ] = newEvent }
	Month [ month ] = append ( Month [ month ], cancelEvent )

	if len ( Types [ category ] ) == 0 { Types [ category ] = newEvent }
	Types [ category ] = append ( Types [ category ], cancelEvent )

	if len ( Source [ source ] ) == 0 { Source [ source ] = newEvent }
	Source [ source ] = append ( Source [ source ], cancelEvent )

	if len ( Year [ year ] ) == 0 { Year [ year ] = newEvent }
	Year [ year ] = append ( Year [ year ], cancelEvent )

}


func YearMonthWeekDay ( dateIn string ) ( int, string, int, string ) {

	layout := "2006-01-02T15:04:05.000Z"
	suffix := "T11:45:26.371Z"

	usefulDate := dateIn + suffix

	t, _ := time.Parse ( layout, usefulDate )

	year, week := t.ISOWeek()

	monthNum := strings.Split ( dateIn, "-" )
	month := monthNum[1]

	x := t.Weekday()
	dayOfWeek := x.String()

	return year, month, week, dayOfWeek
}

