package decorators

import (
	"context"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestChainComponent(t *testing.T) {
	c := NewChainComponent()
	_, err := c.New(context.Background(), c.Settings())
	assert.NoError(t, err)
}

func TestChainComponentUnknown(t *testing.T) {
	c := NewChainComponent()
	s := c.Settings()
	s.Enabled = append(s.Enabled, "unknown")
	_, err := c.New(context.Background(), s)
	assert.Error(t, err)
}

func TestChain(t *testing.T) {
	numDecorators := 5
	callOrder := make([]int, numDecorators+1)
	c := make(Chain, 0)
	for i := 0; i < numDecorators; i++ {
		cpy := i
		c = append(c, func(l domain.Lambda) domain.Lambda {
			return func(ctx context.Context, in domain.SNSInput) error {
				callOrder[cpy] = cpy
				return l(ctx, in)
			}
		})
	}
	base := func(_ context.Context, _ domain.SNSInput) error {
		callOrder[numDecorators] = numDecorators
		return nil
	}
	_ = c.Decorate(base)(context.Background(), make(map[string]interface{}))
	for i, v := range callOrder {
		assert.Equal(t, i, v) // assert that the call order is the same as the decorator chain order
	}
}
