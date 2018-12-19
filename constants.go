package dnsPacket

/*
0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |

*/

//QR
const (
	DNSQuery    = 0 << 15
	DNSResponse = 1 << 15
)

//Opcodes
const (
	OpcodeStandardQuery = 0
	OpcodeInverseQuery  = 1
	OpcodeServerStatus  = 2
)

//Rcodes
const (
	RcodeNoError        = 0
	RcodeFormatError    = 1
	RcodeServerFailure  = 2
	RcodeNameError      = 3
	RcodeNotImplemented = 4
	RcodeRefused        = 5
	RcodeMask           = 15
)

//Qclass
const (
	QclassIN = 1
)

//query params flags
const (
	FlagsOpCodeStandardQuery = 0 << 11
	FlagsOpCodeInverseQuery  = 1 << 11
	FlagsOpCodeServerStatus  = 1 << 12
	FlagsAuthoritativeAnswer = 1 << 10
	FlagsTruncation          = 1 << 9
	FlagsRecurionDesired     = 1 << 8
	FlagsRecursionAvailable  = 1 << 7
	FlagsMask                = 3840
	ZMask                    = 112
)
