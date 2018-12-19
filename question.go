package dnsPacket

import (
	"bytes"
	"strings"
)

type Question struct {
	Qname  string
	Qtype  int
	Qclass int
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
