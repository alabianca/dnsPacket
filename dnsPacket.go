package dnsPacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
+---------------------+
| Header              |
+---------------------+
| Question            | Question for the name server
+---------------------+
| Answer              | Answers to the question
+---------------------+
| Authority           | Not used in this project
+---------------------+
| Additional          | Not used in this project
+---------------------+
*/

/*
DNS HEADERS

0 1  2 3  4 5 6 7 8 9 0 1 2 3 4 5
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| ID                                            |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR| Opcode     |AA|TC|RD|RA| Z | RCODE         |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| QDCOUNT                                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| ANCOUNT                                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| NSCOUNT                                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| ARCOUNT                                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type DNSPacket struct {
	Type       string
	ID         uint16
	Opcode     int
	Rcode      int
	Flags      int
	Qdcount    uint16
	Ancount    uint16
	Nscount    uint16
	Arcount    uint16
	Questions  []Question
	Answers    []Answer
	Authority  []byte
	Additional []byte
}

func (dns *DNSPacket) AddQuestion(q Question) {
	dns.Questions = append(dns.Questions, q)
}

func Encode(dnsPacket *DNSPacket) []byte {
	packet := make([]byte, 0)
	isQuery := dnsPacket.Type == "query"
	var packetType = 0
	if isQuery {
		packetType = DNSQuery
	} else {
		packetType = DNSResponse
	}

	packetID, _ := fromIntToBytes(uint16(dnsPacket.ID))
	params := packetType | dnsPacket.Opcode | dnsPacket.Flags
	queryParms, _ := fromIntToBytes(uint16(params))
	qcount, _ := fromIntToBytes(uint16(dnsPacket.Qdcount))
	ancount, _ := fromIntToBytes(uint16(dnsPacket.Ancount))
	nscount, _ := fromIntToBytes(uint16(dnsPacket.Nscount))
	arcount, _ := fromIntToBytes(uint16(dnsPacket.Arcount))

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

func Decode(packet []byte) *DNSPacket {

	id := decodePart(packet, 0, 2)
	queryParams := decodePart(packet, 2, 4)
	qdCount := decodePart(packet, 4, 6)
	anCount := decodePart(packet, 6, 8)
	nsCount := decodePart(packet, 8, 10)
	arCount := decodePart(packet, 10, 12)

	isQuery := queryParams & DNSQuery
	opcode := (queryParams << 1) >> 12
	flags := (queryParams & FlagsMask) >> 8
	z := (queryParams & ZMask) >> 4
	rcode := (queryParams & RcodeMask)

	var queryType string
	if isQuery > 0 {
		queryType = "query"
	} else {
		queryType = "response"
	}

	fmt.Printf("Id: %d\n", id)
	fmt.Printf("query: %d\n", queryParams)
	fmt.Printf("isQuery: %d\n", isQuery)
	fmt.Printf("opcode: %d\n", decodeOpcode(int(opcode)))

	fmt.Printf("Flags: %d\n", flags)
	fmt.Printf("Z: %d\n", z)
	fmt.Printf("rcode: %d\n", decodeRcode(int(rcode)))
	fmt.Printf("qdCount: %d\n", qdCount)
	fmt.Printf("anCount: %d\n", anCount)
	fmt.Printf("nsCount: %d\n", nsCount)
	fmt.Printf("arCount: %d\n", arCount)

	dnsPacket := DNSPacket{
		Type:    queryType,
		ID:      id,
		Opcode:  decodeOpcode(int(opcode)),
		Flags:   int(flags),
		Qdcount: qdCount,
		Ancount: anCount,
		Nscount: nsCount,
		Arcount: arCount,
	}

	startOfQuestions := 12
	for i := 0; i < int(qdCount); i++ {
		qname, n := decodeQname(packet[startOfQuestions:])

		qTypeStart := startOfQuestions + n
		qTypeEnd := qTypeStart + 2
		qClassStart := qTypeEnd
		qClassEnd := qTypeEnd + 2
		fmt.Println(packet[qClassStart:qClassEnd])
		qtype := decodePart(packet, qTypeStart, qTypeEnd)
		qclass := decodePart(packet, qClassStart, qClassEnd)
		//qclass := binary.BigEndian.Uint16(packet[qClassStart:qClassEnd])
		startOfQuestions = n + 4

		fmt.Printf("Name: %s\n", qname)
		fmt.Printf("QType: %d\n", qtype)
		fmt.Printf("QClass: %d\n", qclass)
	}

	return &dnsPacket

}

func fromIntToBytes(num uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, num)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decodeOpcode(opcode int) int {
	var oc int

	if code := opcode & OpcodeStandardQuery; code > 0 {
		oc = code
	}

	if code := opcode & OpcodeInverseQuery; code > 0 {
		oc = code
	}

	if code := opcode & OpcodeServerStatus; code > 0 {
		oc = code
	}

	return oc
}

func decodeRcode(rcode int) int {
	var rc int

	if code := rcode & RcodeNoError; code > 0 {
		rc = code
	}

	if code := rcode & RcodeFormatError; code > 0 {
		rc = code
	}

	if code := rcode & RcodeServerFailure; code > 0 {
		rc = code
	}

	if code := rcode & RcodeNameError; code > 0 {
		rc = code
	}

	if code := rcode & RcodeNotImplemented; code > 0 {
		rc = code
	}

	if code := rcode & RcodeRefused; code > 0 {
		rc = code
	}

	return rc
}

func decodePart(packet []byte, start int, end int) uint16 {
	return binary.BigEndian.Uint16(packet[start:end])
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

func decodeQname(qname []byte) (string, int) {
	name := new(bytes.Buffer)
	start := 0

	for {
		labelSize := int(qname[start])
		label := qname[start+1 : labelSize+start+1]
		start = start + labelSize + 1

		name.WriteString(string(label))

		if qname[start] == 0 {
			start++
			break
		}

		name.WriteString(".")
	}

	return name.String(), start
}

func encodeQuestion(q Question) []byte {
	question := make([]byte, 0)

	name := encodeQname(q.Qname)
	qtype, _ := fromIntToBytes(uint16(1))  //hard coded
	qclass, _ := fromIntToBytes(uint16(1)) //hard coded

	question = append(question, name...)
	question = append(question, qtype...)
	question = append(question, qclass...)

	return question
}
