module github.com/jbuchbinder/cadmonitor

go 1.15

replace github.com/jbuchbinder/cadmonitor/monitor => ./monitor

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.3.0
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7 // indirect
)
