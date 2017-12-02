package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/miekg/dns"
)

func main() {
	addr := flag.String("addr", "185.230.223.69", "IP addr you want to send out on")
	fname := flag.String("file", "iplist.txt", "The IP list to send queries to")
	flag.Parse()

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{
		Name:   "version.bind.",
		Qtype:  dns.TypeTXT,
		Qclass: dns.ClassCHAOS,
	}
	dnspacket, _ := m1.Pack()

	addrt := *addr

	listener, err := net.ListenPacket("udp4", addrt+":5353")
	if err != nil {
		log.Fatalf("failed to listen on UDP %s", err.Error())
	}

	fd, err := os.Open(*fname)
	if err != nil {
		log.Fatalf("open ip list file %s", err.Error())
	}
	bior := bufio.NewReader(fd)
	go ReadRes(listener)
	for {
		ips, _, err := bior.ReadLine()
		if err != nil {
			break
		}

		// ip := net.ParseIP(string(ips))

		// if ip == nil {
		// 	continue
		// }

		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:53", string(ips)))
		if err != nil {
			continue
		}

		listener.WriteTo(dnspacket, addr)
		time.Sleep(time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	// listener.WriteTo()
}

func ReadRes(conn net.PacketConn) {
	for {
		buf := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("err reading %s", err)
			continue
		}
		msg := &dns.Msg{}
		err = msg.Unpack(buf[:n])

		if err != nil {
			log.Printf("err parsing %s", err)
			continue
		}

		if len(msg.Answer) != 1 {
			continue
		}

		fmt.Printf("%d,%s,%s\n", time.Now().Unix(), addr.String(), msg.Answer[0].String())
	}

}
