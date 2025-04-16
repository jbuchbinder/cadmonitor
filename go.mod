module github.com/jbuchbinder/cadmonitor

go 1.23.0

toolchain go1.24.0

replace github.com/jbuchbinder/cadmonitor/monitor => ./monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20250327183818-e07f6f02833a
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
