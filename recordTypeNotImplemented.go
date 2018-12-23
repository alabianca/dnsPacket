package dnsPacket

type RecordTypeDefault struct {
	Data []byte
}

func (record *RecordTypeDefault) Process(a Answer) {
	record.Data = a.Data
}

func (record *RecordTypeDefault) Type() int {
	return 0
}

func (record *RecordTypeDefault) Encode() []byte {
	return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}
