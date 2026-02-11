package protocol

// 0xxx: System commands
// 1xxx: Database and collection commands
// 2xxx: Key-value engine commands
// 3xxx: Document engine commands
// 4xxx: Object engine commands
// 5xxx: Analytical engine commands
// 6xxx: Broker engine commands

const (
	ServerCommandPing                      uint16 = 1
	ServerCommandSupportedProtocolVersions uint16 = 2
	ServerCommandLogin                     uint16 = 100
	ServerCommandAuthStatus                uint16 = 101

	ServerCommandKVSet               uint16 = 2000
	ServerCommandKVSetTTL            uint16 = 2001
	ServerCommandKVSetMultiple       uint16 = 2002
	ServerCommandKVSetMultipleTTL    uint16 = 2003
	ServerCommandKVSetIfExists       uint16 = 2004
	ServerCommandKVSetIfExistsTTL    uint16 = 2005
	ServerCommandKVSetIfNotExists    uint16 = 2006
	ServerCommandKVSetIfNotExistsTTL uint16 = 2007

	ServerCommandKVDelete         uint16 = 2020
	ServerCommandKVDeleteMultiple uint16 = 2021
	ServerCommandKVDeleteAll      uint16 = 2022

	ServerCommandKVExists         uint16 = 2030
	ServerCommandKVGet            uint16 = 2031
	ServerCommandKVGetMultiple    uint16 = 2032
	ServerCommandKVGetAll         uint16 = 2033
	ServerCommandKVGetTTL         uint16 = 2034
	ServerCommandKVGetMultipleTTL uint16 = 2035
	ServerCommandKVGetAllTTL      uint16 = 2036
	ServerCommandKVKeys           uint16 = 2037
	ServerCommandKVCount          uint16 = 2038
	ServerCommandKVSize           uint16 = 2039

	ServerCommandKVNumIncrement  uint16 = 2520
	ServerCommandKVNumDecrement  uint16 = 2521
	ServerCommandKVNumMultiply   uint16 = 2522
	ServerCommandKVNumDivide     uint16 = 2523
	ServerCommandKVNumModulo     uint16 = 2524
	ServerCommandKVNumLeftShift  uint16 = 2525
	ServerCommandKVNumRightShift uint16 = 2526
	ServerCommandKVNumBitwiseAnd uint16 = 2527
	ServerCommandKVNumBitwiseOr  uint16 = 2528
	ServerCommandKVNumBitwiseXor uint16 = 2529
	ServerCommandKVNumBitwiseNot uint16 = 2530

	ServerCommandStringAppend    uint16 = 2540
	ServerCommandStringPrepend   uint16 = 2541
	ServerCommandStringLength    uint16 = 2542
	ServerCommandStringSubstring uint16 = 2543

	ServerCommandKVListLength    uint16 = 2550
	ServerCommandKVListGet       uint16 = 2551
	ServerCommandKVListSet       uint16 = 2552
	ServerCommandKVListRemove    uint16 = 2553
	ServerCommandKVListLeftPush  uint16 = 2554
	ServerCommandKVListRightPush uint16 = 2555
	ServerCommandKVListLeftPop   uint16 = 2556
	ServerCommandKVListRightPop  uint16 = 2557

	ServerCommandKVMapLength uint16 = 2560
	ServerCommandKVMapSet    uint16 = 2561
	ServerCommandKVMapGet    uint16 = 2562
	ServerCommandKVMapDelete uint16 = 2563
	ServerCommandKVMapExists uint16 = 2564
	ServerCommandKVMapKeys   uint16 = 2565
	ServerCommandKVMapValues uint16 = 2566

	ServerCommandObjectPut         uint16 = 4000
	ServerCommandObjectGet         uint16 = 4001
	ServerCommandObjectGetMetadata uint16 = 4002
	ServerCommandObjectDelete      uint16 = 4003
	ServerCommandObjectExists      uint16 = 4004
	ServerCommandObjectList        uint16 = 4005
	ServerCommandObjectCopy        uint16 = 4006
	ServerCommandObjectMove        uint16 = 4007
	ServerCommandObjectRename      uint16 = 4008

	ServerCommandBrokerSubscribe      uint16 = 6000
	ServerCommandBrokerSubscribeQueue uint16 = 6001
	ServerCommandBrokerUnsubscribe    uint16 = 6002
	ServerCommandBrokerPublish        uint16 = 6003
	ClientCommandBrokerMessage        uint16 = 6004
)
