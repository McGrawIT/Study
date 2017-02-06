package binding

import "unsafe"

/*
#cgo CFLAGS: -I ${SRCDIR}/solclient/include -I ${SRCDIR}/solclient/ex
#cgo LDFLAGS: -L ${SRCDIR}/solclient/lib -lsolclient
#include "solclient.h"
#include <stdlib.h>


	int solClient_session_create();
	int solClient_session_connect();
	tin solClient_session_disconnect();
	int solClient_session_destroy();

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

/*	Creates a new Session within a specified Context.

	solClient_dllExport solClient_returnCode_t solClient_session_create	(

		solClient_propertyArray_pt 				props,
		solClient_opaqueContext_pt 				opaqueContext_p,
		solClient_opaqueSession_pt * 			opaqueSession_p,
		solClient_session_createFuncInfo_t * 	funcInfo_p,
		size_t 									funcInfoSize
	)

	When processing the property list, the API will not modify any of the strings pointed to by props.
	Parameters:

	props 				An array of name/value string pair pointers to configure session properties.
	opaqueContext_p 	The Context in which the Session is to be created.
	opaqueSession_p 	An opaque Session pointer is returned that refers to the created Session.
	funcInfo_p 			A pointer to a structure that provides information on callback functions for events and received messages.
	funcInfoSize 		The size (in bytes) of the passed-in funcInfo structure to allow the structure to grow in the future.

	The session properties are supplied as an array of name/value pointer pairs, where the name and
	value are both strings. Only configuration property names starting with "SESSION_" are processed;
	other property names are ignored. Any values not supplied are set to default values. When the
	Session is created, an opaque Session pointer is returned to the caller, and this value is then
	used for any Session-level operations (for example, sending a message). The passed-in structure
	functInfo_p provides information on the message receive callback function and the Session event
	function which the application has provided for this Session. Both of these callbacks are mandatory.
	The message receive callback is invoked for each received message on this Session. The Session event
	callback is invoked when Session events occur, such as the Session going up or down. Both callbacks
	are invoked in the context of the Context thread to which this Session belongs. Note that the property
	values are stored internally in the API and the caller does not have to maintain the props array
	or the strings that are pointed to after this call completes.

	Returns:

	SOLCLIENT_OK, SOLCLIENT_FAIL
	SubCodes (Unless otherwise noted above, subcodes are only relevant when this function returns SOLCLIENT_FAIL):
	SOLCLIENT_SUBCODE_OUT_OF_RESOURCES - The maximum number of Sessions already created for Context (refer to SOLCLIENT_CONTEXT_PROP_MAX_SESSIONS).

 */
func SessionCreateConnect ( propertyArray_pt 			props,
opaqueContext_pt 			opaqueContext_p,
opaqueSession_pt * 			opaqueSession_p,
session_createFuncInfo_t * 	funcInfo_p,
size_t 						funcInfoSize ) {

	status := solClient_session_create( propertyArray_pt, opaqueContext_pt, opaqueSession_pt, session_createFuncInfo_t, size_t )

	if status != SOLCLIENT_OK {

	}

/*
	Define any "flows" for the Session  ( Does this define message flows? )
 */
	status = SessionFlow()


//	Session Created; Connect to it

	solClient_session_connect()
}

/*
	Flow characteristics and behavior are defined by Flow properties. The Flow properties are supplied
	as an array of name/value pointer pairs, where the name and value are both strings.

	FLOW and ENDPOINT configuration property names are processed; other property names are ignored.
	If the Flow creation specifies a non-durable endpoint, ENDPOINT properties can be used to change the
	default properties on the non-durable endpoint. Any values not supplied are set to default values.

	When the Flow is created, an opaque Flow pointer is returned to the caller, and this value is then
	used for any Flow-level operations (for example, starting/stopping a Flow, getting statistics,
	sending an acknowledgment). The passed-in structure functInfo_p provides information on the message
	receive callback function and the Flow event function which the application has provided for this Flow.
	Both of these callbacks are mandatory. The message receive callback is invoked for each received message
	on this Flow. The Flow event callback is invoked when Flow events occur, such as the Flow going up or down.
	Both callbacks are invoked in the context of the Context thread to which the controlling Session belongs.
 */

/*
	solClient_dllExport solClient_returnCode_t solClient_session_createFlow	(

		solClient_propertyArray_pt 			props,
		solClient_opaqueSession_pt 			opaqueSession_p,
		solClient_opaqueFlow_pt * 			opaqueFlow_p,
		solClient_flow_createFuncInfo_t * 	funcInfo_p,
		size_t 								funcInfoSize
	)


	props 				An array of name/value string pair pointers to configure Flow properties.
	opaqueSession_p 	The Session in which the Flow is to be created.
	opaqueFlow_p 		The returned opaque Flow pointer that refers to the created Flow.
	funcInfo_p 			A pointer to a structure that provides information on callback functions for events and received messages.
	funcInfoSize 		The size of the passed-in funcInfo structure (in bytes) to allow the structure to grow in the future.

 */

func SessionFlow () {

	status := solClient_session_createFlow	()

	if status != SOLCLIENT_OK {

	}
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