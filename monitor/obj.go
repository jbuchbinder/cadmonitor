package monitor

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type UnitStatus struct {
	gorm.Model   `json:"-"`
	CallStatusID string `json:"cs_id" gorm:"uniqueIndex:unit_idx;size:64"`
	Unit         string `json:"unit" gorm:"size:32;uniqueIndex:unit_idx"`
	Status       string `json:"status"`
	DispatchTime string `json:"dispatch_time" gorm:"index"`
	EnRouteTime  string `json:"enroute_time" gorm:"index"`
	ArrivedTime  string `json:"arrived_time" gorm:"index"`
	ClearedTime  string `json:"cleared_time"`
}

type Narrative struct {
	gorm.Model   `json:"-"`
	CallStatusID string    `json:"cs_id" gorm:"uniqueIndex:narrative_idx;size:64"`
	RecordedTime time.Time `json:"recorded_time" gorm:"uniqueIndex:narrative_idx"`
	Message      string    `json:"message"`
	User         string    `json:"user"`
}

type Incident struct {
	gorm.Model     `json:"-"`
	CallStatusID   string `json:"cs_id" gorm:"index;size:64"`
	FDID           string `json:"fdid" gorm:"size:10;uniqueIndex:incident_fdid"`
	IncidentNumber string `json:"incident_number" gorm:"size:32;uniqueIndex:incident_fdid"`
}

type CallStatus struct {
	gorm.Model    `json:"-"`
	ID            string                `json:"id" gorm:"-"`
	CallID        string                `json:"call_id" gorm:"size:64,uniqueIndex:call_idx"`
	CallTime      time.Time             `json:"call_time" gorm:"uniqueIndex:call_idx"`
	DispatchTime  time.Time             `json:"dispatch_time" gorm:"index"`
	ArrivalTime   sql.NullTime          `json:"arrival_time"`
	CallType      string                `json:"call_type"`
	CallerPhone   string                `json:"caller_phone" gorm:"size:32"`
	NatureOfCall  string                `json:"nature_of_call" gorm:"index"`
	Location      string                `json:"location" gorm:"index"`
	District      string                `json:"district"`
	CrossStreets  string                `json:"cross_streets"`
	Priority      int                   `json:"priority" gorm:"index"`
	Incidents     []Incident            `json:"incidents" db:"-" gorm:"foreignKey:CallStatusID"`
	Narratives    []Narrative           `json:"narratives" db:"-" gorm:"foreignKey:CallStatusID"`
	Units         []UnitStatus          `json:"units" db:"-" gorm:"foreignKey:CallStatusID"`
	UnitStatusMap map[string]UnitStatus `json:"unit_status_map" db:"-" gorm:"-" sql:"-"`
	LastUpdated   time.Time             `json:"last_updated"`
	RawHTML       string                `json:"raw_html"`
}
