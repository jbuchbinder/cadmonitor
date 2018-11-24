# QVEC CAD MONITOR

[![Build Status](https://secure.travis-ci.org/jbuchbinder/cadmonitor.png)](http://travis-ci.org/jbuchbinder/cadmonitor)

Emergency services CAD system monitor for apparatuses. Originally designed to interface with [QVEC](http://qvec.org)'s Aegis system

## Development

The default `secrets.go` file is encrypted by the developer, but you can create your own like this:

	package main
	
	const (
	        USER = "username"
	        PASS = "password"
	)

These are for unit tests, and do not affect the functionality of the actual library. 

