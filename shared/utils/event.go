package utils

import (
	"encoding/json"
	"fmt"
	"github.com/delete-ui/NATS_TIRE_LIBRARY/shared/types"
	"github.com/google/uuid"
	"time"
)

type EventHelper struct {
	serviceName string
	version     string
}

func NewEventHelper(serviceName, version string) *EventHelper {
	return &EventHelper{
		serviceName: serviceName,
		version:     version,
	}
}

func (h *EventHelper) NewEventHeader(eventType types.EventType, correlationID int) types.EventHeader {
	return types.EventHeader{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		Timestamp:     time.Now().UTC(),
		Source:        h.serviceName,
		Version:       h.version,
		CorrelationID: correlationID,
	}
}

func (h *EventHelper) CreateMatchBundleEvent(bundle types.MatchBundle, correlationID int) types.MatchBundleEvent {
	return types.MatchBundleEvent{
		EventHeader: h.NewEventHeader(types.EventTypeMatchBundle, correlationID),
		Payload:     bundle,
	}
}

func (h *EventHelper) CreateMatchMonitoringEvent(match types.MatchMonitoring, correlationID int) types.MatchMonitoringEvent {
	return types.MatchMonitoringEvent{
		EventHeader: h.NewEventHeader(types.EventTypeMatchMonitoring, correlationID),
		Payload:     match,
	}
}

func (h *EventHelper) CreateForkFoundEvent(fork types.Fork, correlationID int) types.ForkFoundEvent {
	return types.ForkFoundEvent{
		EventHeader: h.NewEventHeader(types.EventTypeForkFound, correlationID),
		Payload:     fork,
	}
}

func (h *EventHelper) SerializeEvent(event interface{}) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}
	return data, nil
}

func (h *EventHelper) DeserializeEvent(data []byte, eventType types.EventType) (interface{}, error) {
	switch eventType {
	case types.EventTypeMatchBundle:
		var event types.MatchBundleEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return event, nil

	case types.EventTypeMatchMonitoring:
		var event types.MatchMonitoringEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return event, nil

	case types.EventTypeForkFound:
		var event types.ForkFoundEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return event, nil

	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
