package protocol

// 0xxx: System commands
// 1xxx: Database and collection commands
// 2xxx: Key-value engine commands

const (
	CommandPing                      uint16 = 1
	CommandSupportedProtocolVersions uint16 = 2
	CommandLogin                     uint16 = 100
	CommandAuthStatus                uint16 = 101

	CommandKVGet uint16 = 2000
)
