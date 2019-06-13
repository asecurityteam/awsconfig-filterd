package decorators

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func fakeLambda(_ context.Context, _ domain.SNSInput) error {
	return nil
}

func TestComponent(t *testing.T) {
	c := NewSubscriptionComponent()
	_, err := c.New(context.Background(), c.Settings())
	assert.NoError(t, err)
}

func TestSubscription(t *testing.T) {
	tc := []struct {
		Name         string
		Input        map[string]interface{}
		ResponseCode int
		RTError      error
		ExpectError  bool
	}{
		{
			Name: "subscribe success",
			Input: map[string]interface{}{
				"Type":         "SubscriptionConfirmation",
				"SubscribeURL": "http://likeandsubscribe",
			},
			ResponseCode: 200,
		},
		{
			Name: "no type",
			Input: map[string]interface{}{
				"SubscribeURL": "http://likeandsubscribe",
			},
		},
		{
			Name: "bad type",
			Input: map[string]interface{}{
				"Type":         100,
				"SubscribeURL": "http://likeandsubscribe",
			},
		},
		{
			Name: "not a subscription",
			Input: map[string]interface{}{
				"Type": "Notification",
			},
		},
		{
			Name: "malformed notification",
			Input: map[string]interface{}{
				"Type":         "SubscriptionConfirmation",
				"SubscribeURL": 100,
			},
			ExpectError: true,
		},
		{
			Name: "roud trip error",
			Input: map[string]interface{}{
				"Type":         "SubscriptionConfirmation",
				"SubscribeURL": "http://likeandsubscribe",
			},
			RTError:     errors.New(""),
			ExpectError: true,
		},
		{
			Name: "roud trip error",
			Input: map[string]interface{}{
				"Type":         "SubscriptionConfirmation",
				"SubscribeURL": "http://likeandsubscribe",
			},
			RTError:     errors.New(""),
			ExpectError: true,
		},
		{
			Name: "subscribe error",
			Input: map[string]interface{}{
				"Type":         "SubscriptionConfirmation",
				"SubscribeURL": "http://likeandsubscribe",
			},
			ExpectError:  true,
			ResponseCode: 500,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			res := &http.Response{
				StatusCode: tt.ResponseCode,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}
			mockRT := NewMockRoundTripper(ctrl)
			mockRT.EXPECT().RoundTrip(gomock.Any()).Return(res, tt.RTError).AnyTimes()

			s := &Subscription{
				LogFn: testLogFn,
				Client: &http.Client{
					Transport: mockRT,
				},
			}

			e := s.Decorate(fakeLambda)(context.Background(), tt.Input)
			assert.Equal(t, tt.ExpectError, e != nil)
		})
	}

}
