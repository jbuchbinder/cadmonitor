module github.com/jbuchbinder/cadmonitor/cmd/cadbackupingest

go 1.17

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210628184336-b3eabebfc526
	gorm.io/driver/mysql v1.2.1
	gorm.io/gorm v1.22.4
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
)
