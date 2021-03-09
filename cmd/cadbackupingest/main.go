package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jbuchbinder/cadmonitor/monitor"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	debug       = flag.Bool("debug", false, "Enable debugging output")
	dryrun      = flag.Bool("dryrun", false, "Run in dry mode, with no db commits")
	fdid        = flag.String("fdid", "04042", "FDID")
	monitorType = flag.String("monitorType", "aegis", "Type of CAD system being ingested")
	suffix      = flag.String("suffix", "63", "Limit units to ones having a specfic suffix")
	database    = flag.String("db", "cad:cad@/cad", "MySQL CAD backup database")
	backupdir   = flag.String("backupdir", "backup", "Read from backup directory")
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

	if !*dryrun {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: *database,
		}), &gorm.Config{Logger: l})
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

	m, err := monitor.InstantiateCadMonitor(*monitorType)
	if err != nil {
		panic(err)
	}
	err = m.ConfigureFromValues(map[string]string{
		"baseUrl": "",
		//"suffix":  *suffix, // -- disable to get every unit
		"fdid": *fdid,
	})
	if err != nil {
		panic(err)
	}

	dirents, err := os.ReadDir(*backupdir)
	if err != nil {
		panic(err)
	}
	for _, dirent := range dirents {
		if dirent.IsDir() {
			continue
		}
		log.Printf("Processing %s", dirent.Name())
		fullPath := *backupdir + string(os.PathSeparator) + dirent.Name()
		contents, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Printf("ERROR: Reading file %s: %s", dirent.Name(), err.Error())
			continue
		}
		status, err := m.GetStatus(contents, fullPath)

		if !*dryrun {
			tx := db.Create(&status)
			if tx.Error != nil {
				log.Printf("ERROR: %s", tx.Error)
			}
		}
	}
}
