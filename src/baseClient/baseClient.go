package baseClient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	DefaultUrl = "https://www.toggl.com/api/v8"
)

type TogglClient struct {
	client *http.Client
	// cookie *http.Cookie
}

func NewClient() (*TogglClient, error) {
	return &TogglClient{
		client: http.DefaultClient,
	}, nil
}

func request(c *TogglClient, method, endpoint string, b interface{}) (_ *json.RawMessage, err error) {

	var body []byte

	if b != nil {
		if body, err = json.Marshal(b); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// req.AddCookie(c.cookie)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth("apiKey", "api_token")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	js, err := io.ReadAll(resp.Body)

	var raw json.RawMessage
	if json.Unmarshal(js, &raw) != nil {
		return nil, err
	}
	return &raw, err
}

func (c *TogglClient) GetRequest(endpoint string) (*json.RawMessage, error) {
	return request(c, "GET", endpoint, nil)
}
