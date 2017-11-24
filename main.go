package main

import (
	"flag"
	"log"
	"time"
)

var (
	Suffix       = flag.String("suffix", "", "Unit suffix to restrict polling to (i.e. 63 for STA63 units)")
	PollInterval = flag.Int("poll-interval", 15, "Poll interval in seconds")
)

func main() {
	flag.Parse()

	cadbrowser := CadBrowser{}
	log.Printf("Logging into CAD interface")
	err := cadbrowser.Login(USER, PASS)
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
		log.Printf("Sleeping for %d seconds", *PollInterval)
		for interval := 0; interval < *PollInterval; interval++ {
			time.Sleep(1 * time.Second)
		}
	}
}
