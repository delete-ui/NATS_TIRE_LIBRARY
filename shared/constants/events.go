package constants

import "time"

// NATS TOPICS
const (
	TopicBundleMatch     = "events.bundle.match"
	TopicMatchMonitoring = "events.match.monitoring"
	TopicForkFound       = "events.fork.found"

	StreamEvents = "EVENTS"
)

// JetStream Configurations
const (
	MaxMessageAge = 24 * time.Hour
	Replicas      = 1
)

// Subscriber settings
const (
	DefaultAckWait = 30 * time.Second
	MaxDeliver     = 5
	MaxAckPending  = 256
	PullBatchSize  = 100
)

// Timeout settings
const (
	ConnectTimeout = 5 * time.Second
	RequestTimeout = 10 * time.Second
	ReconnectWait  = 1 * time.Second
	MaxReconnects  = -1
)

// Protocol version
const (
	ProtocolVersion = "1.0.0"
)
