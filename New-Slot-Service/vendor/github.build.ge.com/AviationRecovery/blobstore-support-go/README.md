# blobstore
--
    import "github.build.ge.com/aviation-predix-common/blobstore-support-go.git"

Package blobstore provides functions to access the Predix Blobstore via its S3
interface. This is currently just a convenience wrapper around minio that
handles VCAP configuration, and provides simple get/put object functionality.

This package only supports vhost style URLs because the AWS S3 client (used in
the equivalent Java library) fails path style requests with an invalid
certificate error (*.store.gecis.io != store.gecis.io).

## Usage

#### func  GetObject

```go
func GetObject(objectURL string) ([]byte, error)
```
GetObject retrieves an object from the blobstore and returns it as a byte slice
TODO It would be better to return an io.Reader here, but there isn't an easy way
to do that without exposing Minio-implementation specifics.

#### func  PutObject

```go
func PutObject(objectKey string, file io.Reader) (string, error)
```
PutObject writes a new object to the blobstore and returns the URL for the new
object.
