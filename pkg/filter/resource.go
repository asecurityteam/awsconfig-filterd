package filter

import (
	"context"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
)

// ResourceTypeFiltererConfig defines the configuration options for a ResourceTypeFilterer.
type ResourceTypeFiltererConfig struct {
	ValidResourceTypes []string `description:"The set of AWS resource types that are filtered."`
}

// Name is used by the settings library to replace the default naming convention.
func (c *ResourceTypeFiltererConfig) Name() string {
	return "ResourceTypeFilter"
}

// ResourceTypeFiltererComponent satisfies the settings library Component API,
// and may be used by the settings.NewComponent function.
type ResourceTypeFiltererComponent struct{}

// Settings populates a set of default valid resource types for the ResourceTypeFilterer
// if none are provided via config.
func (*ResourceTypeFiltererComponent) Settings() *ResourceTypeFiltererConfig {
	return &ResourceTypeFiltererConfig{
		ValidResourceTypes: []string{
			config.ResourceTypeAwsEc2Instance,
			config.ResourceTypeAwsElasticLoadBalancingLoadBalancer,
			config.ResourceTypeAwsElasticLoadBalancingV2LoadBalancer,
		},
	}
}

// New constructs a ResourceTypeFilterer from a config.
func (*ResourceTypeFiltererComponent) New(_ context.Context, c *ResourceTypeFiltererConfig) (*ResourceTypeFilterer, error) {
	return &ResourceTypeFilterer{
		ValidResourceTypes: c.ValidResourceTypes,
	}, nil
}

// ResourceTypeFilterer filters AWS Config items on resource types.
type ResourceTypeFilterer struct {
	ValidResourceTypes []string
}

// FilterConfig returns true if a AWS Config item matches one of the valid resource types defined.
func (f *ResourceTypeFilterer) FilterConfig(c domain.ConfigurationItem) bool {
	for _, validResourceType := range f.ValidResourceTypes {
		if c.ResourceType == validResourceType {
			return true
		}
	}
	return false
}
