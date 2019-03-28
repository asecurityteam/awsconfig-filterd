package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/asecurityteam/awsconfig-filterd/pkg/logs"
	"github.com/aws/aws-sdk-go/service/configservice"
)

// ConfigNotification is the basic structure of all incoming AWS Config SNS events
// (https://docs.aws.amazon.com/config/latest/developerguide/example-sns-notification.html).
type ConfigNotification struct {
	Message            string `json:"Message"`
	Timestamp          string `json:"Timestamp"`
	ProcessedTimestamp string `json:"ProcessedTimestamp"`
	Type               string `json:"Type"`
	MessageID          string `json:"MessageId"`
	TopicArn           string `json:"TopicArn"`
	Subject            string `json:"Subject"`
	SignatureVersion   string `json:"SignatureVersion"`
	Signature          string `json:"Signature"`
	SigningCertURL     string `json:"SigningCertURL"`
	UnsubscribeURL     string `json:"UnsubscribeURL"`
}

// ConfigEvent represents a single AWS Config Configuration Item.
type ConfigEvent struct {
	ConfigurationItem domain.ConfigurationItem `json:"configurationItem"`
	MessageType       string                   `json:"messageType"`
}

// ConfigFilterHandler applies a filter to AWS Config events.
type ConfigFilterHandler struct {
	ConfigFilterer domain.ConfigFilterer
	LogFn          domain.LogFn
	StatFn         domain.StatFn
}

// Handle accepts Config events, applies a filter, and returns the events that match.
func (h *ConfigFilterHandler) Handle(ctx context.Context, in ConfigNotification) (ConfigNotification, error) {
	logger := h.LogFn(ctx)
	stater := h.StatFn(ctx)

	var event ConfigEvent
	if in.Message == "" {
		return ConfigNotification{}, nil
	}

	e := json.Unmarshal([]byte(in.Message), &event)
	if e != nil {
		logger.Error(logs.InvalidInput{Reason: e.Error()})
		return ConfigNotification{}, e
	}

	if event.MessageType != configservice.MessageTypeConfigurationItemChangeNotification {
		return ConfigNotification{}, nil
	}

	if event.ConfigurationItem.ResourceType == "" {
		e = domain.ErrInvalidInput{Reason: "empty resource type"}
		logger.Error(logs.InvalidInput{Reason: e.Error()})
		return ConfigNotification{}, e
	}
	stater.Count("event.awsconfig.filter.resource_type", 1,
		fmt.Sprintf("type:%s", event.ConfigurationItem.ResourceType))

	if ok := h.ConfigFilterer.FilterConfig(event.ConfigurationItem); !ok {
		stater.Count("event.awsconfig.filter.discarded", 1)
		return ConfigNotification{}, nil
	}

	if ts, err := time.Parse(time.RFC3339Nano, in.Timestamp); err == nil {
		stater.Timing("event.awsconfig.filter.event.delay", time.Since(ts))
	}
	stater.Count("event.awsconfig.filter.accepted", 1)
	in.ProcessedTimestamp = time.Now().Format(time.RFC3339Nano)
	return in, nil
}
