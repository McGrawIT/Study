package binding

import "unsafe"

/*
#cgo CFLAGS: -I ${SRCDIR}/solclient/include -I ${SRCDIR}/solclient/ex
#cgo LDFLAGS: -L ${SRCDIR}/solclient/lib -lsolclient
#include "solclient.h"
#include <stdlib.h>
*/

/*	The question is the "include" that gets this C code to "connect" to the Solace APISs

	What follows is the C code / functions I will call from Go
 */

//	#include <stdio.h>
//
import "C"
import (
	"github.com/jackmanlabs/errors"
	"unsafe"
)

const (

	SOLCLIENT_SESSION_PROP_AUTHENTICATION_SCHEME_BASIC					= 2
	SOLCLIENT_SESSION_PROP_AUTHENTICATION_SCHEME_CLIENT_CERTIFICATE		= 3
	SOLCLIENT_SESSION_PROP_AUTHENTICATION_SCHEME_GSS_KRB				= 4
	OLCLIENT_SESSION_PROP_USERNAME = "jr.mcgraw"
	LCLIENT_SESSION_PROP_PASSWORD = "secret"
)

/*
	To create a Session, the client application must provide the following:

	Session Properties

	Properties used to customize the Session. Any Session property that is not explicitly supplied
	is set to default values. Although the defaults can be used in many cases, some client and router
	parameters require specific input from the client to establish a connection to a Message VPN
	on a Solace router. Refer to Session Properties Required to Establish a Connection.

	Specific Context Instance

	Session Event Callback.  A Session Event Callback must be specified when creating a Session.
	This Callback is invoked for each Session Event

	Message Receive Callback.  A message event callback must be specified when creating a Session.
	This callback is invoked each time a Direct message is received through the Session.

	Once a Session is created, it must be connected.
 */
func SessionCreateConnect () {


	solClient_session_create()
	solClient_session_connect()
}

func SessionClose () {

	solClient_session_disconnect()

//	If the Session object is not destroyed, you may reconnect to it at a later time

	solClient_session_destroy()				// Not Required
}
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

/*
	When a Session is created, the application must provide a Session event callback

	( solClient_session_eventCallbackInfoFunc_t ),

	along with an optional pointer to client data. This callback routine is invoked for router events
	that occur for the Session, such as connection problems, or publish or subscription issues.

	For a complete list of possible Session events,
	refer to the C API Developer Online Reference documentation.

	For an example of how to configure a Session event callback refer to the common.c sample file
 */

func SessionEventCallback () {			// See common.c sample file

}