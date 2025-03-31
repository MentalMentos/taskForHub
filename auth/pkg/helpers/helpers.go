package helpers

const (
	FailedToCreateUser    = "FAILED_TO_CREATE_USER" // FailedToCreateElement contains error message for failed to create element
	FailedToDeleteUser    = "FAILED_TO_DELETE_USER" // FailedToDeleteElement contains error message for failed to delete element
	FailedToGetUser       = "FAILED_TO_GET_USER"    // FailedToGetElements contains error message for failed to get elements
	JSONParseError        = "JSON_PARSE_ERROR"      // JSONParseError contains error message for failed to parse json
	DefaultValueForFields = "NOT FOUND"             // DefaultValueForFields contains default value for fields
	FailedToHashPass      = "FAILED_TO_HASH_PASSWORD"
	FailedToGenJWT        = "FAILED_TO_GEN_JWT"

	AppPrefix      = " [ APP ] "      // AppPrefix contains app prefix
	ResponsePrefix = " [ RESPONSE ] " // ResponsePrefix contains response prefix
	InfoPrefix     = " INFO "         // InfoPrefix contains info prefix
	Success        = "SUCCESS"        // Success contains success message
	RequestError   = "REQUEST_ERROR"  // RequestError contains error message for failed to error
	StatusPrefix   = " STATUS "       // StatusPrefix contains status prefix

	UnmarshalError = "UNMARSHAL_ERROR" // UnmarshalError contains error message for unmarshal error
	ReadBodyError  = "READ_BODY_ERROR" // ReadBodyError contains error message for read body error

	RepoPrefix     = " [ REPOSITORY ] "
	ReconnectDB    = "RECONNECTING TO DATABASE..."
	DisconnectDB   = "DISCONNECTED FROM DATABASE"
	ConnectFailed  = "FAILED TO CONNECT TO DATABASE"
	ConnectSuccess = "SUCCESSFULLY CONNECTED TO MONGO DATABASE"
	FailedToClose  = "FAILED TO CLOSE"
)
