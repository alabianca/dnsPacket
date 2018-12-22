package dnsPacket

import (
	"bytes"
	"strconv"
	"strings"
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

func (record *RecordTypeA) Encode() []byte {
	return encodeIpV4(record.IPv4)
}

func encodeIpV4(ip string) []byte {
	byteIp := make([]byte, 4)
	parts := strings.Split(ip, ".")

	for i := range parts {
		part, _ := strconv.ParseInt(parts[i], 10, 16)
		byteIp[i] = byte(part)
	}

	return byteIp
}

//127.0.0.1
