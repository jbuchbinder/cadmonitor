package monitor

import (
	"time"
)

const (
	URLPREFIX = "http://cadview.qvec.org/NewWorld.CAD.ViewOnly/"
)

type UnitStatus struct {
	Unit         string `json:"unit"`
	Status       string `json:"status"`
	DispatchTime string `json:"dispatch_time"`
	EnRouteTime  string `json:"enroute_time"`
	ArrivedTime  string `json:"arrived_time"`
}

type Narrative struct {
	RecordedTime time.Time `json:"recorded_time"`
	Message      string    `json:"message"`
	User         string    `json:"user"`
}

type CallStatus struct {
	Narratives []Narrative           `json:"narratives"`
	Units      map[string]UnitStatus `json:"units"`
}
