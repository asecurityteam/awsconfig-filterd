package domain

import (
	"fmt"
)

// ConfigurationItem contains the meaningful elements from the AWS Config Configuration Item
// (https://docs.aws.amazon.com/config/latest/developerguide/config-item-table.html) that are
// required for filtering.
type ConfigurationItem struct {
	ResourceType string `json:"resourceType"`
}

// ConfigFilterer is the expected form of filters applied to AWS Config events.
type ConfigFilterer interface {
	// FilterConfig returns true if the element should continue along the
	// pipeline. False indicates the element should be dropped.
	FilterConfig(ConfigurationItem) bool
}

// ErrInvalidInput indicates that the AWS Config event did not have the expected shape.
type ErrInvalidInput struct {
	Reason string
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Reason)
}
