package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jbuchbinder/cadmonitor/monitor"
)

var (
	baseURL     = flag.String("baseUrl", "http://cadview.qvec.org/", "Base URL")
	beginDate   = flag.String("begin-date", "06/02/2019", "Begin date in MM/DD/YYYY format")
	endDate     = flag.String("end-date", "06/02/2019", "End date in MM/DD/YYYY format")
	fdid        = flag.String("fdid", "04042", "FDID")
	monitorType = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	suffix      = flag.String("suffix", "63", "Limit units to ones having a specfic suffix")
	dump        = flag.Bool("dump", false, "Dump results to terminal")
	mysqlHost   = flag.String("mysql-host", "127.0.0.1", "MySQL DB Host (empty disables)")
	mysqlDB     = flag.String("mysql-db", "cadbackup", "Name of database")
)

func main() {
	flag.Parse()

	m, err := monitor.InstantiateCadMonitor(*monitorType)
	if err != nil {
		panic(err)
	}
	err = m.ConfigureFromValues(map[string]string{
		"baseUrl": *baseURL,
		//"suffix":  *suffix, // -- disable to get every unit
		"fdid": *fdid,
	})
	if err != nil {
		panic(err)
	}
	err = m.Login(monitor.USER, monitor.PASS)
	if err != nil {
		panic(err)
	}

	calldata := map[string]monitor.CallStatus{}

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

	for k, v := range calldata {
		log.Printf("Inserting %s / %s : %#v", k, v.ID, v)
		//err := rethinkdb.Table(*rethinkTable).Insert(v).Exec(session)
	}
}

func fetchCallsForDate(m monitor.CadMonitor, dt string) map[string]monitor.CallStatus {
	var calls map[string]string
	calldata := map[string]monitor.CallStatus{}
	calls, err := m.GetClearedCalls(dt)
	if err != nil {
		panic(err)
	}
	for _, v := range calls {
		cs, err := m.GetStatus(v)
		if err != nil {
			panic(err)
		}
		// Record call data
		//log.Printf("Recording %s", cd.CallStatus.CallID)
		//b, _ := json.MarshalIndent(cd, "", "\t")
		//fmt.Println(string(b))
		calldata[cs.CallID] = cs
	}

	return calldata
}
