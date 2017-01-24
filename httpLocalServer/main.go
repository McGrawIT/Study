package main

import (
	"github.build.ge.com/502612370/httpLocalServer/api"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting")
	http.ListenAndServe(":3000", api.Handlers())
}
