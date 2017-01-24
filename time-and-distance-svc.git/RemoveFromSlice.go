package main

import (
		"fmt"
)

var	list	= []int{1,2,3,4,5}

func deleteOne ( member int, members []int ) ( []int ) {

	result := members
	for pos, val := range members {

		if val == member {

			resultLeft := append ( []int{}, members[0:pos]... )
			pos++
			//fmt.Println ( resultLeft )
			result = append ( resultLeft, members[pos:len(members)]... )

		} else { continue }
	}

	return result
}



func main() {

	fmt.Println(list)
	for i:=1;i<8;i++ {

		list1 := deleteOne ( i, list )
		fmt.Println ("Remove", i, "from", list, "==>", list1 )
	}
}
