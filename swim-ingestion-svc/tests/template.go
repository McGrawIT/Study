package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	basepath = flag.String("basepath", "", "The base path of the configuration manager REST API that you want to test.")
)

func main() {

	flag.Parse()
	if *basepath == "" {
		fmt.Println("You must specify a basepath for this tool to work.")
		flag.Usage()
		os.Exit(1)
	}

}
