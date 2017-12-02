package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	baddr   = flag.String("addr", "185.230.223.69", "IP addr you want to send out on")
	ipfname = flag.String("file", "iplist.txt", "The IP list to send queries to")
	dnsbase = flag.String("dbase", "s.flm.me.uk", "the zone that has the NS record pointing to it")
)

func main() {
	flag.Parse()

	rxListener, err := net.ListenPacket("udp4", *baddr+":53")
	if err != nil {
		log.Fatalf("failed to listen on UDP 53 %s", err.Error())
	}

	if !verifyNSsetup(*dnsbase, *baddr) {
		log.Fatalf("failed to confirm NS records are setup correctly, exiting")
	}

	startDNSListener(rxListener)

	http.HandleFunc("/upload", handleUpload)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalf("Failed to listen on HTTP: %s", http.ListenAndServe(*baddr+":5050", nil))
}

func verifyNSsetup(name, ipaddr string) bool {
	nsr, err := net.LookupNS(name)
	if err != nil {
		log.Printf("Unable to check NS setup correctly (1): %s", err.Error())
		return false
	}

	if len(nsr) != 1 {
		log.Printf("NS setup incorrect, I saw %d NS records for %s, I am expecting a single one.", len(nsr), name)
		return false
	}

	ar, err := net.LookupIP(nsr[0].Host)
	if len(ar) != 1 {
		log.Printf("NS setup incorrect, I saw %d A/AAAA records for %s, I am expecting a single one.", len(ar), nsr[0].Host)
		return false
	}

	if ar[0].String() != net.ParseIP(ipaddr).String() {
		log.Printf("NS setup incorrect, %s is suppose to point to %s, but it points to %s", nsr[0].Host, ipaddr, ar[0].String())
		return false
	}

	return true
}
