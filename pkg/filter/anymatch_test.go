package filter

import (
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/stretchr/testify/require"
)

type falseFilter struct{}

func (*falseFilter) FilterConfig(_ domain.ConfigurationItem) bool {
	return false
}

type trueFilter struct{}

func (*trueFilter) FilterConfig(_ domain.ConfigurationItem) bool {
	return true
}

func TestAnyMatch(t *testing.T) {
	f := make(AnyMatch, 0)
	c := domain.ConfigurationItem{ResourceType: ""}
	require.False(t, f.FilterConfig(c))

	f = append(f, &falseFilter{})
	f = append(f, &falseFilter{})
	f = append(f, &falseFilter{})
	f = append(f, &trueFilter{})
	require.True(t, f.FilterConfig(c))
}
