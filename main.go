package main

import (
	"flag"
	"log"
	"time"

	"github.com/jbuchbinder/qveccadmonitor/monitor"
)

var (
	baseURL      = flag.String("baseUrl", "http://cadview.qvec.org/NewWorld.CAD.ViewOnly/", "Base URL")
	monitorType  = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	pollInterval = flag.Int("poll-interval", 15, "Poll interval in seconds")
	suffix       = flag.String("suffix", "", "Unit suffix to restrict polling to (i.e. 63 for STA63 units)")
)

func main() {
	flag.Parse()

	cadbrowser, err := monitor.GetCadMonitor(*monitorType)
	if err != nil {
		panic(err)
	}
	err = cadbrowser.ConfigureFromValues(map[string]string{
		"baseUrl": *baseURL,
		"suffix":  *suffix,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Logging into CAD interface")
	err = cadbrowser.Login(monitor.USER, monitor.PASS)
	if err != nil {
		panic(err)
	}
	for {
		log.Printf("Starting main loop")

		calls, err := cadbrowser.GetActiveCalls()
		if err != nil {
			log.Printf("err: %s", err.Error())
			goto sleeploop
		}

		if len(calls) == 0 {
			log.Printf("No active calls")
			goto sleeploop
		}

		for _, callurl := range calls {
			status, err := cadbrowser.GetStatus(callurl)
			if err != nil {
				log.Printf("err: %s", err.Error())
				continue
			}

			// TODO: process data for call instead of displaying
			log.Printf("Status: %#v", status)
		}

		// Sleep during poll interval
	sleeploop:
		log.Printf("Sleeping for %d seconds", *pollInterval)
		for interval := 0; interval < *pollInterval; interval++ {
			time.Sleep(1 * time.Second)
		}
	}
}
