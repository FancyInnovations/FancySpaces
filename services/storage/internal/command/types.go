package command

// 0xxx: System commands
// 1xxx: Database and collection commands
// 2xxx: Key-value engine commands

const (
	CommandPing                      uint16 = 0001
	CommandSupportedProtocolVersions uint16 = 0002
)
