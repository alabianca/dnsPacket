package dnsPacket

import (
	"bytes"
	"encoding/binary"
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

0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| ID                                            |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR| Opcode |AA|TC|RD|RA| Z | RCODE             |
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

func fromIntToBytes(num uint16) ([]byte, error) {
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
	qtype, _ := fromIntToBytes(uint16(1))  //hard coded
	qclass, _ := fromIntToBytes(uint16(1)) //hard coded

	question = append(question, name...)
	question = append(question, qtype...)
	question = append(question, qclass...)

	return question
}
