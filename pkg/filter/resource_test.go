package filter

import (
	"fmt"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
	"github.com/stretchr/testify/require"
)

func TestResourceTypeFilter(t *testing.T) {
	tc := []struct {
		name         string
		in           domain.ConfigurationItem
		expectedBool bool
		expectedErr  error
	}{
		{
			"EC2",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsEc2Instance},
			true,
			nil,
		},
		{
			"ELB",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsElasticLoadBalancingLoadBalancer},
			true,
			nil,
		},
		{
			"ALB",
			domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsElasticLoadBalancingV2LoadBalancer},
			true,
			nil,
		},
		{
			"unknown resource",
			domain.ConfigurationItem{ResourceType: "AWS::EC3::Instance"},
			false,
			fmt.Errorf("no valid resource type found"),
		},
		{
			"empty resource",
			domain.ConfigurationItem{ResourceType: ""},
			false,
			fmt.Errorf("no valid resource type found"),
		},
		{
			"substring match",
			domain.ConfigurationItem{ResourceType: fmt.Sprintf("%s%s", config.ResourceTypeAwsEc2Instance, config.ResourceTypeAwsEc2Instance)},
			false,
			fmt.Errorf("no valid resource type found"),
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			filter := &ResourceTypeFilter{}
			actualBool, actualErr := filter.Filter(tt.in)
			require.Equal(t, tt.expectedBool, actualBool)
			require.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
