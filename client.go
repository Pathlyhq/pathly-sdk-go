// Package pathly is the official Go client for the Pathly public `/v1` API
// (https://pathlyhq.com — docs: https://pathlyhq.com/en/developers).
//
// Hand-written rather than generated: creates stay idempotent, Retry-After is
// honoured, and a missing resource (404) stays distinct from a transport failure.
package pathly

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Version is the SDK semver (also used in the default User-Agent).
	Version = "0.1.0"

	// DefaultBaseURL is the production API.
	DefaultBaseURL = "https://api.pathlyhq.com"

	defaultUserAgent = "pathly-sdk-go/" + Version

	maxAttempts  = 4
	maxRetryWait = 90 * time.Second
	pageSize     = "200"
	maxPages     = 200
)

// Client is ready to use and safe for concurrent use. The token is never logged.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	userAgent  string
	sleep      func(time.Duration)
}

// Option configures the client at construction time.
type Option func(*Client)

// WithHTTPClient forces an HTTP client, for tests or a corporate proxy.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.httpClient = h }
}

// WithUserAgent identifies the SDK version in the API logs.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithSleep replaces the wait between two attempts (tests inject a no-op).
func WithSleep(f func(time.Duration)) Option {
	return func(c *Client) { c.sleep = f }
}

// WithBaseURL overrides the API base URL.
func WithBaseURL(base string) Option {
	return func(c *Client) {
		if strings.TrimSpace(base) != "" {
			c.baseURL = strings.TrimRight(strings.TrimSpace(base), "/")
		}
	}
}

// WithToken sets the API token explicitly.
func WithToken(token string) Option {
	return func(c *Client) { c.token = strings.TrimSpace(token) }
}

// New builds a client. Token defaults to PATHLY_API_TOKEN; base URL to
// PATHLY_API_URL or DefaultBaseURL.
func New(opts ...Option) (*Client, error) {
	token := strings.TrimSpace(os.Getenv("PATHLY_API_TOKEN"))
	base := strings.TrimSpace(os.Getenv("PATHLY_API_URL"))
	if base == "" {
		base = DefaultBaseURL
	}
	c := &Client{
		baseURL:    strings.TrimRight(base, "/"),
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  defaultUserAgent,
		sleep:      time.Sleep,
	}
	for _, o := range opts {
		o(c)
	}
	if c.token == "" {
		return nil, fmt.Errorf("PATHLY_API_TOKEN is required. Export it, or pass WithToken")
	}
	return c, nil
}

type request struct {
	method         string
	path           string
	body           any
	idempotencyKey string
	out            any
}

// do sends the request and retries transient failures (429, 5xx, transport).
func (c *Client) do(ctx context.Context, r request) error {
	var payload []byte
	if r.body != nil {
		encoded, err := json.Marshal(r.body)
		if err != nil {
			return fmt.Errorf("encoding the body of %s: %w", r.path, err)
		}
		payload = encoded
	}

	for attempt := 1; attempt < maxAttempts; attempt++ {
		retryIn, err := c.attempt(ctx, r, payload, attempt, false)
		if retryIn <= 0 {
			return err
		}
		c.sleep(retryIn)
	}
	_, err := c.attempt(ctx, r, payload, maxAttempts, true)
	return err
}

func (c *Client) attempt(
	ctx context.Context,
	r request,
	payload []byte,
	attempt int,
	last bool,
) (time.Duration, error) {
	var reader io.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, c.baseURL+r.path, reader)
	if err != nil {
		return 0, fmt.Errorf("building request %s %s: %w", r.method, r.path, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if r.idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", r.idempotencyKey)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		wrapped := fmt.Errorf("calling %s %s: %w", r.method, r.path, err)
		if last {
			return 0, wrapped
		}
		return backoff(attempt), wrapped
	}

	body, readErr := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if readErr != nil {
		wrapped := fmt.Errorf("reading the response of %s %s: %w", r.method, r.path, readErr)
		if last {
			return 0, wrapped
		}
		return backoff(attempt), wrapped
	}

	if res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500 {
		failure := apiError(res.StatusCode, r.path, body)
		if last {
			return 0, failure
		}
		return waitFor(res, attempt), failure
	}

	if res.StatusCode >= 400 {
		return 0, apiError(res.StatusCode, r.path, body)
	}

	if r.out == nil || len(body) == 0 {
		return 0, nil
	}
	if err := json.Unmarshal(body, r.out); err != nil {
		return 0, fmt.Errorf("unreadable response from %s %s: %w", r.method, r.path, err)
	}
	return 0, nil
}

func waitFor(res *http.Response, attempt int) time.Duration {
	if raw := res.Header.Get("Retry-After"); raw != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && secs > 0 {
			wait := time.Duration(secs) * time.Second
			if wait > maxRetryWait {
				return maxRetryWait
			}
			return wait
		}
	}
	return backoff(attempt)
}

func backoff(attempt int) time.Duration {
	return time.Duration(attempt) * 500 * time.Millisecond
}

// Ping checks that the token is accepted.
//
// A 403 is not a failure: it proves the key is valid and only signals that it
// does not carry org:read. Treating that as an error would force every
// scenario-only key to request an organization scope.
func (c *Client) Ping(ctx context.Context) error {
	err := c.do(ctx, request{method: http.MethodGet, path: "/v1/usage"})
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusForbidden {
		return nil
	}
	return err
}

func listPaged[T any](ctx context.Context, c *Client, base string) ([]T, error) {
	var all []T
	cursor := ""
	for page := 0; page < maxPages; page++ {
		var out struct {
			Items      []T     `json:"items"`
			NextCursor *string `json:"nextCursor"`
		}
		query := url.Values{"limit": []string{pageSize}}
		if cursor != "" {
			query.Set("cursor", cursor)
		}
		path := base + "?" + query.Encode()
		if err := c.do(ctx, request{method: http.MethodGet, path: path, out: &out}); err != nil {
			return nil, err
		}
		all = append(all, out.Items...)
		if out.NextCursor == nil || *out.NextCursor == "" {
			return all, nil
		}
		cursor = *out.NextCursor
	}
	return nil, fmt.Errorf("%s pagination: too many pages, stopping on purpose", base)
}
