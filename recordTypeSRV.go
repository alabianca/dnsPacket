package dnsPacket

type RecordTypeSRV struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string
}

// 2 bytes   2 bytes    2 bytes length prefixed labels
//Priority  | Weight  | Port   | Target
//
func (record *RecordTypeSRV) Process(a Answer) {

	priority := decodePart(a.Data, 0, 2)
	weight := decodePart(a.Data, 2, 4)
	port := decodePart(a.Data, 4, 6)

	target, _ := decodeQname(a.Data[6:])

	record.Priority = priority
	record.Weight = weight
	record.Port = port
	record.Target = target
}

func (record *RecordTypeSRV) Encode(a *Answer) []byte {
	data := make([]byte, 0)

	priorityBytes, _ := fromIntToBytes(record.Priority)
	weightBytes, _ := fromIntToBytes(record.Weight)
	portBytes, _ := fromIntToBytes(record.Port)

	data = append(data, priorityBytes...)
	data = append(data, weightBytes...)
	data = append(data, portBytes...)
	data = append(data, encodeQname(record.Target)...)

	return data
}
