package baseClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultAuthPassword = "api_token"
	DefaultUrl          = "https://api.track.toggl.com/api/v9"
	SessionCookieName   = "__Host-timer-session"
)

type TogglClient struct {
	client   *http.Client
	Url      string
	password string
	cookie   *http.Cookie
}

func NewClient(apiKey string) (*TogglClient, error) {
	client := &TogglClient{
		client:   http.DefaultClient,
		Url:      DefaultUrl,
		password: DefaultAuthPassword,
	}

	if len(apiKey) < 1 {
		fmt.Printf("%s\n", "valid token required")
		return nil, errors.New("token required")
	}

	if _, err := client.authenticate(apiKey); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *TogglClient) authenticate(apiKey string) (_ []byte, err error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.Url, "me/sessions"), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(apiKey, DefaultAuthPassword)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer func() {
			err = resp.Body.Close()
		}()
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}

	hasSessionCookie := false
	for _, value := range resp.Cookies() {
		if value.Name == SessionCookieName {
			fmt.Printf("Cookie: %s\n", value)
			fmt.Printf("Setting Cookie\n")
			c.cookie = value
			hasSessionCookie = true
			break
		}
	}

	if !hasSessionCookie {
		return nil, fmt.Errorf("auth cookie %q not found on authentication response", SessionCookieName)
	}

	return nil, err
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
	req.AddCookie(c.cookie)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
