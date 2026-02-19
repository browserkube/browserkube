package reportvideo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type RecorderClient interface {
	Stop(ctx context.Context) error
}

type Config struct {
	BaseURL string
	timeout time.Duration
	client  *http.Client
}

type Option func(*Config)

func WithTimeout(t time.Duration) Option {
	return func(c *Config) {
		c.timeout = t
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Config) {
		c.client = client
	}
}

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(cfg *Config, opts ...Option) *Client {
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.timeout == 0 {
		cfg.timeout = time.Second * 5
	}

	if cfg.client == nil {
		cfg.client = &http.Client{
			Timeout: cfg.timeout,
		}
	}

	return &Client{
		client:  cfg.client,
		baseURL: cfg.BaseURL,
	}
}

func (c *Client) Stop(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/recorder/stop", nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to stop recording: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to stop recording: status code: %v", resp.StatusCode)
	}

	return nil
}
