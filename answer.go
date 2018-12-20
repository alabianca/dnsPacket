package dnsPacket

import (
	"bytes"
	"fmt"
)

/*
0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                                               |
/                                               /
/ 1 1 |    NAME COMPRESSED FORMAT               /
|                                               |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      TYPE                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     CLASS                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      TTL                      |
|                                               |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   RDLENGTH                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
/                     RDATA                     /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type Answer struct {
	Name     string
	Type     int
	Class    int
	TTL      uint32
	RdLength int
	Data     []byte
}

func (a Answer) String() string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("Name: %s Type: %d Class: %d TTL: %d DataLength: %d\n", a.Name, a.Type, a.Class, a.TTL, a.RdLength))
	buf.WriteString(fmt.Sprint("Data: ", a.Data))

	return buf.String()
}

func (a *Answer) Process() PacketProcessor {
	var p PacketProcessor

	switch a.Type {
	case DNSRecordTypeA:
		p = &RecordTypeA{}
	default:
		p = &RecordTypeA{}
	}

	p.Process(*a)

	return p

}

type PacketProcessor interface {
	Process(Answer)
}

type RecordTypeA struct {
	IPv4 string
}

func (record *RecordTypeA) Process(a Answer) {
	record.IPv4 = "Hello World"
}
