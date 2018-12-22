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

//Processes an answer based on the type
//and returns the correspnding type to the caller
//for example if we process an 'A' record we get back
//a RecrodTypeA which has a single property IPv4 address
//But if we process a 'SRV' record we get back a RecordTypeSRV
//which contains port,target etc.
//all these types implement the PacketProcessor interface
func (a *Answer) Process() PacketProcessor {
	var p PacketProcessor

	switch a.Type {
	case DNSRecordTypeA:
		p = &RecordTypeA{}

	case DNSRecordTypeSRV:
		p = &RecordTypeSRV{}

	default:
		p = &RecordTypeDefault{}
	}

	p.Process(*a)

	return p

}

func (a *Answer) Encode(offset int) []byte {
	answer := make([]byte, 0)
	//encode the name - should be compressed, but if packet contains 0 questions - we can't compress it
	var name []byte
	if offset != 0 {
		compressed := OffsetMarker | offset
		name, _ = fromIntToBytes(uint16(compressed))
	} else { //do not compress the name
		name = encodeQname(a.Name)
	}

	aType, _ := fromIntToBytes(uint16(a.Type))
	aClass, _ := fromIntToBytes(uint16(a.Class))
	ttl, _ := fromUint32ToBytes(a.TTL)
	rdLength, _ := fromIntToBytes(uint16(a.RdLength))

	answer = append(answer, name...)
	answer = append(answer, aType...)
	answer = append(answer, aClass...)
	answer = append(answer, ttl...)
	answer = append(answer, rdLength...)
	answer = append(answer, a.Data...)

	return answer
}

//PacketProcessor interface. All recordType's should
//implement this interface
type PacketProcessor interface {
	Process(Answer)
}
