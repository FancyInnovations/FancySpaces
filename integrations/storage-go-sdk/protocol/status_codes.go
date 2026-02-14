package protocol

// 0xxx: Success codes
// 1xxx: Client error codes
// 2xxx: Server error codes

const (
	// StatusOK indicates that the operation was successful.
	StatusOK uint16 = 1

	// StatusInvalidMessage indicates that the message is invalid or malformed.
	StatusInvalidMessage uint16 = 1000

	// StatusCommandNotFound indicates that the requested command was not found.
	StatusCommandNotFound uint16 = 1001

	// StatusInvalidCredentials indicates that the provided credentials are invalid.
	StatusInvalidCredentials uint16 = 1002

	// StatusBadRequest indicates that the request could not be understood by the server due to incorrect syntax.
	StatusBadRequest uint16 = 1003

	// StatusUnauthorized indicates that authentication is required.
	StatusUnauthorized uint16 = 1004

	// StatusDatabaseNotFound indicates that the specified database was not found.
	StatusDatabaseNotFound uint16 = 1005

	// StatusCollectionNotFound indicates that the specified collection was not found.
	StatusCollectionNotFound uint16 = 1006

	// StatusCommandNotAllowed indicates that the command is not allowed in the current context.
	StatusCommandNotAllowed uint16 = 1007

	// StatusNotFound indicates that the requested resource was not found.
	StatusNotFound uint16 = 1008

	// StatusForbidden indicates that the client does not have permission to access the requested resource.
	StatusForbidden uint16 = 1009

	// StatusInternalServerError indicates that an internal server error occurred.
	StatusInternalServerError uint16 = 2000
)
