module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.15

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/headzoo/ut v0.0.0-20181013193318-a13b5a7a02ca // indirect
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.11
)
