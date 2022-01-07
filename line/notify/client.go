package notify

import "net/http"

type Client struct {
	HTTPClient *http.Client

	accessToken string
	endpoint    string
}

func New(accessToken string) *Client {
	return &Client{
		HTTPClient:  http.DefaultClient,
		accessToken: accessToken,
	}
}
