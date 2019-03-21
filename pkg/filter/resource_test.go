package filter

import (
	"context"
	"fmt"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
	"github.com/stretchr/testify/require"
)

func TestResourceTypeFilterer(t *testing.T) {
	tc := []struct {
		name     string
		in       domain.ConfigurationItem
		expected bool
	}{
		{
			"EC2",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsEc2Instance},
			true,
		},
		{
			"ELB",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsElasticLoadBalancingLoadBalancer},
			true,
		},
		{
			"ALB",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsElasticLoadBalancingV2LoadBalancer},
			true,
		},
		{
			"unknown resource",
			domain.ConfigurationItem{ResourceType: "AWS::EC3::Instance"},
			false,
		},
		{
			"empty resource",
			domain.ConfigurationItem{ResourceType: ""},
			false,
		},
		{
			"substring match",
			domain.ConfigurationItem{ResourceType: fmt.Sprintf("%s%s", config.ResourceTypeAwsEc2Instance, config.ResourceTypeAwsEc2Instance)},
			false,
		},
	}

	resourceTypeFiltererComponent := &ResourceTypeFiltererComponent{}
	resourceTypeFilterer, err := resourceTypeFiltererComponent.New(context.Background(), resourceTypeFiltererComponent.Settings())
	require.Nil(t, err)

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {

			filter := resourceTypeFilterer
			actual := filter.FilterConfig(tt.in)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestResourceTypeFiltererNoValidTypes(t *testing.T) {
	filter := &ResourceTypeFilterer{}
	out := filter.FilterConfig(domain.ConfigurationItem{
		ResourceType: config.ResourceTypeAwsEc2Instance})
	require.Equal(t, false, out)
}

func TestName(t *testing.T) {
	resourceTypeFiltererConfig := ResourceTypeFiltererConfig{}
	require.Equal(t, "ResourceTypeFilter", resourceTypeFiltererConfig.Name())
}
