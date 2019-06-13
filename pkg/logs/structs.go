package logs

// InvalidInput is logged when the input provided is not valid
type InvalidInput struct {
	Reason  string `logevent:"reason"`
	Message string `logevent:"message,default=invalid-input"`
}

// SubscriptionError is logged when there is an error confirming the SNS subscription
type SubscriptionError struct {
	Reason  string `logevent:"reason"`
	Message string `logevent:"message,default=subscription-error"`
}
