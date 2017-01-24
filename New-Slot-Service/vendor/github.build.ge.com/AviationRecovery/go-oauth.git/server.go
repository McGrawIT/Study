package OAuth2

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"log"
	"io/ioutil"
	"io"
)

/*
SERVICE_CREDENTIAL is found by the following command:
	echo -n {CLIENT_ID}:{CLIENT_SECRET} | base64
  example:
	echo -n flydubai-service-client:iM2NrWVd | base64
*/

const (
	DEBUG		= false
	OLD_DEBUG	= false

	ENABLE_FAST_TOKEN_CHECK = true
)

const VERSION = "0.0.9"

var (
	SCOPES = []string{""}
)

var CLIENT_ID string
var CLIENT_SECRET string
var SERVICE_CREDENTIAL string
var uaaHostName string

func init() {

	oauthKeyCache = make(map[string]string)

	cihn := "uaa_client_id"
	cshn := "uaa_client_secret"
	hn := "uaa_url"

	uaa_client_id := GetServiceHostName(cihn)
	uaa_client_secret := GetServiceHostName(cshn)
	uaaHostName = GetServiceHostName(hn)

	if uaa_client_id != "" && !strings.ContainsAny(uaa_client_id, "[]") {
		CLIENT_ID = uaa_client_id
	}
	if uaa_client_secret != "" && !strings.ContainsAny(uaa_client_secret, "[]") {
		CLIENT_SECRET = uaa_client_secret
	}
	if uaaHostName == "" {
		log.Println("Error: No uaa_uri specified in Environment")
	}

	//fmt.Println("DBG-> CLIENT_ID: ",CLIENT_ID,"; CLIENT_SECRET: ",CLIENT_SECRET)

	//  calculate SERVICE_CREDENTIAL
	message := CLIENT_ID + ":" + CLIENT_SECRET
	uaa_service_credential := base64.StdEncoding.EncodeToString([]byte(message))
	SERVICE_CREDENTIAL = uaa_service_credential

}

func GetVersion() (version string) {
	return VERSION
}

func ParseTokenRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")[7:]
	return token
}

func GetAuthenticationToken(r *http.Request) (token string, exists bool) {

	exists = true
	if len(r.Header.Get("Authorization")) == 0 {
		exists = false
		return
	}

	token = ParseTokenRequest(r)
	return
}

func CheckAuthentication(r *http.Request) (statusCode int) {

	host_url := uaaHostName + "/oauth/check_token"
	if uaaHostName == "" {
		statusCode = http.StatusUnauthorized //  401
		log.Println("Error: No uaa_uri specified in Environment")
		return
	}

	token, exists := GetAuthenticationToken(r)
	if !exists {
		statusCode = http.StatusUnauthorized //  401
		return
	}

	//  check for FastToken processing here??
	if ENABLE_FAST_TOKEN_CHECK {
		if FastTokenVerify(token) {
			if DEBUG {
				fmt.Println("OAuth2_DBG-> Verified Request Authorization using FastToken")
			}
			statusCode = http.StatusOK
			return
		}
	}

	data := url.Values{}
	data.Set("token", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", host_url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Basic "+SERVICE_CREDENTIAL)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if (resp != nil) { defer resp.Body.Close() }
	if resp == nil || err != nil {
		statusCode = http.StatusUnauthorized
		return
	}

	//	fmt.Println("DBG-> resp: ", resp)
	statusCode = resp.StatusCode
	if statusCode != http.StatusOK {
		log.Println("OAuth2_DBG-> Verify Request Authorization Failed for token (",token,"); Code: ", statusCode)
		return
	}

	//  need to process body here
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		log.Println("Error-CheckAuthentication(): Unable to read token validation body")
////		statusCode = http.StatusUnauthorized
//		return
	}

	if DEBUG { fmt.Println("DBG-> body: ", string(body)) }

	if DEBUG { fmt.Println("OAuth2_DBG-> Verified Request Authorization") }
	return
}

func IsAuthenticated() negroni.Handler {
	au := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		statusCode := CheckAuthentication(r)
		if statusCode != http.StatusOK {
			//Handle the different response codes appropriately
			w.WriteHeader(statusCode)

			ct := r.Header.Get("Content-Type")
			if strings.Contains(strings.ToLower(ct), "json") {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				fmt.Fprintf(w, "{ \"authentication-status-code\" : "+strconv.Itoa(statusCode)+" }")
			} else {
				w.Header().Set("Content-Type", "plain/text")
				fmt.Fprintf(w, "Authentication Status Code: "+strconv.Itoa(statusCode))
			}
			return
		}

		next(w, r)
	}
	return negroni.HandlerFunc(au)
}
