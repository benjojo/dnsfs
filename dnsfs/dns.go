package main

import (
	"net"

	"strings"

	"log"

	"github.com/miekg/dns"
)

func startDNSListener(socket net.PacketConn) {
	go DNSLoop(socket)
}

type storageRequest struct {
	storageNotifications chan bool
	content              string
	replications         int
}

var uploadPendingMap map[string]storageRequest

func DNSLoop(socket net.PacketConn) {
	for {
		dnsin := make([]byte, 1500)
		inbytes, inaddr, err := socket.ReadFrom(dnsin)

		inmsg := &dns.Msg{}

		if unpackErr := inmsg.Unpack(dnsin[inbytes:]); unpackErr != nil {
			log.Printf("Unable to unpack DNS request %s", err.Error())
			continue
		}

		if len(inmsg.Question) != 1 {

			log.Printf("More than one quesion in query (%d), droppin %+v", len(inmsg.Question), inmsg)
			continue
		}

		if !strings.Contains(inmsg.Question[0].Name, *dnsbase) {
			log.Printf("question is not for us '%s' vs expected '%s'", inmsg.Question[0].Name, *dnsbase)

			continue
		}

		outmsg := &dns.Msg{}

		ostring := make([]string, 1)
		ostring[0] = "aaaaaaaaaaaaaaaaaaaaa"

		inmsg.SetReply(outmsg)
		outmsg.Answer = make([]dns.RR, 1)
		outmsg.Answer[0] = &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   inmsg.Question[0].Name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    2147483646},
			Txt: ostring,
		}
		outputb, err := outmsg.Pack()

		if err != nil {
			log.Printf("unable to pack response to thing")
			continue
		}

		socket.WriteTo(outputb, inaddr)
	}
}

func verifyNSsetup(name string) bool {
	s, err := net.LookupTXT("tokentest" + name)
	if err != nil {
		return false
	}

	if len(s) != 1 {
		return false
	}

	return true
}
