package filter

import (
	"context"
	"fmt"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
	"github.com/stretchr/testify/require"
)

func TestResourceType(t *testing.T) {
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

	cmp := NewResourceTypeComponent()
	f, err := cmp.New(context.Background(), cmp.Settings())
	require.Nil(t, err)

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := f.FilterConfig(tt.in)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestResourceTypeFiltererNoValidTypes(t *testing.T) {
	f := &ResourceType{}
	out := f.FilterConfig(domain.ConfigurationItem{
		ResourceType: config.ResourceTypeAwsEc2Instance,
	})
	require.False(t, out)
}

func TestName(t *testing.T) {
	conf := ResourceTypeConfig{}
	require.Equal(t, "resourcetype", conf.Name())
}
