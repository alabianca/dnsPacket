package dnsPacket

type RecordTypeDefault struct {
	Data []byte
}

func (record *RecordTypeDefault) Process(a Answer) {
	record.Data = a.Data
}

func (record *RecordTypeDefault) Encode(a *Answer) []byte {
	return []byte{}
}
