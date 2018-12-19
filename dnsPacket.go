package dnsPacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	Z          int
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

func (dns *DNSPacket) AddQuestion(name string, qclass int, qtype int) *Question {
	question := Question{
		Qname:  name,
		Qclass: qclass,
		Qtype:  qtype,
	}

	dns.Questions = append(dns.Questions, question)

	return &question
}

func (dns *DNSPacket) AddAnswer(name string, aclass int, atype int, ttl uint32, dataLength int, data []byte) *Answer {
	answer := Answer{
		Name:     name,
		Class:    aclass,
		Type:     atype,
		TTL:      ttl,
		RdLength: dataLength,
	}
	answer.Data = make([]byte, len(data))
	for i := range data {
		answer.Data[i] = data[i]
	}

	dns.Answers = append(dns.Answers, answer)

	return &answer
}

func (dns *DNSPacket) IsAuthoritativeAnswer() bool {
	if (dns.Flags & AuthoritativeAnswerMask) > 0 {
		return true
	}

	return false
}

func (dns *DNSPacket) IsTruncated() bool {
	if (dns.Flags & TruncationMask) > 0 {
		return true
	}

	return false
}

func (dns *DNSPacket) IsRecursionDesired() bool {
	if (dns.Flags & RecursionDesiredMask) > 0 {
		return true
	}

	return false
}

func (dns *DNSPacket) IsRecursionAvailable() bool {
	if (dns.Flags & RecursionAvailableMask) > 0 {
		return true
	}

	return false
}

func (dns DNSPacket) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("Type: %s\n", dns.Type))
	buf.WriteString(fmt.Sprintf("ID: %d\n", dns.ID))
	buf.WriteString(fmt.Sprintf("Opcode: %d\n", dns.Opcode))
	buf.WriteString(fmt.Sprintf("Z: %d\n", dns.Z))
	buf.WriteString(fmt.Sprintf("Rcode: %d\n", dns.Rcode))
	buf.WriteString(fmt.Sprintf("Flags:\n"))
	buf.WriteString(fmt.Sprintf(" --AA: %t\n", dns.IsAuthoritativeAnswer()))
	buf.WriteString(fmt.Sprintf(" --TC: %t\n", dns.IsTruncated()))
	buf.WriteString(fmt.Sprintf(" --RD: %t\n", dns.IsRecursionDesired()))
	buf.WriteString(fmt.Sprintf(" --RA: %t\n", dns.IsRecursionAvailable()))
	buf.WriteString(fmt.Sprintf("Question Count: %d\n", dns.Qdcount))
	buf.WriteString(fmt.Sprintf("Answer Count: %d\n", dns.Ancount))
	buf.WriteString(fmt.Sprintf("NS Count: %d\n", dns.Nscount))
	buf.WriteString(fmt.Sprintf("AR Count: %d\n", dns.Arcount))
	buf.WriteString(fmt.Sprintf("Questions:\n"))

	for i, q := range dns.Questions {
		buf.WriteString(fmt.Sprintf("%d - %s\n", i, q))
	}

	buf.WriteString(fmt.Sprintf("Answers:\n"))

	for i, a := range dns.Answers {
		buf.WriteString(fmt.Sprintf("%d - %s\n", i, a))
	}

	buf.WriteString("\n")

	return buf.String()
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
	params := packetType | dnsPacket.Opcode | dnsPacket.Flags | dnsPacket.Z | dnsPacket.Rcode
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

	//header values
	id := decodePart(packet, 0, 2)
	queryParams := decodePart(packet, 2, 4)
	qdCount := decodePart(packet, 4, 6)
	anCount := decodePart(packet, 6, 8)
	nsCount := decodePart(packet, 8, 10)
	arCount := decodePart(packet, 10, 12)

	isQuery := queryParams & DNSQuery
	opcode := (queryParams << 1) >> 12
	flags := queryParams & FlagsMask
	z := (queryParams & ZMask) >> 4
	rcode := (queryParams & RcodeMask)

	var queryType string
	if isQuery > 0 {
		queryType = "query"
	} else {
		queryType = "response"
	}

	dnsPacket := DNSPacket{
		Type:    queryType,
		ID:      id,
		Opcode:  decodeOpcode(int(opcode)),
		Flags:   int(flags),
		Qdcount: qdCount,
		Ancount: anCount,
		Nscount: nsCount,
		Arcount: arCount,
		Rcode:   int(rcode),
		Z:       int(z),
	}

	//process questions
	startOfQuestions := 12
	for i := 0; i < int(qdCount); i++ {
		qname, n := decodeQname(packet[startOfQuestions:])

		qTypeStart := startOfQuestions + n
		qTypeEnd := qTypeStart + 2
		qClassStart := qTypeEnd
		qClassEnd := qTypeEnd + 2

		qtype := decodePart(packet, qTypeStart, qTypeEnd)
		qclass := decodePart(packet, qClassStart, qClassEnd)
		startOfQuestions = startOfQuestions + n + 4

		dnsPacket.AddQuestion(qname, int(qclass), int(qtype))

	}

	//process answers
	startOfAnswers := startOfQuestions
	for i := 0; i < int(anCount); i++ {
		compressedAnswerName := decodePart(packet, startOfAnswers, startOfAnswers+2)
		offset := compressedAnswerName & CompressedAnswerMask
		answerName, _ := decodeQname(packet[offset:])

		startOfAnswerType := startOfAnswers + 2
		endOfAnswerType := startOfAnswerType + 2
		startOfAnswerClass := endOfAnswerType
		endOfAnswerClass := startOfAnswerClass + 2
		startOfTTL := endOfAnswerClass
		endOfTTL := startOfTTL + 4
		startOfDataLength := endOfTTL
		endOfDataLength := startOfDataLength + 2
		startOfData := endOfDataLength
		anType := decodePart(packet, startOfAnswerType, endOfAnswerType)
		anClass := decodePart(packet, startOfAnswerClass, endOfAnswerClass)
		ttl := binary.BigEndian.Uint32(packet[startOfTTL:endOfTTL])
		dataLength := decodePart(packet, startOfDataLength, endOfDataLength)
		endOfData := startOfData + int(dataLength)

		dnsPacket.AddAnswer(answerName, int(anClass), int(anType), ttl, int(dataLength), packet[startOfData:endOfData])

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
