package main

import (
	"strings"
	"fmt"
)

func main () {

	fmt.Println ( "Begin Test" )

	slotFileName := "[20160331-073917]_[FDB-CC10316_3bc60539-68df-4808-b11c-2e64bc4fe6fd]_DXB Summer S16 (27Mar-29Oct16)_OBFUSCATED.xlsx"

	slotFileName = slotNameExtract ( slotFileName )

	fmt.Println ( "Filename:", slotFileName )

	fmt.Println ( "End Test" )

}

func slotNameExtract ( slotFileName string ) ( string ) {

	fileMetadata := strings.Split ( slotFileName, "_" )

//	for _, value := range fileMetadata { fmt.Println ( value ) }
	return fileMetadata[3]
}
