module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.23.0

toolchain go1.24.0

replace (
	github.com/jbuchbinder/cadmonitor => ../..
	github.com/jbuchbinder/cadmonitor/monitor => ../../monitor
)

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20230220185920-fbf01137b3f4
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/PuerkitoBio/goquery v1.10.2 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
