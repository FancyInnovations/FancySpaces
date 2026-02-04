package protocol

// 0xxx: Success codes
// 1xxx: Client error codes
// 2xxx: Server error codes

const (
	// StatusOK indicates that the operation was successful.
	StatusOK uint16 = 0001

	// StatusInvalidMessage indicates that the message is invalid or malformed.
	StatusInvalidMessage uint16 = 1000

	// StatusCommandNotFound indicates that the requested command was not found.
	StatusCommandNotFound uint16 = 1001

	// StatusInternalServerError indicates that an internal server error occurred.
	StatusInternalServerError uint16 = 2000
)
