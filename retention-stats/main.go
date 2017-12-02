package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ProbeInfo struct {
	ReportedTimes map[int64]TimeReport
	UsedResolvers map[string]int
}

type TimeReport struct {
	FirstSeen int
	LastSeen  int
}

func main() {
	detailedmode := flag.Bool("detailed", false, "Print individual values")
	csvmode := flag.Bool("csv", false, "Print results as CSV")
	flag.Parse()

	probes := make(map[string]ProbeInfo)

	reader := bufio.NewReader(os.Stdin)

	for {
		lnb, _, err := reader.ReadLine()

		if err != nil {
			break
		}

		parts := strings.Split(string(lnb), ",")
		if len(parts) != 3 {
			continue
		}

		// datapoint := RipeData{}
		// err = json.Unmarshal(lnb, &datapoint)
		if err != nil {
			log.Printf("Unable to decode data point %s", err.Error())
			continue
		}

		dnsval := strings.Split(parts[2], "\"")
		if len(dnsval) != 3 {
			continue
		}
		i, _ := strconv.ParseInt(dnsval[1], 10, 64)
		timestamp64, _ := strconv.ParseInt(parts[0], 10, 64)
		timestamp := int(timestamp64)
		probeid := parts[1]

		if probes[probeid].ReportedTimes == nil {
			tmp := probes[probeid]
			tmp.ReportedTimes = make(map[int64]TimeReport)
			tmp.UsedResolvers = make(map[string]int)
			probes[probeid] = tmp
		}

		rt := probes[probeid].ReportedTimes
		if rt[i].FirstSeen == 0 {
			rttmp := TimeReport{}
			rttmp.FirstSeen = timestamp
			rt[i] = rttmp
		} else if rt[i].LastSeen < timestamp {
			rttmp := rt[i]
			rttmp.LastSeen = timestamp
			rt[i] = rttmp
		} else if rt[i].FirstSeen > timestamp {
			rttmp := rt[i]
			rttmp.FirstSeen = timestamp
			rt[i] = rttmp
		}

	}

	if *csvmode {
		fmt.Print("timeheld,probeid,resolvers\n")
	}

	for probeid, probedata := range probes {
		if !(*csvmode) {
			fmt.Printf("For Probe ID %s:\n", probeid)
			fmt.Printf("Resolvers %v\n", probedata.UsedResolvers)
		}
		recordTime := 0.0

		for timestamp, tsdata := range probedata.ReportedTimes {
			if tsdata.LastSeen == 0 {
				if *detailedmode {
					fmt.Printf("\tFor value %d, never got cached\n", timestamp)
				}
			} else {
				fs := time.Unix(int64(tsdata.FirstSeen), 0)
				ls := time.Unix(int64(tsdata.LastSeen), 0)
				between := ls.Sub(fs)
				if *detailedmode {
					fmt.Printf("\tFor value %d, was in cache for %f min\n", timestamp, between.Minutes())
				}
				if between.Minutes() > recordTime {
					recordTime = between.Minutes()
				}
			}
		}

		if !(*csvmode) {
			fmt.Printf("For Probe ID %s held on at best for %f min\n\n", probeid, recordTime)
		} else {
			fmt.Printf("%0.f,%s\r\n", recordTime, probeid)
		}

	}
}
