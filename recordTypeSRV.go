package dnsPacket

import "fmt"

type RecordTypeSRV struct {
	Data string
}

func (record *RecordTypeSRV) Process(a Answer) {
	fmt.Println(a.Data)

	record.Data = string(a.Data)
}
