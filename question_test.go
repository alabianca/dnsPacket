package dnsPacket

import (
	"testing"
)

func TestDecodeQname(t *testing.T) {

	tables := []struct {
		in  []byte
		out string
	}{
		{
			[]byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0},
			"google.com",
		},
		{
			[]byte{7, 95, 103, 111, 100, 114, 111, 112, 4, 95, 116, 99, 112, 5, 108, 111, 99, 97, 108, 0},
			"_godrop._tcp.local",
		},
	}

	for _, table := range tables {
		decoded, _ := decodeQname(table.in)

		if decoded != table.out {
			t.Errorf("Fail\nGot: %s\nWant: %s", decoded, table.out)
		}
	}
}

func TestEncodeQname(t *testing.T) {
	tables := []struct {
		out []byte
		in  string
	}{
		{
			[]byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0},
			"google.com",
		},
		{
			[]byte{7, 95, 103, 111, 100, 114, 111, 112, 4, 95, 116, 99, 112, 5, 108, 111, 99, 97, 108, 0},
			"_godrop._tcp.local",
		},
	}

	for _, table := range tables {
		encoded := encodeQname(table.in)

		for i := range encoded {
			if encoded[i] != table.out[i] {
				t.Errorf("Fail\nGot: %b\nWant: %b\n", encoded, table.out)
			}
		}
	}
}

func TestDecodeQuestion(t *testing.T) {
	questionBytes := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	qname, qtype, qclass, _ := decodeQuestion(questionBytes)

	if qname != "google.com" {
		t.Errorf("Fail\nGot: %s\nWant: %s\n", qname, "google.com")
	}

	if qtype != 1 {
		t.Errorf("Fail\nGot: %d\n Want: %d\n", qtype, 1)
	}

	if qclass != 1 {
		t.Errorf("Fail\nGot: %d\n Want: %d\n", qclass, 1)
	}

}

func TestEncodeQuestion(t *testing.T) {
	question := Question{
		Qname:  "google.com",
		Qclass: 1,
		Qtype:  1,
	}

	questionBytes := encodeQuestion(question)
	expectedBytes := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}

	for i := range questionBytes {
		if questionBytes[i] != expectedBytes[i] {
			t.Errorf("Fail\nGot: %b\nWant: %b\n", questionBytes, expectedBytes)
		}
	}
}
