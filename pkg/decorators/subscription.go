package decorators

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/asecurityteam/awsconfig-filterd/pkg/logs"
	"github.com/asecurityteam/components"
	"github.com/asecurityteam/runhttp"
)

const typeNotification = "SubscriptionConfirmation"

// SubscriptionConfig contains settings for the SubscriptionComponent.
type SubscriptionConfig struct {
	HTTP *components.HTTPConfig
}

// Name of the configuration root.
func (*SubscriptionConfig) Name() string {
	return decoratorSubscription
}

// SubscriptionComponent is the component for the subscription decorator
type SubscriptionComponent struct {
	HTTP *components.HTTPComponent
}

// NewSubscriptionComponent generates a SubscriptionComponent.
func NewSubscriptionComponent() *SubscriptionComponent {
	return &SubscriptionComponent{
		HTTP: components.NewHTTPComponent(),
	}
}

// Settings generates the default configuration.
func (c *SubscriptionComponent) Settings() *SubscriptionConfig {
	return &SubscriptionConfig{
		HTTP: c.HTTP.Settings(),
	}
}

// New generates a Subcription decorator.
func (c *SubscriptionComponent) New(ctx context.Context, conf *SubscriptionConfig) (*Subscription, error) {
	rt, e := c.HTTP.New(ctx, conf.HTTP)
	if e != nil {
		return nil, e
	}
	return &Subscription{
		LogFn: runhttp.LoggerFromContext,
		Client: &http.Client{
			Transport: rt,
		},
	}, nil
}

// Subscription is a lambda decorator which will check for S=subscription confirmation messages
type Subscription struct {
	LogFn  domain.LogFn
	Client *http.Client
}

// Decorate wraps the input lambda with a Subscription decorator
func (s *Subscription) Decorate(l domain.Lambda) domain.Lambda {
	return func(ctx context.Context, in domain.SNSInput) error {
		logger := s.LogFn(ctx)
		t, ok := in["Type"]
		if !ok {
			logger.Info(logs.InvalidInput{Reason: "Missing notification Type"})
			return nil
		}
		val, ok := t.(string)
		if !ok {
			logger.Info(logs.InvalidInput{Reason: "Unable to parse notification Type"})
			return nil
		}
		if !strings.EqualFold(val, typeNotification) {
			return l(ctx, in)
		}
		var sub snsSubscription
		b, _ := json.Marshal(in)
		if e := json.Unmarshal(b, &sub); e != nil {
			logger.Info(logs.InvalidInput{Reason: e.Error()})
			return e
		}
		res, e := s.Client.Get(sub.SubscribeURL)
		if e != nil {
			logger.Error(logs.SubscriptionError{Reason: e.Error()})
			return e
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(res.Body)
			e := fmt.Errorf("Received unexpected error confirming the subscription [%d]: %s", res.StatusCode, b)
			logger.Error(logs.SubscriptionError{Reason: e.Error()})
			return e
		}
		return nil
	}
}

type snsSubscription struct {
	Type             string `json:"Type"`
	MessageID        string `json:"MessageId"`
	Token            string `json:"Token"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	SubscribeURL     string `json:"SubscribeURL"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
}
