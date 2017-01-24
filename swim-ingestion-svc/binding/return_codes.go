package binding

const (
	SOLCLIENT_OK          = 0  /**< The API call was successful. */
	SOLCLIENT_WOULD_BLOCK = 1  /**< The API call would block, but non-blocking was requested. */
	SOLCLIENT_IN_PROGRESS = 2  /**< An API call is in progress (non-blocking mode). */
	SOLCLIENT_NOT_READY   = 3  /**< The API could not complete as an object is not ready (for example, the Session is not connected). */
	SOLCLIENT_EOS         = 4  /**< A getNext on a structured container returned End-of-Stream. */
	SOLCLIENT_NOT_FOUND   = 5  /**< A get for a named field in a MAP was not found in the MAP. */
	SOLCLIENT_NOEVENT     = 6  /**< solClient_context_processEventsWait returns this if wait is zero and there is no event to process */
	SOLCLIENT_INCOMPLETE  = 7  /**< The API call completed some, but not all, of the requested function. */
	SOLCLIENT_ROLLBACK    = 8  /**< solClient_transactedSession_commit returns this when the transaction has been rolled back. */
	SOLCLIENT_FAIL        = -1 /**< The API call failed. */
)
