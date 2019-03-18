package domain

import (
	"fmt"
)

// ConfigurationItem contains the meaningful elements from the AWS Config Configuration Item
// (https://docs.aws.amazon.com/config/latest/developerguide/config-item-table.html) that are
// required for filtering.
type ConfigurationItem struct {
	ResourceType string `json:"ResourceType"`
}

// ConfigFilter is the expected form of filters applied to AWS Config events
type ConfigFilter interface {
	Filter(ConfigurationItem) (bool, error)
}

// ErrEventDiscarded indicates that the event did not match any filters
type ErrEventDiscarded struct {
	Reason string
}

func (e ErrEventDiscarded) Error() string {
	return fmt.Sprintf("event was discarded: %s", e.Reason)
}
