package domain

import (
	"context"
)

// SNSInput represents the incoming SNS payload to our lambda handler
type SNSInput map[string]interface{}

// Lambda is the function signature of the lambda handler in this project
type Lambda func(context.Context, SNSInput) error

// Decorator returns a new Lambda which decorates in the input function
type Decorator func(Lambda) Lambda
