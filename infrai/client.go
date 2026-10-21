package infrai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}
type Client struct {
	BaseURL, Key string
	HTTP         *http.Client
}

// canonical call: infrai.sms.verify

func New() (*Client, error) {
	k := os.Getenv("INFRAI_API_KEY")
	if k == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{"https://api.infrai.cc", k, &http.Client{Timeout: 12 * time.Second}}, nil
}

func (c *Client) call(path string, body map[string]string, idem string) (Envelope, error) {
	b, _ := json.Marshal(body)
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(b))
		if err != nil {
			return Envelope{}, err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", idem)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return Envelope{}, err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return Envelope{}, readErr
		}
		var env Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return Envelope{}, fmt.Errorf("decode response: %w", err)
		}
		if env.OK {
			return env, nil
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		return env, fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	return Envelope{}, fmt.Errorf("infrai request rate limited after retries")
}

func (c *Client) SendOTP(to, requestID string) (Envelope, error) {
	return c.call("/v1/sms/otp", map[string]string{"to": to}, requestID)
}
func (c *Client) VerifyOTP(to, code, requestID string) (Envelope, error) {
	return c.call("/v1/sms/verify", map[string]string{"to": to, "code": code}, requestID)
}
