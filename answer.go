package dnsPacket

type Answer struct {
	Name     string
	Type     string
	Class    string
	Ttl      int
	RdLength int
	Data     []byte
}
