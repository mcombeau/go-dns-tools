package decoder

import (
	"testing"

	"github.com/mcombeau/go-dns-tools/testutils"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDecodeDomainName(t *testing.T) {
	tests := []struct {
		data     []byte
		offset   int
		expected string
		length   int
	}{
		{
			data:     []byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
			offset:   0,
			expected: "example.com.",
			length:   13,
		},
		{
			data:     []byte{3, 'w', 'w', 'w', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0},
			offset:   0,
			expected: "www.example.com.",
			length:   17,
		},
		{
			data:     []byte{0},
			offset:   0,
			expected: ".",
			length:   1,
		},
	}

	for _, test := range tests {
		name, length, err := decodeDomainName(test.data, test.offset)

		assert.NoError(t, err)

		assert.Equal(t, test.expected, name)
		assert.Equal(t, test.length, length)
	}
}

func TestDecodeDomainNameInDNSMessage(t *testing.T) {
	mockResponse, err := testutils.MockDNSResponse()
	if err != nil {
		t.Fatalf("Failed to create mock response: %v\n", err)
	}

	var unpackedMockResponse dns.Msg
	err = unpackedMockResponse.Unpack(mockResponse)

	if err != nil {
		t.Fatalf("Failed to unpack mock response: %v\n", err)
	}

	want := unpackedMockResponse.Question[0].Name
	wantBytes := len(unpackedMockResponse.Question[0].Name) + 1
	got, gotBytes, err := decodeDomainName(mockResponse, 12)

	assert.NoError(t, err)

	assert.Equal(t, want, got)
	assert.Equal(t, wantBytes, gotBytes)
}

func TestDecodeDomainNameInCompressedDNSMessage(t *testing.T) {
	mockResponse, err := testutils.MockDNSCompressedResponse()
	if err != nil {
		t.Fatalf("Failed to create mock response: %v\n", err)
	}

	var unpackedMockResponse dns.Msg
	err = unpackedMockResponse.Unpack(mockResponse)
	if err != nil {
		t.Fatalf("Failed to unpack mock response: %v\n", err)
	}

	headerLength := 12
	domainLength := len(unpackedMockResponse.Question[0].Name) + 1
	questionLength := domainLength + 4 // +4 for QTYPE and QCLASS
	answerOffset := headerLength + questionLength

	want := unpackedMockResponse.Answer[0].Header().Name
	got, gotBytes, err := decodeDomainName(mockResponse, answerOffset)
	wantBytes := 2 // Pointer size is 2 bytes

	assert.NoError(t, err)

	assert.Equal(t, want, got)
	assert.Equal(t, wantBytes, gotBytes)
}