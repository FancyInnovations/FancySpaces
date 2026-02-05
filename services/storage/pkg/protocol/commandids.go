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

	CommandKVDelete         uint16 = 2020
	CommandKVDeleteMultiple uint16 = 2021
	CommandKVDeleteAll      uint16 = 2022

	CommandKVExists      uint16 = 2030
	CommandKVGet         uint16 = 2031
	CommandKVGetMultiple uint16 = 2032
	CommandKVGetAll      uint16 = 2033
	CommandKVKeys        uint16 = 2034

	CommandKVNumIncrement  uint16 = 2520
	CommandKVNumDecrement  uint16 = 2521
	CommandKVNumMultiply   uint16 = 2522
	CommandKVNumDivide     uint16 = 2523
	CommandKVNumModulo     uint16 = 2524
	CommandKVNumLeftShift  uint16 = 2525
	CommandKVNumRightShift uint16 = 2526
	CommandKVNumBitwiseAnd uint16 = 2527
	CommandKVNumBitwiseOr  uint16 = 2528
	CommandKVNumBitwiseXor uint16 = 2529
	CommandKVNumBitwiseNot uint16 = 2530

	CommandStringAppend    uint16 = 2540
	CommandStringPrepend   uint16 = 2541
	CommandStringLength    uint16 = 2542
	CommandStringSubstring uint16 = 2543

	CommandKVListLength    uint16 = 2550
	CommandKVListGet       uint16 = 2551
	CommandKVListSet       uint16 = 2552
	CommandKVListRemove    uint16 = 2553
	CommandKVListLeftPush  uint16 = 2554
	CommandKVListRightPush uint16 = 2555
	CommandKVListLeftPop   uint16 = 2556
	CommandKVListRightPop  uint16 = 2557

	CommandKVMapLength uint16 = 2560
	CommandKVMapSet    uint16 = 2561
	CommandKVMapGet    uint16 = 2562
	CommandKVMapDelete uint16 = 2563
	CommandKVMapExists uint16 = 2564
	CommandKVMapKeys   uint16 = 2565
	CommandKVMapValues uint16 = 2566
)
