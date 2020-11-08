package monitor

import (
	"time"
)

type UnitStatus struct {
	CallStatusID string `json:"cs_id"`
	Unit         string `json:"unit"`
	Status       string `json:"status"`
	DispatchTime string `json:"dispatch_time"`
	EnRouteTime  string `json:"enroute_time"`
	ArrivedTime  string `json:"arrived_time"`
	ClearedTime  string `json:"cleared_time"`
}

type Narrative struct {
	CallStatusID string    `json:"cs_id"`
	RecordedTime time.Time `json:"recorded_time"`
	Message      string    `json:"message"`
	User         string    `json:"user"`
}

type Incident struct {
	CallStatusID   string `json:"cs_id"`
	FDID           string `json:"fdid"`
	IncidentNumber string `json:"incident_number"`
}

type CallStatus struct {
	ID           string                `json:"id"`
	CallID       string                `json:"call_id"`
	CallTime     time.Time             `json:"call_time"`
	DispatchTime time.Time             `json:"dispatch_time"`
	ArrivalTime  time.Time             `json:"arrival_time"`
	CallType     string                `json:"call_type"`
	CallerPhone  string                `json:"caller_phone"`
	NatureOfCall string                `json:"nature_of_call"`
	Location     string                `json:"location"`
	District     string                `json:"district"`
	CrossStreets string                `json:"cross_streets"`
	Priority     int                   `json:"priority"`
	Incidents    []Incident            `json:"incidents" db:"-"`
	Narratives   []Narrative           `json:"narratives" db:"-"`
	Units        map[string]UnitStatus `json:"units" db:"-"`
	LastUpdated  time.Time             `json:"last_updated"`
	RawHTML      string                `json:"raw_html"`
}
