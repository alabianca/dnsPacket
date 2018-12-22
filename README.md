# Decode and Encode DNS Packets

## Type - DNSPacket

```go
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
```

#### DNSPacket.Type
Determines if a packet is a query or a response. Set this to "query" or "response"

#### DNSPacket.ID
The transaction Id of the packet (max uint16)

#### DNSPacket.Opcode
Determines the query type. Possible values are:
       * `0`: Standard Query
       * `1`: Inverse query
       * `2`: Server status request
       * `3-15`: Reserved for future use

#### DNSPacket.Z
Reserved for future use. Set this field to 0

#### DNSPacket.Rcode
The response code. Possible values are:
    * `0`: No Error condition
    * `1`: Format Error
    * `2`: Server failure
    * `3`: Name Error
    * `4`: Not Implemented
    * `5`: Refused

#### DNSPacket.Flags
Sets the bit for the following flags: `AA` (Authoritative Answer), `TC` (Truncation), `RD` (Recursion Desired) and `RA` (Recursion Available)

Example of setting the `RD` flag.
```go
    packet := DNSPacket {
        ...
        Flags: dnsPacket.FlagsRecursionDesired
        ...
    }
```
If you need to set multiple flags simply bitwise OR them together
```go
    packet := DNSPacket {
        ...
        Flags: dnsPacket.FlagsRecursionDesired | dnsPacket.FlagsAuthoritativeAnswer
        ...
    }
```

#### DNSPacket.Qdcount
How many questions are on this packet

#### DNSPacket.Ancount
How many answers are on this packet

#### DNSPacket.Nscount
How many name server resource records are in this packet

#### DNSPacket.Arcount
How many additional records are in this packet

#### DNSPacket.Questions
These are the questions on the packet. A question has the following format.

```go
type Question struct {
	Qname  string
	Qtype  int
	Qclass int
}
```
#### Question.Qname
The domain name we are querying example: `google.com` or `_someService._tcp.local`

#### Question.Qtype
The DNS Record type (`A`, `SRV` etc.) Set these to their numerical equivalents. 
For example if you are looking for `A` records set this field to `1`, but if you are looking for `SRV` records set this field to `33`. 

#### Question.Qclass
The class we are looking for. For example `IN` . Also set this to their numerical value. If you are looking for `IN` set it to `1`. (I should set this to the default...)

#### DNSPacket.Answers
Resource records answering the question. Answers have the follwing format

```go
type Answer struct {
	Name     string
	Type     int
	Class    int
	TTL      uint32
	RdLength int
	Data     []byte
}
```

#### Answer.Name
The domain name answering

#### Answer.Type
Same as in `Question`

#### Answer.Class
Same as in `Question`

#### Answer.TTL
Time to live for this response (in seconds)

#### Answer.RdLength
Length of the `Answer.Data` field. 

#### Answer.Data
The data of the response in pure bytes. For example if you are querying for an `A` record, this field will contain the Ipv4 address. To get a concrete type out of the `Data` field, call `.Process()` method on answer (see example below) .

#### DNSPacket.Additional
Not yet implemented...

#### DNSPacket.Authority
Not yet implemented...

## Methods - DNSPacket

#### AddQuestion(name string, qclass int, qtype int) *Question
Adds a quesion to the packet. 

Example: add a question for `google.com` class `IN` and type `A`
```go
packet.AddQuestion("google.com", 1, 1)
```

#### AddAnswer(name string, aclass int, atype int, ttl uint32, dataLength int, data []byte) *Answer
Adds an answer to the packet

## Functions

#### Encode(dnsPacket *DNSPacket) []byte
Encodes a packet and get back the raw bytes

#### Decode(packet []byte) *DNSPacket
Decodes bytes and returns a pointer to a DNSPacket



