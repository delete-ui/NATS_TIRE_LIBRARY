package types

import "time"

type EventType string

const (
	EventTypeMatchBundle     EventType = "match.bundle"
	EventTypeMatchMonitoring EventType = "match.monitoring"
	EventTypeForkFound       EventType = "fork.found"
)

type EventHeader struct {
	EventID       string    `json:"event_id"`
	EventType     EventType `json:"event_type"`
	Timestamp     time.Time `json:"timestamp"`
	Source        string    `json:"source"`
	Version       string    `json:"version"`
	CorrelationID int       `json:"correlation_id"`
}

type MatchBundleEvent struct {
	EventHeader EventHeader `json:"event_header"`
	Payload     MatchBundle `json:"payload"`
}

type MatchMonitoringEvent struct {
	EventHeader EventHeader     `json:"event_header"`
	Payload     MatchMonitoring `json:"payload"`
}

type ForkFoundEvent struct {
	EventHeader EventHeader `json:"event_header"`
	Payload     Fork        `json:"payload"`
}
