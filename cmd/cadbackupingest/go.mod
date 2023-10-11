module github.com/jbuchbinder/cadmonitor/cmd/cadbackupingest

go 1.18

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20220626150718-6edcab5606c8
	gorm.io/driver/mysql v1.4.7
	gorm.io/gorm v1.24.5
)

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.17.0 // indirect
)
