package filter

import (
	"context"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	config "github.com/aws/aws-sdk-go/service/configservice"
	"github.com/stretchr/testify/require"
)

func TestFilterComponent(t *testing.T) {
	cmp := NewFilterComponent()
	f, err := cmp.New(context.Background(), cmp.Settings())
	require.Nil(t, err)

	c := domain.ConfigurationItem{ResourceType: "not match"}
	require.False(t, f.FilterConfig(c))
	c = domain.ConfigurationItem{ResourceType: config.ResourceTypeAwsEc2Instance}
	require.True(t, f.FilterConfig(c))
}
