package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/mcombeau/go-dns-tools/decoder"
	"github.com/mcombeau/go-dns-tools/dns"
	"github.com/mcombeau/go-dns-tools/encoder"
	"github.com/mcombeau/go-dns-tools/printer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <domain> [question type]")
		return
	}

	domain := os.Args[1]
	questionType := dns.A
	if len(os.Args) == 3 {
		questionType = dns.GetCodeFromTypeString(os.Args[2])
	}
	dnsServer := "8.8.8.8:53" // Google's public DNS server

	message := &dns.Message{
		Header: &dns.Header{
			Id:            1234,
			Flags:         &dns.Flags{RecursionDesired: true},
			QuestionCount: 1,
		},
		Questions: []dns.Question{
			{
				Name:   domain,
				QType:  questionType,
				QClass: dns.IN,
			},
		},
	}

	data, err := encoder.EncodeDNSMessage(message)
	if err != nil {
		log.Fatalf("Failed to encode DNS message: %v\n", err)
	}

	startTime := time.Now()

	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		log.Fatalf("Failed to connect to DNS server: %v\n", err)
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Fatalf("Failed to send DNS query: %v\n", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	response := make([]byte, 512)
	n, err := conn.Read(response)
	if err != nil {
		log.Fatalf("Failed to read DNS response: %v\n", err)
	}

	queryTime := time.Since(startTime)

	decodedMessage, err := decoder.DecodeDNSMessage(response[:n])
	if err != nil {
		log.Fatalf("Failed to decode DNS response: %v\n", err)
	}

	printer.PrintDNSMessage(decodedMessage, domain)
	printer.PrintDNSQueryInfo(dnsServer, queryTime)
}