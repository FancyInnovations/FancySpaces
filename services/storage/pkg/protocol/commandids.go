package protocol

// 0xxx: System commands
// 1xxx: Database and collection commands
// 2xxx: Key-value engine commands

const (
	CommandPing                      uint16 = 1
	CommandSupportedProtocolVersions uint16 = 2
	CommandLogin                     uint16 = 100
	CommandAuthStatus                uint16 = 101

	CommandKVSet               uint16 = 2000
	CommandKVSetTTL            uint16 = 2001
	CommandKVSetMultiple       uint16 = 2002
	CommandKVSetMultipleTTL    uint16 = 2003
	CommandKVSetIfExists       uint16 = 2004
	CommandKVSetIfExistsTTL    uint16 = 2005
	CommandKVSetIfNotExists    uint16 = 2006
	CommandKVSetIfNotExistsTTL uint16 = 2007
	CommandKVDelete            uint16 = 2008
	CommandKVExists            uint16 = 2100
	CommandKVGet               uint16 = 2101
	CommandKVGetMultiple       uint16 = 2102
	CommandKVGetAll            uint16 = 2103
	CommandKVKeys              uint16 = 2104
)
