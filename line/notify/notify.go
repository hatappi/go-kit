package notify

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	endpoint = "https://notify-api.line.me/api/notify"
)

type NotifyResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (c *Client) Notify(ctx context.Context, message string) error {
	form := url.Values{}
	form.Add("message", message)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var notifyRes NotifyResponse
	if err = json.NewDecoder(res.Body).Decode(&notifyRes); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(notifyRes.Message)
	}

	return nil
}
