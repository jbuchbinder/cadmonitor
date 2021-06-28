module github.com/jbuchbinder/cadmonitor

go 1.15

replace github.com/jbuchbinder/cadmonitor/monitor => ./monitor

require (
	github.com/PuerkitoBio/goquery v1.7.0 // indirect
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210511195108-30492af1bd31
	github.com/joho/godotenv v1.3.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	gorm.io/gorm v1.21.11 // indirect
)
