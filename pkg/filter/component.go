package filter

import (
	"context"
	"fmt"
	"strings"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
)

// FilterConfig contains settings for the FilterComponent.
type FilterConfig struct {
	Enabled      []string `description:"List of enabled filters."`
	ResourceType *ResourceTypeConfig
}

// Name of the configuration root.
func (*FilterConfig) Name() string {
	return "filter"
}

// FilterComponent is the top-level container for all filter types.
type FilterComponent struct {
	ResourceType *ResourceTypeComponent
}

// NewFilterComponent generates a FilterComponent.
func NewFilterComponent() *FilterComponent {
	return &FilterComponent{
		ResourceType: NewResourceTypeComponent(),
	}
}

// Settings generates the default configuration.
func (c *FilterComponent) Settings() *FilterConfig {
	return &FilterConfig{
		Enabled:      []string{"resourcetype"},
		ResourceType: c.ResourceType.Settings(),
	}
}

// New generates a filter.
func (c *FilterComponent) New(ctx context.Context, conf *FilterConfig) (domain.ConfigFilterer, error) {
	f := make(AnyMatch, 0, len(conf.Enabled))
	for _, enabled := range conf.Enabled {
		switch {
		case strings.EqualFold(enabled, "resourcetype"):
			fi, err := c.ResourceType.New(ctx, conf.ResourceType)
			if err != nil {
				return nil, err
			}
			f = append(f, fi)
		default:
			return nil, fmt.Errorf("unknown filter type %s", enabled)
		}
	}
	return f, nil
}
