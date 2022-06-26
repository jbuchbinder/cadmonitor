module github.com/jbuchbinder/cadmonitor/monitor

go 1.18

replace github.com/jeffail/tunny => /opt/go-local/src/github.com/Jeffail/tunny

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/headzoo/surf v1.0.1
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1
	gorm.io/gorm v1.23.5
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/headzoo/ut v0.0.0-20181013193318-a13b5a7a02ca // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
)
