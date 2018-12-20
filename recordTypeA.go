package dnsPacket

import (
	"bytes"
	"strconv"
)

type RecordTypeA struct {
	IPv4 string
}

func (record *RecordTypeA) Process(a Answer) {
	buf := new(bytes.Buffer)

	//4 bytes (ipv4 address)
	for i := range a.Data {
		buf.WriteString(strconv.Itoa(int(a.Data[i])))

		if i <= 2 {
			buf.WriteString(".")
		}
	}

	record.IPv4 = buf.String()
}

//127.0.0.1
