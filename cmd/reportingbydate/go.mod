module github.com/jbuchbinder/cadmonitor/cmd/reportingbydate

go 1.16

replace github.com/jbuchbinder/cadmonitor/monitor => ../../monitor

require (
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-00010101000000-000000000000
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/rethinkdb/rethinkdb-go.v5 v5.1.0
)
