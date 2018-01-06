package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	baddr        = flag.String("addr", "185.230.223.69", "IP addr you want to send out on")
	ipfname      = flag.String("file", "iplist.txt", "The IP list to send queries to")
	dnsbase      = flag.String("dbase", "s.flm.me.uk", "the zone that has the NS record pointing to it")
	ipList       = make([]string, 0)
	globalSender net.PacketConn
)

func main() {
	flag.Parse()

	txListener, err := net.ListenPacket("udp4", "0.0.0.0:34123")
	if err != nil {
		log.Fatalf("failed to listen on UDP 34123 (for tx) %s", err.Error())
	}
	globalSender = txListener

	ipList = parseIPList(*ipfname)

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

func parseIPList(path string) []string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read IP list, %s", err.Error())
	}

	lines := strings.Split(string(bytes), "\n")
	for ln, s := range lines {
		if s == "" {
			continue
		}
		t := net.ParseIP(s)
		if t == nil {
			log.Fatalf("Error in IP list on line %d", ln)
		}
	}

	return lines
}
