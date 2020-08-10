package filter

import (
	"context"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/aws/aws-sdk-go/service/configservice"
)

// ResourceTypeConfig contains settings for the AWS resource type filter.
type ResourceTypeConfig struct {
	Allowed []string `description:"List of AWS resource types allowed to pass through."`
}

// Name is used by the settings library to replace the default naming convention.
func (c *ResourceTypeConfig) Name() string {
	return "resourcetype"
}

// ResourceTypeComponent satisfies the settings library Component API,
// and may be used by the settings.NewComponent function.
type ResourceTypeComponent struct{}

// NewResourceTypeComponent constructs a ResourceTypeComponent.
func NewResourceTypeComponent() *ResourceTypeComponent {
	return &ResourceTypeComponent{}
}

// Settings populates a set of default valid resource types for the ResourceType
// if none are provided via config.
func (*ResourceTypeComponent) Settings() *ResourceTypeConfig {
	return &ResourceTypeConfig{
		Allowed: []string{
			configservice.ResourceTypeAwsEc2Instance,
			configservice.ResourceTypeAwsElasticLoadBalancingLoadBalancer,
			configservice.ResourceTypeAwsElasticLoadBalancingV2LoadBalancer,
			configservice.ResourceTypeAwsEc2NetworkInterface,
		},
	}
}

// New constructs a ResourceType from a config.
func (*ResourceTypeComponent) New(_ context.Context, c *ResourceTypeConfig) (*ResourceType, error) {
	return &ResourceType{
		Allowed: c.Allowed,
	}, nil
}

// ResourceType is a filter that only allows a known set of resource types
// to pass through.
type ResourceType struct {
	Allowed []string
}

// FilterConfig returns true if a AWS Config item matches one of the allowed
// resource types defined.
func (f *ResourceType) FilterConfig(c domain.ConfigurationItem) bool {
	for _, allowed := range f.Allowed {
		if c.ResourceType == allowed {
			return true
		}
	}
	return false
}
