package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jbuchbinder/cadmonitor/monitor"
	rethinkdb "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var (
	baseURL           = flag.String("baseUrl", "http://cadview.qvec.org/", "Base URL")
	beginDate         = flag.String("begin-date", "06/02/2019", "Begin date in MM/DD/YYYY format")
	endDate           = flag.String("end-date", "06/02/2019", "End date in MM/DD/YYYY format")
	fdid              = flag.String("fdid", "04042", "FDID")
	monitorType       = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	suffix            = flag.String("suffix", "63", "Limit units to ones having a specfic suffix")
	dump              = flag.Bool("dump", false, "Dump results to terminal")
	rethinkURL        = flag.String("rethinkdb-url", "10.1.1.60:28015", "Rethink DB URL (empty disables)")
	rethinkDB         = flag.String("rethinkdb-db", "STA63_CALLS", "Name of database")
	rethinkTable      = flag.String("rethinkdb-table", "calls", "Name of table")
	rethinkClearTable = flag.Bool("rethinkdb-clear-table", false, "Clear table before inserting")
)

// CallData represents a single call
type CallData struct {
	CallStatus        monitor.CallStatus
	ID                string `rethinkdb:"id"`
	DispatchTime      string
	EnrouteTime       string
	ArrivalTime       string
	ClearedTime       string
	DispatchToEnroute time.Duration
	DispatchToArrival time.Duration
	TimeOnScene       time.Duration
}

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

	calldata := map[string]CallData{}

	startDate := monitor.CadDateTimeShort(*beginDate)
	stopDate := monitor.CadDateTimeShort(*endDate)
	log.Printf("Range : %s - %s", startDate.Format(time.RFC1123), stopDate.Format(time.RFC1123))
	thisDt := startDate
	for {
		if thisDt.After(stopDate) {
			break
		}
		log.Printf("fetchCallsForDate %s", thisDt.Format("01/02/2006"))
		callsSlice := fetchCallsForDate(m, thisDt.Format("01/02/2006"))
		for k, v := range callsSlice {
			calldata[k] = v
		}
		thisDt = thisDt.Add(time.Hour * time.Duration(24))
	}

	if *dump {
		b, _ := json.MarshalIndent(calldata, "", "  ")
		fmt.Println(string(b))
	}

	if *rethinkURL != "" {
		session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
			Address:  *rethinkURL,
			Database: *rethinkDB,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// Basic setup -- ignore errors
		rethinkdb.DBCreate(*rethinkDB)
		if *rethinkClearTable {
			rethinkdb.TableDrop(*rethinkTable)
		}
		rethinkdb.TableCreate(*rethinkTable)

		for k, v := range calldata {
			v.ID = v.CallStatus.CallID
			log.Printf("rethinkdb: Inserting %s / %s", k, v.CallStatus.CallID)
			err := rethinkdb.Table(*rethinkTable).Insert(v).Exec(session)
			if err != nil {
				log.Printf("%s", err.Error())
			}
		}
	}
}

func fetchCallsForDate(m monitor.CadMonitor, dt string) map[string]CallData {
	var calls map[string]string
	calldata := map[string]CallData{}
	calls, err := m.GetClearedCalls(dt)
	if err != nil {
		panic(err)
	}
	for _, v := range calls {
		cs, err := m.GetStatus(v)
		cd := CallData{CallStatus: cs}
		if err != nil {
			panic(err)
		}
		//fmt.Println("-------------------------------------------------------------------------------------------")
		//fmt.Printf("CALL %s (PRI %d @ %s)\n", k, cd.CallStatus.Priority, cd.CallStatus.Location)

		// Dispatch time is first arrived time from units
		{
			for _, u := range cs.Units {
				if u.Unit == *suffix+"FAST" || u.Unit == "STA"+*suffix || u.Unit == "RES"+*suffix {
					if cd.DispatchTime == "" {
						cd.DispatchTime = u.DispatchTime
					}

					if monitor.CadDateTime(u.DispatchTime).Unix() < monitor.CadDateTime(cd.DispatchTime).Unix() &&
						u.DispatchTime != "" {
						cd.DispatchTime = u.DispatchTime
					}
				}
			}
		}

		// Arrived/Cleared time is first arrived time from units
		{
			for _, u := range cs.Units {
				if strings.HasSuffix(u.Unit, *suffix) {
					if cd.EnrouteTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						cd.EnrouteTime = u.EnRouteTime
					}

					if monitor.CadDateTime(u.EnRouteTime).Unix() < monitor.CadDateTime(cd.EnrouteTime).Unix() &&
						u.EnRouteTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						cd.EnrouteTime = u.EnRouteTime
					}
				}
			}
		}

		{
			for _, u := range cd.CallStatus.Units {
				if strings.HasSuffix(u.Unit, *suffix) {
					if cd.ArrivalTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						cd.ArrivalTime = u.ArrivedTime
					}

					if monitor.CadDateTime(u.ArrivedTime).Unix() < monitor.CadDateTime(cd.ArrivalTime).Unix() &&
						u.ArrivedTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", arrivalTime, u.ArrivedTime, u.Unit)
						cd.ArrivalTime = u.ArrivedTime
					}
				}
			}
		}
		{
			for _, u := range cd.CallStatus.Units {
				if strings.HasSuffix(u.Unit, *suffix) {
					if cd.ClearedTime == "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", clearedTime, u.ClearedTime, u.Unit)
						cd.ClearedTime = u.ClearedTime
					}

					if monitor.CadDateTime(u.ClearedTime).Unix() < monitor.CadDateTime(cd.ClearedTime).Unix() &&
						monitor.CadDateTime(u.ClearedTime).Unix() > monitor.CadDateTime(cd.ArrivalTime).Unix() &&
						u.ClearedTime != "" {
						//fmt.Printf("Old time : %s, new time : %s, unit : %s\n", clearedTime, u.ClearedTime, u.Unit)
						cd.ClearedTime = u.ClearedTime
					}
				}
			}
		}

		cd.DispatchToEnroute = monitor.CadDateTime(cd.EnrouteTime).Sub(monitor.CadDateTime(cd.DispatchTime))
		cd.DispatchToArrival = monitor.CadDateTime(cd.ArrivalTime).Sub(monitor.CadDateTime(cd.DispatchTime))
		cd.TimeOnScene = monitor.CadDateTime(cd.ClearedTime).Sub(monitor.CadDateTime(cd.ArrivalTime))

		// Record call data
		//log.Printf("Recording %s", cd.CallStatus.CallID)
		//b, _ := json.MarshalIndent(cd, "", "\t")
		//fmt.Println(string(b))
		calldata[cd.CallStatus.CallID] = cd
	}

	return calldata
}
