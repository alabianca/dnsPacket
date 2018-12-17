package dnsPacket

import (
	"bytes"
	"encoding/binary"
	"strings"
)

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

type Question struct {
	Qname  string
	Qtype  string
	Qclass string
}

type DNSPacket struct {
	Type      string
	ID        int16
	Opcode    int
	Flags     int
	Qdcount   int16
	Ancount   int16
	Nscount   int16
	Arcount   int16
	Questions []Question
}

func (dns *DNSPacket) AddQuestion(q Question) {
	dns.Questions = append(dns.Questions, q)
}

func Encode(dnsPacket *DNSPacket) []byte {
	packet := make([]byte, 0)
	isQuery := dnsPacket.Type == "query"
	var packetType = 0
	if isQuery {
		packetType = 0
	} else {
		packetType = 1
	}

	packetID, _ := FromIntToBytes(uint16(dnsPacket.ID))
	params := packetType | dnsPacket.Opcode | dnsPacket.Flags
	queryParms, _ := FromIntToBytes(uint16(params))
	qcount, _ := FromIntToBytes(uint16(dnsPacket.Qdcount))
	ancount, _ := FromIntToBytes(uint16(dnsPacket.Ancount))
	nscount, _ := FromIntToBytes(uint16(dnsPacket.Nscount))
	arcount, _ := FromIntToBytes(uint16(dnsPacket.Arcount))

	packet = append(packet, packetID...)
	packet = append(packet, queryParms...)
	packet = append(packet, qcount...)
	packet = append(packet, ancount...)
	packet = append(packet, nscount...)
	packet = append(packet, arcount...)

	for _, q := range dnsPacket.Questions {
		packet = append(packet, encodeQuestion(q)...)
	}

	return packet

}

func FromIntToBytes(num uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, num)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func encodeQname(qname string) []byte {
	name := make([]byte, 0)

	sections := strings.Split(qname, ".")

	for i := 0; i < len(sections); i++ {
		length := len(sections[i])

		name = append(name, byte(length))

		for j := 0; j < length; j++ {
			name = append(name, byte(sections[i][j]))
		}
	}
	name = append(name, byte(0))
	return name
}

func encodeQuestion(q Question) []byte {
	question := make([]byte, 0)

	name := encodeQname(q.Qname)
	qtype, _ := FromIntToBytes(uint16(1))  //hard coded
	qclass, _ := FromIntToBytes(uint16(1)) //hard coded

	question = append(question, name...)
	question = append(question, qtype...)
	question = append(question, qclass...)

	return question
}
