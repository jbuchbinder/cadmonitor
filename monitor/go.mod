module github.com/jbuchbinder/cadmonitor/monitor

go 1.23.0

toolchain go1.24.0

replace github.com/jeffail/tunny => /opt/go-local/src/github.com/Jeffail/tunny

require (
	github.com/PuerkitoBio/goquery v1.10.2
	github.com/headzoo/surf v1.0.1
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1
	gorm.io/gorm v1.25.12
)

require (
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/headzoo/ut v0.0.0-20181013193318-a13b5a7a02ca // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
