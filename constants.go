package dnsPacket

const (
	FlagsQuery               = 0 << 15
	FlagsResponse            = 1 << 15
	FlagsOpCodeStandardQuery = 0 << 11
	FlagsOpCodeInverseQuery  = 1 << 11
	FlagsOpCodeServerStatus  = 1 << 12
	FlagsAuthoritativeAnswer = 1 << 10
	FlagsTruncation          = 1 << 9
	FlagsRecurionDesired     = 1 << 8
	FlagsRecursionAvailable  = 1 << 7
)
