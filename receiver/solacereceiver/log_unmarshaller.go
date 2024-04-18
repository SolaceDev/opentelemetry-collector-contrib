// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package solacereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver"

import (
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	logs_v1 "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver/internal/model/logs/v1"
)

const logTopicPrefix = "_telemetry/broker/logs/"

// tracesUnmarshaller deserializes the message body.
type logsUnmarshaller interface {
	// unmarshal the amqp-message into traces.
	// Only valid traces are produced or error is returned
	unmarshal(message *inboundMessage) (plog.Logs, error)
}

// newUnmarshalleer returns a new unmarshaller ready for message unmarshalling
func newLogsUnmarshaller(logger *zap.Logger, metrics *opencensusMetrics) logsUnmarshaller {
	return &solaceLogsUnmarshaller{
		logger:  logger,
		metrics: metrics,
		// v1 unmarshaller is implemented by solaceMessageUnmarshallerV1
		receiveUnmarshallerV1: &brokerTraceReceiveUnmarshallerV1{
			logger:  logger,
			metrics: metrics,
		},
		egressUnmarshallerV1: &brokerTraceEgressUnmarshallerV1{
			logger:  logger,
			metrics: metrics,
		},
	}
}

// solaceTracesUnmarshaller implements tracesUnmarshaller.
type solaceLogsUnmarshaller struct {
	logger                *zap.Logger
	metrics               *opencensusMetrics
	receiveUnmarshallerV1 tracesUnmarshaller
	egressUnmarshallerV1  tracesUnmarshaller
}

// unmarshal will unmarshal an *solaceMessage into ptrace.Traces.
// It will make a decision based on the version of the message which unmarshalling strategy to use.
// For now, only receive v1 messages are used.
func (u *solaceLogsUnmarshaller) unmarshal(message *inboundMessage) (plog.Logs, error) {
	const (
		logTopic = "_telemetry/broker/logs/v1"
	)
	if message.Properties == nil || message.Properties.To == nil {
		// no topic
		u.logger.Error("Received message with no topic")
		return plog.Logs{}, errUnknownTopic
	}
	var topic string = *message.Properties.To
	// Multiplex the topic string. For now we only have a single type handled
	if strings.HasPrefix(topic, logTopic) {
		// we can unmarshal
		logData, err := u.unmarshalToLogData(message)
		if err != nil {
			return plog.Logs{}, err
		}
		// logData := getDummyData()
		logs := plog.NewLogs()
		u.populateLogs(logData, logs)
		return logs, nil
	}
	// unknown topic, do not require an upgrade
	u.logger.Error("Received message with unknown topic", zap.String("topic", *message.Properties.To))
	return plog.Logs{}, errUnknownTopic
}

func (u *solaceLogsUnmarshaller) unmarshalToLogData(message *inboundMessage) (*logs_v1.LogsData, error) {
	var data = message.GetData()
	if len(data) == 0 {
		return nil, errEmptyPayload
	}
	var logRecord logs_v1.LogsData
	if err := proto.Unmarshal(data, &logRecord); err != nil {
		return nil, err
	}
	return &logRecord, nil
}

func (u *solaceLogsUnmarshaller) populateLogs(logsData *logs_v1.LogsData, logs plog.Logs) {
	for _, resourceLogData := range logsData.GetResourceLogs() {
		resourceLogs := logs.ResourceLogs().AppendEmpty()
		u.populateAttrs(resourceLogData.Resource.Attributes, resourceLogs.Resource().Attributes())
		for _, scopeLogData := range resourceLogData.ScopeLogs {
			scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
			for _, logRecordData := range scopeLogData.LogRecords {
				logRecord := scopeLogs.LogRecords().AppendEmpty()
				logRecord.SetTimestamp(pcommon.Timestamp(logRecordData.TimeUnixNano))
				logRecord.SetObservedTimestamp(pcommon.NewTimestampFromTime(time.Now()))
				logRecord.SetSeverityNumber(plog.SeverityNumber(logRecordData.SeverityNumber))
				// TODO check to make sure this is a string
				if logRecordData.Body != nil && logRecordData.Body.Value != nil {
					logRecord.Body().SetStr(logRecordData.Body.GetStringValue())
				} else {
					u.logger.Warn("no body")
				}
				u.populateAttrs(logRecordData.Attributes, logRecord.Attributes())
			}
		}
	}
}

func (u *solaceLogsUnmarshaller) populateAttrs(attributesFrom []*logs_v1.KeyValue, attributesTo pcommon.Map) {
	for _, attr := range attributesFrom {
		switch casted := attr.Value.Value.(type) {
		case *logs_v1.AnyValue_StringValue:
			attributesTo.PutStr(attr.Key, casted.StringValue)
		case *logs_v1.AnyValue_BoolValue:
			attributesTo.PutBool(attr.Key, casted.BoolValue)
		case *logs_v1.AnyValue_IntValue:
			attributesTo.PutInt(attr.Key, casted.IntValue)
		case *logs_v1.AnyValue_DoubleValue:
			attributesTo.PutDouble(attr.Key, casted.DoubleValue)
		default:
			u.logger.Warn("Unknown attribute type", zap.String("type", fmt.Sprintf("%t", attr.Value.Value)))
			attributesTo.PutEmpty(attr.Key)
		}
	}
}

// func getDummyData() *logs_v1.LogRecord {
// 	return &logs_v1.LogRecord{
// 		TimeUnixNano: uint64(time.Now().UnixNano()),
// 		SeverityNumber: 1,
// 		Body: &logs_v1.AnyValue{
// 			Value: &logs_v1.AnyValue_StringValue{
// 				StringValue: "CLIENT_CLIENT_CONNECT",
// 			},
// 		},
// 		// attributes:
// 		//   solace.event.scope --> string: "CLIENT"
// 		//   solace.event.module --> string: "CLIENT"
// 		//   solace.event.name --> string: "CONNECT"
// 		//   solace.broker.name --> string: "vmr-133-9" (appears in logs today after the severity string)
// 		//   solace.broker.message-vpn --> string: "default" (ex; same as how we include VPN in existing logs)
// 		//   solace.client.name --> string: "#client" (ex; same as how we include client name today)
// 		//   solace.client.id --> int64: 1234
// 		//   solace.client.name --> string: "#client"
// 		//   solace.client.clientUsername --> string: "default"
// 		//   solace.client.originalClientUsername --> string: "default"
// 		//   solace.client.webSessionId --> string: "12345"
// 		//   solace.client.applianceAddr --> string: "192.168.133.9"
// 		//   solace.client.clientAddr --> string: "192.168.0.1"
// 		//   solace.client.capabilities --> string: "cap1,cap2,cap3"
// 		//   solace.client.tlsValidityNotAfter --> string: "2022-01-01T00:00:00Z"
// 		Attributes: []*logs_v1.KeyValue{
// 			{
// 				Key: "solace.event.scope",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "CLIENT",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.event.module",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "CLIENT",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.event.name",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "CONNECT",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.broker.name",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "vmr-133-9",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.broker.message-vpn",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "default",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.name",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "#client",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.id",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_IntValue{
// 						IntValue: 1234,
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.clientUsername",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "default",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.originalClientUsername",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "default",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.webSessionId",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "12345",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.applianceAddr",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "192.168.133.9",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.clientAddr",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "192.168.0.1",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.capabilities",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "cap1,cap2,cap3",
// 					},
// 				},
// 			},
// 			{
// 				Key: "solace.client.tlsValidityNotAfter",
// 				Value: &logs_v1.AnyValue{
// 					Value: &logs_v1.AnyValue_StringValue{
// 						StringValue: "2022-01-01T00:00:00Z",
// 					},
// 				},
// 			},
// 		},
// 	}
// }


// time_unix_nano: <64-bit ns timestamp since Epoch>
// severity_number: SEVERITY_NUMBER_INFO
// body: string: "CLIENT_CLIENT_CONNECT" (mgmtEvent::eventIdToString(eventId))
// attributes:
//   solace.event.scope --> string: "CLIENT"
//   solace.event.module --> string: "CLIENT"
//   solace.event.name --> string: "CONNECT"
//   solace.broker.name --> string: "vmr-133-9" (appears in logs today after the severity string)
//   solace.broker.message-vpn --> string: "default" (ex; same as how we include VPN in existing logs)
//   solace.client.name --> string: "#client" (ex; same as how we include client name today)
//   solace.client.id --> int64: 1234
//   solace.client.name --> string
//   solace.client.clientUsername --> string
//   solace.client.originalClientUsername --> string
//   solace.client.webSessionId --> string
//   solace.client.applianceAddr --> string
//   solace.client.clientAddr --> string
//   solace.client.capabilities --> string
//   solace.client.tlsValidityNotAfter --> string
// * See Rules below