package notify

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type testRoundTripper struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (trt testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return trt.roundTripFunc(req)
}

func TestNotify(t *testing.T) {
	testCases := []struct {
		name          string
		roundTripFunc func(*http.Request) (*http.Response, error)
		wantErr       bool
	}{
		{
			name: "success",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"status":200, "message":"ok"}`)),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "api returns invalid json",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`test`)),
				}, nil
			},
			wantErr: true,
		},
		{
			name: "api returns 401",
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"status":401,"message":"Invalid access token"}`)),
				}, nil
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client := &Client{
				HTTPClient: &http.Client{
					Transport: testRoundTripper{
						roundTripFunc: tc.roundTripFunc,
					},
				},
			}

			ctx := context.Background()
			err := client.Notify(ctx, "test")
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}
		})
	}
}
