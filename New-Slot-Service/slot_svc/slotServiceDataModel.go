package slot_svc

/*	Slot Data File ( from Config Service in JSON )

	Configuration Service will pull the User-Generated / Maintained Slot Data File
	out of Blob Store and convert the Excel to JSON
*/

type	SlotJSON struct {

	FlightNumber			string	`json:"flightnumber"`

	SunFlightsInSeason		int		`json:"sunFlightsInSeason"`
	SunCurrentMax			float64	`json:"sunCurrentMax"`

	MonFlightsInSeason		int		`json:"monFlightsInSeason"`
	MonCurrentMax			float64	`json:"monCurrentMax"`

	TueFlightsInSeason		int		`json:"tueFlightsInSeason"`
	TueCurrentMax			float64	`json:"tueCurrentMax"`

	WedFlightsInSeason		int		`json:"wedFlightsInSeason"`
	WedCurrentMax			float64	`json:"wedCurrentMax"`

	ThuFlightsInSeason		int		`json:"thuFlightsInSeason"`
	ThuCurrentMax			float64	`json:"thuCurrentMax"`

	FriFlightsInSeason		int		`json:"friFlightsInSeason"`
	FriCurrentMax			float64	`json:"friCurrentMax"`

	SatFlightsInSeason		int		`json:"satFlightsInSeason"`
	SatCurrentMax			float64	`json:"satCurrentMax"`
}

//	Primary Key to a Slot Data File in the Slot DB

/*	Key Structures in use....Airport, Carrier, and Season
	This is "one per" Slot Data File, so add any "Slot File Specific"
	Fields (e.g., Created and Most Recent Updated Dates)
*/
/*	Store Parsed / Converted Dates

	(1) From File Name (season / date range, as MonthYear2 ( e.g., Oct16 ), parsed
	to MonthYear4 ( e.g., Oct2016 ), with Begin Date's Day set to 01 and End Date set to
	28 / 30 / 31 ( skipping Lear Year for MVP1 ),
	(2) From Queries for Operation Date ( Break up the YYYY-MM-DD )
*/

type SlotKey struct {

	Airport					string		// 3-Digit Code ( LAX, SFO, etc. )
	Carrier					string		// 2-Digit Code ( FZ, DL, etc. )
	SeasonStart				string		// YYYY-MM-DD format
	SeasonEnd				string		// YYYY-MM-DD format
}


//	The In-Memory Slot Data Fie DB ( Persisted in Blob Store )

var	SlotDataDB 		map[ SlotKey ][]FlightSlot

type FlightSlot 	struct {

	flightNumber 	string
	slotWeek     	[]SlotRecord // Sunday through Saturday
}

type SlotRecord struct {

	Weekday				string	`json:"weekday"`
	OpDate				string	`json:"opDate"`
	FlightsInSeason		int		`json:"flightsInSeason"`

	OriginalMax		int	`json:"originalMaxCancels"`
	CurrentMax		int	`json:"maxCancels"`

	SlotIndex		float64	`json:"slotIndex"`
}


/*	This Row is the fully denormalized

	It forms a row in the Result Set (for Queries) and gets
	marshaled into JSON for response content

	For now, Airport and Date Range are specified in Slot File Name
	Carrier is not specified in filename or file contents (Default: FZ)

	Flight Number is the Key to the Week's Schedule
*/

/*	THE Slot Record
*/

type slots struct {

/*
	For now, we are only dealing with flyDubai (FZ) @ Dubai International
	Carrier code is not in source file (for now), defaults to FZ

 */

/*	Current Slot Data Service code (this file) DOES support >1 Airport
	and >1 Carrier, even if the source data is limited to DXB / FZ
 */
	airportCode			string	// Code from File Name
	carrierCode			string	// Airline Code (FZ only)
	beginDate			string
	endDate				string
	flightNumber		string

	weekday				string	// Sunday, Monday, ...
	opDate				string	// Operation Date (if applicable)

	flightsInSeason		int		// For each Weekday
	currentMaxCancels	int		// For each Weekday
	originalMaxCancels	int		// From first Slot Data File instance

//	slotIndex = ( Original Max Cancels - Current Max Cancels ) / Flights In Season

	slotIndex 			float64

}
/*
	Return Structure, "flattened" ( Denormalized )
 */
type slotPartFlat struct {

	CarrierCode					string	`json:"carrierCode"`
	AirportCode					string	`json:"airportCode"`

	SeasonStart					string		// YYYY-MM-DD format
	SeasonEnd					string		// YYYY-MM-DD format

	Weekday						string	`json:"weekday"`

//	Operation Date will only have a value if one has been supplied as a Query Parameter

	OpDate						string	`json:"OpDate"`

	FlightNumber				string	`json:"flightNumber"`
	FlightsInSeason				int		`json:"flightsInSeason"`
	SlotIndex					float64	`json:"slotIndex"`
	CurrentMaxCancels			int		`json:"maxCancels"`
	OriginalMaxCancels			int		`json:"originalMaxCancels"`
}

/*	slotPart{} not only "flattens" the Slot Data (for Result Sets),
	but it also uses a slightly different field name for Cancels
*/

type slotPart struct {

	CarrierCode					string	`json:"carrierCode"`
	AirportCode					string	`json:"airportCode"`
	Weekday						string	`json:"weekday"`
	OpDate						string	`json:"opDate"`
	FlightNumber				string	`json:"flightNumber"`
	FlightsInSeason				int		`json:"flightsInSeason"`
	SlotIndex					float64	`json:"slotIndex"`
	CurrentMaxCancels			int		`json:"maxCancels"`
	OriginalMaxCancels			int		`json:"originalMaxCancels"`
}


const(

	DEBUG 				= false
	CLIENT_ID 			= "aviation-flydubai-service-client"
	CLIENT_SECRET 		= "yu7c2EScCsb9bEpj"
	SERVICE_CREDENTIAL 	= "YXZpYXRpb24tZmx5ZHViYWktc2VydmljZS1jbGllbnQ6eXU3YzJFU2NDc2I5YkVwag=="

	Sun	=	"Sunday"
	Mon	=	"Monday"
	Tue	=	"Tuesday"
	Wed	=	"Wednesday"
	Thu	=	"Thursday"
	Fri	=	"Friday"
	Sat	=	"Saturday"

)

//	Slot Data File specific information

var		SlotFileInfo	map[SlotKey]SlotDataFileDetails
var		sizeSlotDataDB	int

type	SlotDataFileDetails	struct {

	createdDate				string	// First load of this Slot File
	lastUpdatedDate			string	// Last time Slot details changed
	lastReceivedDate		string	// Last Time Slot File was pulled
}

var (

	SlotSeasons			map[ string ]string
	Stub				bool
	slotsLoaded			bool
	dumpAllSlots		bool
	doDebug				bool
	initialLoad			bool
	stubSlotIndex		float64
	showEndPoints 		bool

	//	slotFileURL: slotConfigHostName + slotConfigEndpoint + currentSlotFile

	slotFileURL			string		// Host Name + Endpoint + Current File

	slotConfigHostName 	string		// Host Name ( Config Service )
	slotConfigEndpoint	string		// End Point ( Config Service )
	currentSlotFile		string

)

type	slotFileMetadataJSON struct {

	AirlineCode					string	`json:"airlineName"`
	URI							string	`json:"fileURI"`
	Version						string	`json:"version"`
	DataType					string	`json:"dataType"`
	Format						string	`json:"fileFormat"`
	AssetRef					string	`json:"assetRef"`
	OriginalFileDate			int		`json:"originalFileDate"`
	OriginalFileName			string	`json:"originalFileName"`
	Metadata					string	`json:"metadata"`
}

//	Potential Key ( IDs specific Flight in SlotData File Data Range )

type Airport struct {

	AirportCode				string
	CountryCode				string
	Location				string
	TimeZone				string		// May have Date implications
	AirportDetail			string		// Expand
}

//	Is there Carrier-specific Slot details / content?

type Carrier struct {

	CarrierCode				string
	CarrierDetail			string		// Expand
}

//	Possible Secondary Key? ( IDs specific Flight in SlotData File? )

type flightKey struct {

	Airport					string
	Carrier					string
	SeasonStart				string
	SeasonEnd				string
	Flight					string
}


//	Future Use?

var	(

	CarrierOriginations		map[ string][]Airport		// Carrier flies out of these
	CarrierDestinations        map[ string][]Carrier	// Carrier flies into these Airports
	AirportCarriers        map[ string][]Carrier		// Airports has these Carriers
)

type slotDay struct {

	weekday				string	// Positional Columns
	flightsInSeason		int		// For each Weekday	
	currentMaxCancels	int		// For each Weekday
}


/*	SlotData Query Parameters ( GET / SlotData )

	For now, filtering only requires AND of 1-n parameters
	None = ALL
*/

type SlotDate struct {		

	month	int
	day		int
	year	int
}

type	Slot	struct {

	Day		string
	Max		int
}

type	unpackKey	struct {

	AP			string
	Begin		string
	End			string
}

type	ForJSON	struct {

	Key		unpackKey
	Slots	[]Slot
}

type	ForPartJSON	struct {

	Key		SlotKey
	Slots	[]slotPart
}



