package decorators

import (
	"context"
	"fmt"
	"strings"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
)

const (
	decoratorSubscription = "subscription"
)

// ChainConfig contains settings for the decorator chain.
type ChainConfig struct {
	Enabled      []string `description:"List of enabled decorators."`
	Subscription *SubscriptionConfig
}

// Name of the configuration root.
func (*ChainConfig) Name() string {
	return "decorator"
}

// ChainComponent is the top-level container for all decorator types.
type ChainComponent struct {
	Subscription *SubscriptionComponent
}

// NewChainComponent generates a ChainComponent.
func NewChainComponent() *ChainComponent {
	return &ChainComponent{
		Subscription: NewSubscriptionComponent(),
	}
}

// Settings generates the default configuration.
func (c *ChainComponent) Settings() *ChainConfig {
	return &ChainConfig{
		Enabled:      []string{decoratorSubscription},
		Subscription: c.Subscription.Settings(),
	}
}

// New generates a decorator chain.
func (c *ChainComponent) New(ctx context.Context, conf *ChainConfig) (Chain, error) {
	chain := make(Chain, 0, len(conf.Enabled))
	for _, enabled := range conf.Enabled {
		switch {
		case strings.EqualFold(enabled, decoratorSubscription):
			s, e := c.Subscription.New(ctx, conf.Subscription)
			if e != nil {
				return nil, e
			}
			chain = append(chain, s.Decorate)
		default:
			return nil, fmt.Errorf("unknown decorator type %s", enabled)
		}
	}
	return chain, nil
}

// Chain is a slice of decorators which can be chained together
type Chain []domain.Decorator

// Decorate applies the chain of decorators to the provided base lambda
func (c Chain) Decorate(base domain.Lambda) domain.Lambda {
	entry := base
	for i := len(c) - 1; i >= 0; i-- {
		entry = c[i](entry)
	}
	return entry
}
