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

	startDNSListener(rxListener)

	// if !verifyNSsetup(*dnsbase) {
	// 	log.Fatalf("failed to confirm NS records are setup correctly, exiting")
	// }

	http.HandleFunc("/upload", handleUpload)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatalf("Failed to listen on HTTP: %s", http.ListenAndServe(*baddr+":5050", nil))
}
