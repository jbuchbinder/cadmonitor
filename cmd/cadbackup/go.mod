module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.24.0

replace (
	github.com/jbuchbinder/cadmonitor => ../..
	github.com/jbuchbinder/cadmonitor/monitor => ../../monitor
)

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20250416130253-acada1016507
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/PuerkitoBio/goquery v1.11.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)
