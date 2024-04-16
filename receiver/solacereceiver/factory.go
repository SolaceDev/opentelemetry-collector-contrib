// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package solacereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver/internal/metadata"
)

const (
	// default value for max unaked messages
	defaultMaxUnaked int32 = 1000
	// default value for host
	defaultHost string = "localhost:5671"
)

// NewFactory creates a factory for Solace receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, metadata.TracesStability),
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha),
	)
}

var receivers = make(map[*Config]*solaceReceiver)
var receiversLock = new(sync.Mutex)

// createDefaultConfig creates the default configuration for receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Broker:     []string{defaultHost},
		MaxUnacked: defaultMaxUnaked,
		Auth:       Authentication{},
		TLS: configtls.ClientConfig{
			InsecureSkipVerify: false,
			Insecure:           false,
		},
		Flow: FlowControl{
			DelayedRetry: &FlowControlDelayedRetry{
				Delay: 10 * time.Millisecond,
			},
		},
	}
}

// CreateTracesReceiver creates a trace receiver based on provided config. Component is not shared
func createTracesReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	receiverConfig component.Config,
	nextConsumer consumer.Traces,
) (receiver.Traces, error) {
	cfg, ok := receiverConfig.(*Config)
	if !ok {
		return nil, component.ErrDataTypeIsNotSupported
	}

	receiversLock.Lock()
	defer receiversLock.Unlock()
	r := receivers[cfg]
	if r == nil {
		var err error
		r, err = newReceiver(cfg, params)
		if err != nil {
			return nil, err
		}
		receivers[cfg] = r
	}
	r.setTraceConsumer(nextConsumer)
	// pass cfg, params and next consumer through
	return r, nil
}

func createLogsReceiver(_ context.Context, params receiver.CreateSettings,
	receiverConfig component.Config, nextConsumer consumer.Logs) (receiver.Logs, error) {
	cfg, ok := receiverConfig.(*Config)
	if !ok {
		return nil, component.ErrDataTypeIsNotSupported
	}

	receiversLock.Lock()
	defer receiversLock.Unlock()
	r := receivers[cfg]
	if r == nil {
		var err error
		r, err = newReceiver(cfg, params)
		if err != nil {
			return nil, err
		}
		receivers[cfg] = r
	}
	r.setLogsConsumer(nextConsumer)
	// pass cfg, params and next consumer through
	return r, nil
}
