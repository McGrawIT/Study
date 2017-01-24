package slot_svc
//package main

import (
	"fmt"
	"strconv"
)

type	testPointer	 struct {

	f1		string
	f2		string
	f3		int
}


func learning() {
//func main () {

	y := testPointer{}
	z := testPointer{}

	x := y.test1 ( "Boston", "Red Sox", 2007 )
	fmt.Println ( x )


	res := test2 ( z )
	fmt.Println( res )

	v3 := test3 ( z )
	fmt.Println( *v3 )

	m := message{}

	msg := m.test4 ( "Boston", "Red Sox", 2007 )
	fmt.Println( msg )
}

func ( testPointer ) test1 ( f1, f2 string, f3 int ) ( string ) {

	result := f1 + f2 + strconv.Itoa( f3 )

	return result
}

func test2 ( val testPointer ) ( testPointer ) {

	val.f1 = "Hello"
	val.f2 = "World."
	val.f3 = 1960

	return val
}



func test3 ( v1 testPointer ) ( *testPointer ) {

	v2 := &v1

	v2.f1 = "Hello, this is"
	v2.f2 = "Water's World."
	v2.f3 = 2016

	return v2
}

type	message 	struct {

	sender		string
	result		string
}

func ( result *message ) test4 ( f1, f2 string, f3 int ) ( message ) {


	//	msg := message{}

	result.sender = "Fan"
	result.result = f1 + f2 + strconv.Itoa( f3 )

	return *result
}

