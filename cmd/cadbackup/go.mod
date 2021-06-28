module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.15

replace (
	github.com/jbuchbinder/cadmonitor => ../..
	github.com/jbuchbinder/cadmonitor/monitor => ../../monitor
)

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210511195108-30492af1bd31
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)
