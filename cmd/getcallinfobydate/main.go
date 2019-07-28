package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jbuchbinder/cadmonitor/monitor"
)

var (
	activeCalls = flag.Bool("active", false, "Show active calls")
	baseURL     = flag.String("baseUrl", "http://cadview.qvec.org/", "Base URL")
	dateFlag    = flag.String("date", "06/02/2018", "Date in MM/DD/YYYY format")
	fdid        = flag.String("fdid", "04042", "FDID")
	monitorType = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	suffix      = flag.String("suffix", "", "Limit units to ones having a specfic suffix")
)

func main() {
	flag.Parse()

	m, err := monitor.InstantiateCadMonitor(*monitorType)
	if err != nil {
		panic(err)
	}
	err = m.ConfigureFromValues(map[string]string{
		"baseUrl": *baseURL,
		"suffix":  *suffix,
		"fdid":    *fdid,
	})
	if err != nil {
		panic(err)
	}
	err = m.Login(monitor.USER, monitor.PASS)
	if err != nil {
		panic(err)
	}

	var calls map[string]string
	var c []string
	if *activeCalls {
		c, err = m.GetActiveCalls()
		for k, v := range c {
			calls[fmt.Sprintf("%d", k)] = v
		}
	} else {
		calls, err = m.GetClearedCalls(*dateFlag)
	}
	if err != nil {
		panic(err)
	}
	for k, v := range calls {
		cs, err := m.GetStatus(v)
		if err != nil {
			panic(err)
		}
		fmt.Println("-------------------------------------------------------------------------------------------")
		fmt.Printf("CALL %s (PRI %d @ %s)\n", k, cs.Priority, cs.Location)

		// Find dispatch time from STA63 or RES63 or 63OFF
		var dispatchTime string
		{
			for _, u := range cs.Units {
				if u.Unit == "STA63" || u.Unit == "RES63" || u.Unit == "63OFF" {
					dispatchTime = u.DispatchTime
					break
				}
			}
			fmt.Printf("  - Dispatch Time        : %s\n", dispatchTime)
		}

		// Arrived/Cleared time is first arrived time from *63 units
		var arrivalTime, clearedTime, enrouteTime string
		{
			for _, u := range cs.Units {
				if strings.HasSuffix(u.Unit, "63") && u.Unit != "STA63" && u.Unit != "RES63" {
					if enrouteTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						enrouteTime = u.EnRouteTime
					}

					if monitor.CadDateTime(u.EnRouteTime).Unix() < monitor.CadDateTime(enrouteTime).Unix() &&
						u.EnRouteTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						enrouteTime = u.EnRouteTime
					}
				}
			}
			fmt.Printf("  - Enroute Time         : %s\n", enrouteTime)
		}

		{
			for _, u := range cs.Units {
				if strings.HasSuffix(u.Unit, "63") && u.Unit != "STA63" && u.Unit != "RES63" {
					if arrivalTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						arrivalTime = u.ArrivedTime
					}

					if monitor.CadDateTime(u.ArrivedTime).Unix() < monitor.CadDateTime(arrivalTime).Unix() &&
						u.ArrivedTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						arrivalTime = u.ArrivedTime
					}
				}
			}
			fmt.Printf("  - Arrival Time         : %s\n", arrivalTime)
		}
		{
			for _, u := range cs.Units {
				if strings.HasSuffix(u.Unit, "63") && u.Unit != "STA63" && u.Unit != "RES63" {
					if clearedTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", clearedTime, u.ClearedTime, u.Unit)
						clearedTime = u.ClearedTime
					}

					if monitor.CadDateTime(u.ClearedTime).Unix() < monitor.CadDateTime(clearedTime).Unix() &&
						monitor.CadDateTime(u.ClearedTime).Unix() > monitor.CadDateTime(arrivalTime).Unix() &&
						u.ClearedTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", clearedTime, u.ClearedTime, u.Unit)
						clearedTime = u.ClearedTime
					}
				}
			}
			fmt.Printf("  - Cleared Time         : %s\n", clearedTime)
		}

		dispatchToEnroute := monitor.CadDateTime(enrouteTime).Sub(monitor.CadDateTime(dispatchTime))
		fmt.Printf("  - Dispatch to Enroute  : %s\n", dispatchToEnroute.String())
		dispatchToArrival := monitor.CadDateTime(arrivalTime).Sub(monitor.CadDateTime(dispatchTime))
		fmt.Printf("  - Dispatch to Arrival  : %s\n", dispatchToArrival.String())
		timeOnScene := monitor.CadDateTime(clearedTime).Sub(monitor.CadDateTime(arrivalTime))
		fmt.Printf("  - Time on Scene        : %s\n", timeOnScene.String())

		fmt.Printf("  %8s: %19s|%19s|%19s|%19s\n", "Unit", "Dispatch", "En Route", "Arrival", "Cleared")
		for _, u := range cs.Units {
			fmt.Printf("  %8s: %19s|%19s|%19s|%19s\n", u.Unit, u.DispatchTime, u.EnRouteTime, u.ArrivedTime, u.ClearedTime)
		}
	}
}
