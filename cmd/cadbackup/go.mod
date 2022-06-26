module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.18

replace (
	github.com/jbuchbinder/cadmonitor => ../..
	github.com/jbuchbinder/cadmonitor/monitor => ../../monitor
)

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210628184336-b3eabebfc526
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.23.5
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
)
