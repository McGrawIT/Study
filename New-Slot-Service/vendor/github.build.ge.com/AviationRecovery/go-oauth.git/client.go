package OAuth2

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const CacheHysteresis = 10

var cachedClient *http.Client
var cachedTS = int64(0)

func RequestAuthorization() (client *http.Client) {

	if DEBUG {
		fmt.Println("OAuth2_DBG-> Including Authorization in Request")
	}

	if (cachedClient != nil) && (time.Now().Unix() <= cachedTS) {
		client = cachedClient
		return
	}

	if DEBUG {
		fmt.Println("OAuth2_DBG-> Requesting New Token")
	}

	hn := "uaa_url"
	host_name := GetServiceHostName(hn)
	host_url := host_name + "/oauth/token"
	if host_name == "" {
		log.Println("Error: No uaa_uri specified in Environment")
		return
	}

	config := clientcredentials.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		TokenURL:     host_url,
		Scopes:       SCOPES,
	}

	// GET TOKEN TO DECODE
	t := &oauth2.Token{}
	t, _ = config.Token(context.Background())
	if OLD_DEBUG {
		fmt.Println("OAuth2_DBG-> ADF-1 API Token:", t.AccessToken)

		result := strings.Split(t.AccessToken, ".")
		d, _ := base64.StdEncoding.DecodeString(result[1])
		fmt.Println("OAuth2_DBG-> Decoded ADF-1 API JWT: ", string(d))
	}
	claims, _ := decode(t.AccessToken)
	cachedTS = claims.Exp - CacheHysteresis			//  give us a hysteresis window

	// GENERATE NEW REST CLIENT TO GET DATA - the client will update its token if it's expired
	client = config.Client(context.Background())
	cachedClient = client

	return
}
