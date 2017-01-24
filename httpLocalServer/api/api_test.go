package api_test

import (
	"fmt"
	"github.build.ge.com/502612370/httpLocalServer/api"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

var (
	server   *httptest.Server
	reader   io.Reader
	usersUrl string
)

func init() {
	server = httptest.NewServer(api.Handlers())

	usersUrl = fmt.Sprintf("%s/users", server.URL)
}

func TestCreateUser(t *testing.T) {

	userJson := `{"username": "dennis", "balance": 200}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", usersUrl, reader)

	res, err := http.DefaultClient.Do(request)

	t.Log( "Response Body:", res.Body )

	if err != nil { t.Error(err) }

	if res.StatusCode != 201 { t.Errorf("Success expected: %d", res.StatusCode) }
}


func TestUniqueUsername(t *testing.T) {

	userJson := `{"username": "dennis", "balance": 200}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", usersUrl, reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil { t.Error(err) }

	if res.StatusCode != 400 { t.Error("Bad Request expected: %d", res.StatusCode) }
}


func TestListUsers(t *testing.T) {

	var resultSet = []api.User{}
	responseBody := "[{0 dennis 200}]"

	//	reader = strings.NewReader("")

	//	request, err := http.NewRequest("GET", usersUrl, reader)

	//	res, err := http.DefaultClient.Do(request)

	//	body, err := ioutil.ReadAll( res.Body )
	//	err = json.Unmarshal( body, &resultSet )

	//	if err != nil { t.Error(err) }

	//	if res.StatusCode != 200 { t.Errorf("Success expected: %d", res.StatusCode) }

	//	userJson := `{"username": "dennis", "balance": 200}`
	//	responseBody := "[{0 dennis 200}]"

	actualResult, errorMessage := genericTest( usersUrl, "GET", "", responseBody, resultSet, http.StatusOK )

	t.Log( "Response Body:", actualResult )
	t.Log( "Expected:", responseBody )

	if errorMessage != "" { t.Error( errorMessage ) }
}

func genericTest ( endPoint,requestType, request, responseBody string, resultType interface{}, statusExpected int ) ( resultSet, errorMessage string ){

	resultSet = resultType

	reader = strings.NewReader( request )

	newRequest, err := http.NewRequest( requestType, endPoint, reader )

	if err != nil { errorMessage = "NewRequest failed" }
	res, _ := http.DefaultClient.Do( newRequest )

	if res.StatusCode != statusExpected { errorMessage = "Expected Status:" + strconv.Itoa( statusExpected ) + " Got: " + strconv.Itoa( res.StatusCode ) }

	body, err := ioutil.ReadAll( res.Body )
	err = json.Unmarshal( body, &resultSet )

	return string ( resultSet ), errorMessage
}
