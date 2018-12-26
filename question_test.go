package dnsPacket

import "testing"

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
