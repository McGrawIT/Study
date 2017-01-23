// Package blobstore provides functions to access the Predix Blobstore
// via its S3 interface. This is currently just a convenience wrapper
// around minio that handles VCAP configuration, and provides simple
// get/put object functionality.
//
// This package only supports vhost style URLs because the AWS S3 client
// (used in the equivalent Java library) fails path style requests with
// an invalid certificate error (*.store.gecis.io != store.gecis.io).
package main

import (
	"errors"
	"io"
	"log"
	"regexp"

	"github.build.ge.com/aviation-predix-common/vcap-support"
	"github.com/minio/minio-go"
	"os"
	//"fmt"
)

// Cloud Foundry environment variable constants
const blobstoreVCAPKey = "predix-blobstore"
const endpointVCAPKey = "url"
const accessKeyVCAPKey = "access_key_id"
const secretVCAPKey = "secret_access_key"

// Package local vars
var endpoint, host, bucket string
var s3Client *minio.Client

// URL regexes
// Using raw strings here to avoid having to escape backslashes
var endpointRegex = regexp.MustCompile(`^(https?)://([^.]+)\.(.+)$`)
var objectRegex = regexp.MustCompile(`^https?://([^.]+)\.([^/]+)/(.+)$`)

// init loads the blobstore connection configuration from VCAP
// This just fails silently if no blobstore is bound. The alternative
// is to panic, and crash the calling service, which may be able to
// continue without blobstore access.
func init() {
	// Get VCAP environment variables
	svc, err := vcap.LoadServices()
	if err != nil {
		log.Println("Error loading vcap services " + err.Error())
		return
	}

	if len(svc[blobstoreVCAPKey]) == 0 {
		log.Println("Error loading Blobstore config. No Blobstore service bound.")
		return
	}

	// Assume only one blobstore is bound
	credentials := svc[blobstoreVCAPKey][0].Credentials

	accessKey := credentials[accessKeyVCAPKey].(string)
	secret := credentials[secretVCAPKey].(string)

	// Parse the endpoint into host and bucket
	endpoint = credentials[endpointVCAPKey].(string)

	matches := endpointRegex.FindStringSubmatch(endpoint)
	if len(matches) != 4 {
		log.Println("Error parsing Blobstore endpoint.")
		return
	}

	// matches[0] is the entire match
	https := true
	if matches[1] == "http" {
		https = false
	}
	bucket = matches[2]
	host = matches[3]

	// Initialise s3 client
	s3Client, err = minio.NewV2(host, accessKey, secret, https)
	if err != nil {
		log.Println("Failed to create s3 client")
		return
	}
	s3Client.TraceOn(os.Stdout)

}

// PutObject writes a new object to the blobstore and returns
// the URL for the new object.
func PutObject(objectKey string, file io.Reader) (string, error) {

	_, err := s3Client.PutObject(bucket, objectKey, file, "")

	if err != nil {
		return "", err
	}

	return endpoint + "/" + objectKey, nil
}

// GetObject retrieves an object from the blobstore and returns
// it as a byte slice
// TODO It would be better to return an io.Reader here, but there
// isn't an easy way to do that without exposing Minio-implementation
// specifics.
func GetObject(objectURL string) ([]byte, error) {

	// Parse object URL
	matches := objectRegex.FindStringSubmatch(objectURL)
	if len(matches) != 4 {
		message := "Error parsing object URL."
		log.Println(message)
		return nil, errors.New(message)
	}

	// Disabling host check for now to allow path style requests
	/*
	// matches[0] is the entire match
	objBucket := matches[1]
	objHost := matches[2]
	objectKey := matches[3]
	if objBucket != bucket || objHost != host {
		return nil, errors.New("Requested object " + objectURL + " is not stored in the bound blobstore.")
	}
	*/

	// Get the object
	objectKey := matches[3]
	obj, err := s3Client.GetObject(bucket, objectKey)
	if err != nil {
		log.Println("Error on s3Client.GetObject: ", err)
		return nil, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size)
	_, err = obj.Read(buf)

	//if err != nil {
	//	return nil, err
	//}


	return buf, nil

}
