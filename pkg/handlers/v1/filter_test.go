package v1

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	domain "github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var eventTimestamp = "2019-02-22T20:43:11.479Z"

func dataFileToString(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("failed to read file '%s': %s", filename, err)
	}
	return string(data)
}
func TestHandle(t *testing.T) {
	validEvent := dataFileToString(t, "config.valid.json")
	invalidResourceType := dataFileToString(t, "config.invalid-resourceType.json")
	noResourceType := dataFileToString(t, "config.no-resourceType.json")
	invalidMessageType := dataFileToString(t, "config.invalid-messageType.json")

	tc := []struct {
		name         string
		in           string
		t            string
		expectedOut  string
		expectedErr  error
		filterCalled bool
		filterOK     bool
	}{
		{
			name:         "success",
			in:           validEvent,
			t:            "Notification",
			expectedOut:  validEvent,
			expectedErr:  nil,
			filterCalled: true,
			filterOK:     true,
		},
		{
			name:         "no type",
			in:           validEvent,
			expectedOut:  validEvent,
			expectedErr:  nil,
			filterCalled: false,
			filterOK:     false,
		},
		{
			name:         "invalid resource type",
			in:           invalidResourceType,
			t:            "Notification",
			expectedOut:  "",
			expectedErr:  nil,
			filterCalled: true,
			filterOK:     false,
		},
		{
			name:         "no resource type",
			in:           noResourceType,
			t:            "Notification",
			expectedOut:  "",
			expectedErr:  domain.ErrInvalidInput{Reason: "empty resource type"},
			filterCalled: false,
			filterOK:     false,
		},
		{
			name:         "invalid message type",
			in:           invalidMessageType,
			t:            "Notification",
			expectedOut:  "",
			expectedErr:  nil,
			filterCalled: false,
			filterOK:     false,
		},
		{
			name:         "no message",
			in:           "",
			t:            "Notification",
			expectedOut:  "",
			expectedErr:  nil,
			filterCalled: false,
			filterOK:     false,
		},
		{
			name:         "cannot unmarshal ConfigEvent",
			in:           "0",
			t:            "Notification",
			expectedOut:  "",
			expectedErr:  &json.UnmarshalTypeError{Value: "number", Offset: 1, Type: reflect.TypeOf(ConfigEvent{})},
			filterCalled: false,
			filterOK:     false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bites, _ := json.Marshal(configNotification{
				Type:      tt.t,
				Message:   tt.in,
				Timestamp: eventTimestamp,
			})
			mockFilterer := NewMockConfigFilterer(ctrl)
			if tt.filterCalled {
				mockFilterer.EXPECT().FilterConfig(gomock.Any()).Return(tt.filterOK)
			}
			mockProducer := NewMockProducer(ctrl)
			if tt.filterOK {
				mockProducer.EXPECT().Produce(gomock.Any(), gomock.Any()).Do(
					func(_ context.Context, event interface{}) {
						require.Equal(t, tt.in, event.(configNotification).Message)
					},
				).Return(nil, nil)
			}

			configFilterHandler := &ConfigFilter{
				LogFn:          testLogFn,
				StatFn:         testStatFn,
				ConfigFilterer: mockFilterer,
				Producer:       mockProducer,
			}
			var input domain.SNSInput
			_ = json.Unmarshal(bites, &input)
			actualErr := configFilterHandler.Handle(context.Background(), input)
			require.IsType(t, tt.expectedErr, actualErr)
		})
	}
}
