package decoder

import (
	"errors"
	"fmt"
)

type DNSMessage struct {
	Header      *DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSResourceRecord
	NameServers []DNSResourceRecord
	Additionals []DNSResourceRecord
}

func DecodeDNSMessage(data []byte) (*DNSMessage, error) {
	if len(data) < 12 {
		return nil, errors.New("invalid DNS message: too short")
	}

	header, err := DecodeDNSHeader(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS header: %v", err)
	}
	offset := 12

	questions := make([]DNSQuestion, 0, header.QuestionCount)
	for i := 0; i < int(header.QuestionCount); i++ {
		question, newOffset, err := decodeDNSQuestion(data, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DNS question: %v", err)
		}
		questions = append(questions, *question)
		offset = newOffset
	}

	answers, offset, err := decodeResourceRecords(data, offset, header.AnswerRRCount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS answer: %v", err)
	}

	nameServers, offset, err := decodeResourceRecords(data, offset, header.NameserverRRCount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS authority name server: %v", err)
	}

	additionals, _, err := decodeResourceRecords(data, offset, header.AdditionalRRCount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS answer: %v", err)
	}

	return &DNSMessage{
		Header:      header,
		Questions:   questions,
		Answers:     answers,
		NameServers: nameServers,
		Additionals: additionals,
	}, nil
}

func decodeResourceRecords(data []byte, offset int, count uint16) ([]DNSResourceRecord, int, error) {
	records := make([]DNSResourceRecord, 0, count)
	for i := 0; i < int(count); i++ {
		record, newOffset, err := decodeDNSResourceRecord(data, offset)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, *record)
		offset = newOffset
	}
	return records, offset, nil
}