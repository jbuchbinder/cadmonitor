package main

import (
	"time"
)

type UnitStatus struct {
	Unit         string
	Status       string
	DispatchTime string
	EnRouteTime  string
	ArrivedTime  string
}

type Narrative struct {
	RecordedTime time.Time
	Message      string
	User         string
}

type CallStatus struct {
	Narratives []Narrative
	Units      map[string]UnitStatus
}
