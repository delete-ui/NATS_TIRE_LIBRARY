package nats

import (
	"NATS_TIRE_LIBRARY/shared/constants"
	"NATS_TIRE_LIBRARY/shared/types"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

func (c *Client) SubscribeToMatchBundle(handler types.EventHandler) error {
	return c.subscribe(constants.TopicBundleMatch, handler, c.handleMatchBundleMessage)
}

func (c *Client) SubscribeToMatchMonitoring(handler types.EventHandler) error {
	return c.subscribe(constants.TopicMatchMonitoring, handler, c.handleMatchMonitoring)
}

func (c *Client) SubscribeToForkFound(handler types.EventHandler) error {
	return c.subscribe(constants.TopicForkFound, handler, c.handleForkFound)
}

func (c *Client) handleMatchBundleMessage(msg *nats.Msg, handler types.EventHandler) error {
	var event types.MatchBundleEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal match bundle event: %w", err)
	}

	if err := handler.HandleMatchBundleFound(event); err != nil {
		return fmt.Errorf("failed to handle match bundle event: %w", err)
	}

	if err := msg.Ack(); err != nil {
		return fmt.Errorf("failed to acknowledge match bundle event: %w", err)
	}

	c.logger.Debug("match event processed",
		zap.Int("match_id", event.Payload.CorrelationID),
		zap.String("event_id", event.EventHeader.EventID))

	return nil
}

func (c *Client) handleMatchMonitoring(msg *nats.Msg, handler types.EventHandler) error {
	var event types.MatchMonitoringEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal match monitoring event: %w", err)
	}

	if err := handler.HandleMatchMonitoring(event); err != nil {
		return fmt.Errorf("failed to handle match monitoring event: %w", err)
	}

	if err := msg.Ack(); err != nil {
		return fmt.Errorf("failed to acknowledge match monitoring event: %w", err)
	}

	c.logger.Debug("match event processed",
		zap.Int("match_id", event.Payload.CorrelationID),
		zap.String("event_id", event.EventHeader.EventID))

	return nil
}

func (c *Client) handleForkFound(msg *nats.Msg, handler types.EventHandler) error {
	var event types.ForkFoundEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal fork found event: %w", err)
	}

	if err := handler.HandleForkFound(event); err != nil {
		return fmt.Errorf("failed to handle fork found event: %w", err)
	}

	if err := msg.Ack(); err != nil {
		return fmt.Errorf("failed to acknowledge fork found event: %w", err)
	}

	c.logger.Debug("match event processed",
		zap.Int("match_id", event.Payload.CorrelationID),
		zap.String("event_id", event.EventHeader.EventID))

	return nil
}

func (c *Client) subscribe(subject string, handler types.EventHandler, msgHandler func(*nats.Msg, types.EventHandler) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return fmt.Errorf("NATS client is closed")
	}

	consumerConfig := &nats.ConsumerConfig{
		Durable:       c.config.ConsumerName,
		DeliverGroup:  c.config.ConsumerGroup,
		AckWait:       c.config.AckWait,
		MaxDeliver:    c.config.MaxDeliver,
		MaxAckPending: c.config.MaxAckPending,
		FilterSubject: subject,
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverNewPolicy,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	}

	_, err := c.jetStream.AddConsumer(c.config.StreamName, consumerConfig)
	if err != nil {
		if err == nats.ErrConsumerNameAlreadyInUse { // ✅ Правильно!
			c.logger.Warn("consumer already exists, using existing",
				zap.String("consumer", c.config.ConsumerName))
		} else {
			return fmt.Errorf("failed to add consumer: %w", err)
		}
	}

	sub, err := c.jetStream.PullSubscribe(
		subject,
		c.config.ConsumerName,
		[]nats.SubOpt{
			nats.Bind(c.config.StreamName, c.config.ConsumerName),
			nats.ManualAck(),
		}...,
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	c.handlers[subject] = handler
	c.subscriptions = append(c.subscriptions, sub)

	// Запускаем обработку сообщений в отдельной горутине
	c.wg.Add(1)
	go c.processSubscription(sub, subject, msgHandler)

	c.logger.Info("subscribed to subject",
		zap.String("subject", subject),
		zap.String("consumer", c.config.ConsumerName),
		zap.String("stream", c.config.StreamName))

	return nil
}

func (c *Client) processSubscription(sub *nats.Subscription, subject string, msgHandler func(*nats.Msg, types.EventHandler) error) {
	defer c.wg.Done()

	// Создаем тикер для пауз между итерациями
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debug("stopping subscription processing",
				zap.String("subject", subject))
			return

		case <-ticker.C:
			// Получаем пакет сообщений с таймаутом
			msgs, err := sub.Fetch(c.config.PullBatchSize, nats.MaxWait(5*time.Second))
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}
				// Если нет сообщений или другие ошибки - логируем и продолжаем
				if err.Error() != "nats: no messages" {
					c.logger.Debug("fetch messages result",
						zap.Error(err),
						zap.String("subject", subject))
				}
				continue
			}

			// Обрабатываем каждое сообщение
			for _, msg := range msgs {
				handler := c.getHandler(subject)
				if handler == nil {
					c.logger.Error("no handler found for subject",
						zap.String("subject", subject))
					msg.Ack() // Не блокируем очередь
					continue
				}

				if err := msgHandler(msg, handler); err != nil {
					c.logger.Error("failed to process message",
						zap.Error(err),
						zap.String("subject", subject))

					// Можно NAK сообщение для повторной обработки
					if nakErr := msg.Nak(); nakErr != nil {
						c.logger.Error("failed to nak message", zap.Error(nakErr))
					}
				}
			}
		}
	}
}

func (c *Client) getHandler(subject string) types.EventHandler {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.handlers[subject]
}

func (c *Client) Unsubscribe() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var lastErr error
	for _, sub := range c.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			lastErr = err
			c.logger.Error("failed to unsubscribe",
				zap.Error(err),
				zap.String("subject", sub.Subject))
		}
	}

	c.subscriptions = nil
	c.handlers = make(map[string]types.EventHandler)

	if lastErr != nil {
		return fmt.Errorf("failed to unsubscribe some subscriptions: %w", lastErr)
	}
	return nil
}
