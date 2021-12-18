module github.com/jbuchbinder/cadmonitor/cmd/cadbackup

go 1.15

replace (
	github.com/jbuchbinder/cadmonitor => ../..
	github.com/jbuchbinder/cadmonitor/monitor => ../../monitor
)

require (
	github.com/jbuchbinder/cadmonitor/monitor v0.0.0-20210628184336-b3eabebfc526
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.4
)
