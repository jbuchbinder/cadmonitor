module github.com/jbuchbinder/cadmonitor/cmd/cadbackupingest

go 1.17

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-00010101000000-000000000000
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.21.3
)
