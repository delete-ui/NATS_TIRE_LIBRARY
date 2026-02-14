package types

// EventHandler
type EventHandler interface {
	HandleMatchBundleFound(event MatchBundleEvent) error
	HandleMatchMonitoring(event MatchMonitoringEvent) error
	HandleForkFound(event ForkFoundEvent) error
}

// Publisher
type Publisher interface {
	PublishMatchBundle(bundle MatchBundle) error
	PublishMatchMonitoring(monitoring MatchMonitoring) error
	PublishForkFound(fork Fork) error
}

// Consumer
type Consumer interface {
	SubscribeToMatchBundle(handler EventHandler) error
	SubscribeToMatchMonitoring(handler EventHandler) error
	SubscribeToForkFound(handler EventHandler) error
}
