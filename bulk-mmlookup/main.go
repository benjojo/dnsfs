package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	typ := flag.String("type", "country", "use country or isp")
	flag.Parse()

	db, err := geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		fmt.Printf("Unable to open GeoLite2-ASN.mmdb , %s\n", err.Error())
		os.Exit(1)
	}
	db2, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		fmt.Printf("Unable to open GeoLite2-Country.mmdb , %s\n", err.Error())
		os.Exit(1)
	}

	linereader := bufio.NewReader(os.Stdin)

	for {
		l, _, err := linereader.ReadLine()
		if err != nil {
			break
		}

		ipaddr := net.ParseIP(string(l))
		if ipaddr == nil {
			continue
		}

		if *typ == "country" {
			c, err := db2.Country(ipaddr)
			if err != nil {
				continue
			}
			fmt.Printf("%s\t%s\n", ipaddr, c.Country.Names["en"])
		} else if *typ == "isp" {
			c, err := db.ASN(ipaddr)
			if err != nil {
				continue
			}
			fmt.Printf("%s\t%d\t%s\n", ipaddr, c.AutonomousSystemNumber, c.AutonomousSystemOrganization)
		} else {
			fmt.Println("Invalid -type, use country or isp")
			os.Exit(1)
		}

	}

}
