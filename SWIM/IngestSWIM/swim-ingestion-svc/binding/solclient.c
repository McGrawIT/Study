#include "solclient/solClient.h"
#include "_cgo_export.h"
#include <stdlib.h>

#define BUFFER_SIZE 1048576 // 1MB

struct SolaceAttachment {
solClient_opaquePointer_pt* Pointer
solClient_uint32_t* Length
}

/**
 * This function prints C API version to STDOUT.
 * Borrowed from common.c in the Solace examples folder.
 */
char* version(void)
{
    solClient_version_info_pt version_p;
    char *str = malloc(1024);

    solClient_returnCode_t rc = SOLCLIENT_OK;
    rc = solClient_version_get ( &version_p );
    if (rc != SOLCLIENT_OK ) {
        sprintf(str, "Unknown library version, solClient_version_get returns FAIL\n\n");
    } else {
        sprintf(str, "CCSMP Version %s (%s)\tVariant: %s\n\n", version_p->version_p, version_p->dateTime_p, version_p->variant_p);
    }

    // printf(str);

    return str;
}

/**
 * solClient needs to be initialized before any other API calls.
 */
int initialize(){
    solClient_returnCode_t rc = SOLCLIENT_OK;
    rc = solClient_initialize(SOLCLIENT_LOG_DEFAULT_FILTER, NULL);
    return rc;
}

solClient_rxMsgCallback_returnCode_t
messageReceiveCallback ( solClient_opaqueSession_pt opaqueSession_p, solClient_opaqueMsg_pt msg_p, void *user_p )
{
    solClient_returnCode_t rc = SOLCLIENT_OK;

    solClient_opaquePointer_pt attachmentPointer
    solClient_uint32_t *attachmentLength

    rc = solClient_msg_getBinaryAttachmentPtr(msg_p, attachmentPointer, attachmentLength)
    if (rc != SOLCLIENT_OK){
        return rc
    }

    printf ( "Received message:\n" );
    solClient_msg_dump ( msg_p, NULL, 0 );
    printf ( "\n" );


    SolaceAttachment attachment
    attachment.Pointer = attachmentPointer
    attachment.Length = attachmentLength

    MessageReceiveCallback(attachment)




    return SOLCLIENT_CALLBACK_OK;
}

void
eventCallback ( solClient_opaqueSession_pt opaqueSession_p,
                solClient_session_eventCallbackInfo_pt eventInfo_p, void *user_p )
{
    printf("Session EventCallback() called:  %s\n", solClient_session_eventToString ( eventInfo_p->sessionEvent));
}