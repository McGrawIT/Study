package binding

/*
#cgo CFLAGS: -I ${SRCDIR}/solclient/include -I ${SRCDIR}/solclient/ex
#cgo LDFLAGS: -L ${SRCDIR}/solclient/lib -lsolclient
#include "solclient.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/jackmanlabs/errors"
	"unsafe"
)

// This function gets the Solace C API version from the library.
func Version() string {
	var ver_ *C.char = C.version()
	defer C.free(unsafe.Pointer(ver_))
	ver := C.GoString(ver_)
	return ver
}

// The Solace Client library needs to be initialized before any other API calls.
func Initialize() error {
	var rc C.int = C.initialize()

	if rc != SOLCLIENT_OK {
		return errors.New("Solace Client initialization failed.")
	}

	return nil
}

//
func EventCallback() {

}


func MessageReceiveCallback(C.SolaceAttachment){

}