package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	// Autoload .env files
	"github.com/jbuchbinder/cadmonitor/monitor"
	_ "github.com/joho/godotenv/autoload"
)

var (
	baseURL       = flag.String("baseUrl", "http://cadview.qvec.org/", "Base URL")
	monitorType   = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	pollInterval  = flag.Int("poll-interval", 5, "Poll interval in seconds")
	inactivityMin = flag.Int("inactivity", 10, "Inactivity in minutes before culling")
	suffix        = flag.String("suffix", "", "Unit suffix to restrict polling to (i.e. 63 for STA63 units)")
	fdid          = flag.String("fdid", "04042", "FDID for agency")

	mutex       sync.Mutex
	activeCalls map[string]monitor.CallStatus
)

func main() {
	flag.Parse()

	cadbrowser, err := monitor.InstantiateCadMonitor(*monitorType)
	if err != nil {
		panic(err)
	}
	err = cadbrowser.ConfigureFromValues(map[string]string{
		"baseUrl": *baseURL,
		"suffix":  *suffix,
		"fdid":    *fdid,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Logging into CAD interface")
	err = cadbrowser.Login(os.Getenv("CADUSER"), os.Getenv("CADPASS"))
	if err != nil {
		panic(err)
	}
	log.Printf("Starting main loop")

	activeCalls = map[string]monitor.CallStatus{}

	// Cull active calls every poll interval x 4
	go func() {
		log.Printf("Culling loop started")
		for {

			log.Printf("Culling job running")

			mutex.Lock()

			for k := range activeCalls {
				if time.Since(activeCalls[k].LastUpdated) > (time.Duration(*inactivityMin) * time.Minute) {
					log.Printf("Removing call %s due to inactivity", k)
					delete(activeCalls, k)
				}
			}

			mutex.Unlock()

			for iter := 0; iter < *pollInterval*2; iter++ {
				time.Sleep(time.Second)
				if cadbrowser.TerminateMonitor() {
					log.Printf("Terminating culling thread")
					break
				}
			}

		}
	}()

	cadbrowser.Monitor(func(call monitor.CallStatus) error {
		mutex.Lock()
		defer mutex.Unlock()

		log.Printf("Callback triggered with ID : %s", call.ID)

		// Check to see if it's in the active map
		if _, ok := activeCalls[call.ID]; !ok {
			// Store new copy, act accordingly
			activeCalls[call.ID] = call

			// Notify that there's a new call
			return notifyNewCall(call)
		}

		if call.District != activeCalls[call.ID].District {
			// Update if district has been updated
			err := notifyCallDifferences(activeCalls[call.ID], call)
			activeCalls[call.ID] = call
			return err
		}

		// See if there are any differences
		if len(call.Units) != len(activeCalls[call.ID].Units) {
			err := notifyUnitDifferences(activeCalls[call.ID], call)
			activeCalls[call.ID] = call
			return err
		}

		// Update last updated to keep it frosty
		activeCalls[call.ID] = call

		log.Printf("Status: %#v", call)
		return nil
	}, *pollInterval)

}

func notifyNewCall(call monitor.CallStatus) error {
	log.Printf("notifyNewCall: %#v", call)
	return nil
}

func notifyCallDifferences(orig monitor.CallStatus, updated monitor.CallStatus) error {
	log.Printf("notifyCallDifferences: %#v -> %#v", orig, updated)
	return nil
}

func notifyUnitDifferences(orig monitor.CallStatus, updated monitor.CallStatus) error {
	log.Printf("notifyUnitDifferences: %#v -> %#v", orig, updated)
	return nil
}
