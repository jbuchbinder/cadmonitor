module github.com/jbuchbinder/cadmonitor/cmd/cadbackupingest

go 1.17

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210511195108-30492af1bd31
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

require (
	github.com/PuerkitoBio/goquery v1.7.0 // indirect
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/headzoo/surf v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
)
