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
	ClearedTime  string `json:"cleared_time"`
}

type Narrative struct {
	RecordedTime time.Time `json:"recorded_time"`
	Message      string    `json:"message"`
	User         string    `json:"user"`
}

type CallStatus struct {
	CallType     string                `json:"call_type"`
	NatureOfCall string                `json:"nature_of_call"`
	Location     string                `json:"location"`
	CrossStreets string                `json:"cross_streets"`
	Priority     int                   `json:"priority"`
	Narratives   []Narrative           `json:"narratives"`
	Units        map[string]UnitStatus `json:"units"`
}
