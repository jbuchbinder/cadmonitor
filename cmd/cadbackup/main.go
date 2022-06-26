package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jbuchbinder/cadmonitor/monitor"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	debug       = flag.Bool("debug", false, "Enable debugging output")
	baseURL     = flag.String("baseUrl", "http://cadview.qvec.org/", "Base URL")
	beginDate   = flag.String("begin-date", "06/02/2019", "Begin date in MM/DD/YYYY format")
	endDate     = flag.String("end-date", "06/02/2019", "End date in MM/DD/YYYY format")
	fdid        = flag.String("fdid", "04042", "FDID")
	monitorType = flag.String("monitorType", "aegis", "Type of CAD system being monitored")
	suffix      = flag.String("suffix", "63", "Limit units to ones having a specfic suffix")
	dump        = flag.Bool("dump", false, "Dump results to terminal")
	diroutput   = flag.Bool("diroutput", false, "Directory output -- use 'database' flag to specify directory")
	database    = flag.String("db", "cadbackup.db", "CAD backup database (empty disables)")
)

func main() {
	flag.Parse()

	// db initialization
	var l logger.Interface
	l = logger.Default.LogMode(logger.Warn)
	if *debug {
		l = logger.Default
	}

	var db *gorm.DB
	var err error

	if !*diroutput {
		db, err = gorm.Open(sqlite.Open(*database), &gorm.Config{Logger: l})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&monitor.CallStatus{})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&monitor.Incident{})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&monitor.Narrative{})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&monitor.UnitStatus{})
		if err != nil {
			panic(err)
		}
	}

	if *diroutput {
		err = os.MkdirAll(*database, 0755)
		if err != nil {
			panic(err)
		}
	}

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
	err = m.Login(os.Getenv("CADUSER"), os.Getenv("CADPASS"))
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

	if *diroutput {
		return
	}

	for k, v := range calldata {
		if *debug {
			log.Printf("Inserting %s / %s : %#v", k, v.ID, v)
		} else {
			log.Printf("Inserting %s / %s", k, v.ID)

		}
		tx := db.Create(&v)
		if tx.Error != nil {
			log.Printf("ERROR: %s", tx.Error)
		}
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
		cs, err := m.GetStatusFromURL(v)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			continue
		}
		incidentNumber := lookupFdidFromIncidents(cs, *fdid)
		if *diroutput {
			ioutil.WriteFile(*database+string(os.PathSeparator)+incidentNumber, []byte(cs.RawHTML), 0644)
		}
		// Record call data
		log.Printf("Recording FDID %s Incident %s (%s)", *fdid, incidentNumber, cs.CallID)

		// Make sure that call ID is set univerally
		for _, v := range cs.Incidents {
			v.CallStatusID = cs.CallID
		}
		for _, v := range cs.Units {
			v.CallStatusID = cs.CallID
		}
		for k := range cs.Narratives {
			cs.Narratives[k].CallStatusID = cs.CallID
		}

		if *debug {
			b, _ := json.MarshalIndent(cs, "", "\t")
			fmt.Println(string(b))
		}
		calldata[cs.CallID] = cs
	}

	return calldata
}

func lookupFdidFromIncidents(cs monitor.CallStatus, fdid string) string {
	for _, i := range cs.Incidents {
		if i.FDID == fdid {
			return i.IncidentNumber
		}
	}
	return ""
}
