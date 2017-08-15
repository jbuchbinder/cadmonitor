# QVEC CAD MONITOR

[![Build Status](https://secure.travis-ci.org/jbuchbinder/qveccadmonitor.png)](http://travis-ci.org/jbuchbinder/qveccadmonitor)


[QVEC](http://qvec.org) CAD system monitor for apparatuses

## Development

The default `secrets.go` file is encrypted by the developer, but you can create your own like this:

	package main
	
	const (
	        USER = "username"
	        PASS = "password"
	)

