package nats

import (
	"NATS_TIRE_LIBRARY/shared/config"
	"NATS_TIRE_LIBRARY/shared/types"
	"fmt"
)

func Factory(cfg *config.Config) (types.Publisher, types.Consumer, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create NATS client: %w", err)
	}

	return client, client, nil
}

func NewPublisher(cfg *config.Config) (types.Publisher, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewConsumer(cfg *config.Config) (types.Consumer, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}
