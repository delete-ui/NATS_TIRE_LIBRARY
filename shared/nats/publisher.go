package nats

import (
	"NATS_TIRE_LIBRARY/shared/constants"
	"NATS_TIRE_LIBRARY/shared/types"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

func (c *Client) PublishMatchBundle(bundle types.MatchBundle) error {
	return c.publishEvent(constants.TopicBundleMatch, bundle, types.EventTypeMatchBundle)
}

func (c *Client) PublishMatchMonitoring(monitoring types.MatchMonitoring) error {
	return c.publishEvent(constants.TopicMatchMonitoring, monitoring, types.EventTypeMatchMonitoring)
}

func (c *Client) PublishForkFound(fork types.Fork) error {
	return c.publishEvent(constants.TopicForkFound, fork, types.EventTypeForkFound)
}

func (c *Client) publishEvent(topic string, payload interface{}, eventType types.EventType) error {
	var event interface{}

	// Создаем соответствующее событие
	switch eventType {
	case types.EventTypeMatchBundle:
		match, ok := payload.(types.MatchBundle)
		if !ok {
			return fmt.Errorf("invalid payload type for match event")
		}
		event = c.events.CreateMatchBundleEvent(match, 0)

	case types.EventTypeMatchMonitoring:
		odds, ok := payload.(types.MatchMonitoring)
		if !ok {
			return fmt.Errorf("invalid payload type for odds event")
		}
		event = c.events.CreateMatchMonitoringEvent(odds, 0)

	case types.EventTypeForkFound:
		fork, ok := payload.(types.Fork)
		if !ok {
			return fmt.Errorf("invalid payload type for fork event")
		}
		event = c.events.CreateForkFoundEvent(fork, 0)

	default:
		return fmt.Errorf("unsupported event type: %s", eventType)
	}

	// Сериализуем событие
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Публикуем
	return c.publish(topic, data)
}

func (c *Client) publish(subject string, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isClosed {
		return fmt.Errorf("client is  closed")
	}

	ack, err := c.jetStream.Publish(subject, data)
	if err != nil {
		c.logger.Error("failed to publish message",
			zap.Error(err),
			zap.String("subject", subject))
		return err
	}

	c.logger.Debug("message published successfully",
		zap.String("subject", subject),
		zap.String("stream", ack.Stream),
		zap.Uint64("sequence", ack.Sequence))

	return nil
}
