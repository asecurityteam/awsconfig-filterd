package filter

import (
	"fmt"
	"sync"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
)

var (
	// defaultValidResourceType are the set of AWS resource types
	// that are filtered by default
	defaultValidResourceTypes = []string{
		config.ResourceTypeAwsEc2Instance,
		config.ResourceTypeAwsElasticLoadBalancingLoadBalancer,
		config.ResourceTypeAwsElasticLoadBalancingV2LoadBalancer,
	}
)

// ResourceTypeFilter filters AWS Config items on resource types.
type ResourceTypeFilter struct {
	ValidResourceTypes []string
	once               sync.Once
}

// init adds a default set of valid resource types if none are provided.
func (f *ResourceTypeFilter) init() {
	if len(f.ValidResourceTypes) == 0 {
		f.ValidResourceTypes = defaultValidResourceTypes
	}
}

// Filter returns true if a AWS Config item matches one of the valid resource types defined.
func (f *ResourceTypeFilter) Filter(c domain.ConfigurationItem) (bool, error) {
	f.once.Do(f.init)

	for _, validResourceType := range f.ValidResourceTypes {
		if c.ResourceType == validResourceType {
			return true, nil
		}
	}
	return false, fmt.Errorf("no valid resource type found")
}
