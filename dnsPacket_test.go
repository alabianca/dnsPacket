package dnsPacket

import (
	"reflect"
	"testing"
)

//Test Encoding a query with a single question
func TestEncodeSingleQuestion(t *testing.T) {
	packet := DNSPacket{
		Type:    "query",
		ID:      1,
		Opcode:  0,
		Z:       0,
		Rcode:   0,
		Flags:   FlagsRecurionDesired,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}

	packet.AddQuestion("google.com", 1, 1)

	encoded := []byte{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	data := Encode(&packet)

	for i := range data {
		if data[i] != encoded[i] {
			t.Errorf("Failed. \nGot:  %b \nWant: %b \n", data, encoded)
		}
	}
}

func TestEncodeMultipleQuestions(t *testing.T) {
	packet := DNSPacket{
		Type:    "query",
		ID:      1,
		Opcode:  0,
		Z:       0,
		Rcode:   0,
		Flags:   FlagsRecurionDesired,
		Qdcount: 2,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}

	packet.AddQuestion("google.com", 1, 1)
	packet.AddQuestion("google.com", 1, 1)

	encoded := []byte{0, 1, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	data := Encode(&packet)

	for i := range data {
		if data[i] != encoded[i] {
			t.Errorf("Failed. \nGot:  %b \nWant: %b \n", data, encoded)
		}
	}

}

func TestDecodeQueryWithSingleQuestion(t *testing.T) {
	data := []byte{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	packet := Decode(data)

	compare := DNSPacket{
		Type:    "query",
		ID:      1,
		Opcode:  0,
		Z:       0,
		Rcode:   0,
		Flags:   FlagsRecurionDesired,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}

	compare.AddQuestion("google.com", 1, 1)

	if !reflect.DeepEqual(packet, &compare) {
		t.Errorf("Fail.\n Got: \n%s\n Want: \n%s\n", packet, compare)
	}
}

func TestDecodeQueryWithMultipleQuestions(t *testing.T) {
	data := []byte{0, 1, 1, 0, 0, 2, 0, 0, 0, 0, 0, 0, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	packet := Decode(data)

	compare := DNSPacket{
		Type:    "query",
		ID:      1,
		Opcode:  0,
		Z:       0,
		Rcode:   0,
		Flags:   FlagsRecurionDesired,
		Qdcount: 2,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}

	compare.AddQuestion("google.com", 1, 1)
	compare.AddQuestion("google.com", 1, 1)

	if !reflect.DeepEqual(packet, &compare) {
		t.Errorf("Fail.\n Got: \n%s\n Want: \n%s\n", packet, compare)
	}
}

func TestPacketFlags(t *testing.T) {
	packet := DNSPacket{
		Type:    "query",
		ID:      1,
		Opcode:  0,
		Z:       0,
		Rcode:   0,
		Flags:   FlagsRecurionDesired | FlagsAuthoritativeAnswer | FlagsTruncation | FlagsRecursionAvailable,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}

	packet.AddQuestion("google.com", 1, 1)

	isRd := packet.IsRecursionDesired()
	isAA := packet.IsAuthoritativeAnswer()
	isTC := packet.IsTruncated()
	isRA := packet.IsRecursionAvailable()

	if !(isRd && isAA && isTC && isRA) {
		t.Errorf("Failed.\nGot: %t %t %t %t \nWant: %t %t %t %t", isRd, isAA, isTC, isRA, true, true, true, true)
	}
}
