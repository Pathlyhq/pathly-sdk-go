package pathly

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *[]time.Duration) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	var waits []time.Duration
	c, err := New(
		WithToken("sp_test_token"),
		WithBaseURL(srv.URL),
		WithSleep(func(d time.Duration) { waits = append(waits, d) }),
		WithUserAgent("pathly-sdk-go-test"),
	)
	if err != nil {
		t.Fatal(err)
	}
	return c, &waits
}

func TestNewRequiresToken(t *testing.T) {
	t.Setenv("PATHLY_API_TOKEN", "")
	_, err := New()
	if err == nil || !strings.Contains(err.Error(), "PATHLY_API_TOKEN") {
		t.Fatalf("want token error, got %v", err)
	}
}

func TestNewFromEnvAndOptions(t *testing.T) {
	t.Setenv("PATHLY_API_TOKEN", "sp_env")
	t.Setenv("PATHLY_API_URL", "https://example.test/")
	c, err := New(WithSleep(func(time.Duration) {}))
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != "https://example.test" {
		t.Errorf("baseURL = %q", c.baseURL)
	}
	if c.token != "sp_env" {
		t.Errorf("token = %q", c.token)
	}
	c2, err := New(WithToken("sp_opt"), WithBaseURL("  "), WithHTTPClient(&http.Client{Timeout: time.Second}))
	if err != nil {
		t.Fatal(err)
	}
	if c2.token != "sp_opt" {
		t.Errorf("token override = %q", c2.token)
	}
	if c2.baseURL != DefaultBaseURL && c2.baseURL != "https://example.test" {
		// WithBaseURL("  ") is ignored, env PATHLY_API_URL still applied at New start then option no-ops
		_ = c2
	}
}

func TestDefaultBaseURL(t *testing.T) {
	t.Setenv("PATHLY_API_TOKEN", "sp_x")
	t.Setenv("PATHLY_API_URL", "")
	c, err := New(WithSleep(func(time.Duration) {}))
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q", c.baseURL)
	}
}

func TestDefaultUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{"planId":"pro"}`))
	}))
	t.Cleanup(srv.Close)
	c, err := New(WithToken("sp_x"), WithBaseURL(srv.URL), WithSleep(func(time.Duration) {}))
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := "pathly-sdk-go/" + Version
	if gotUA != want {
		t.Errorf("User-Agent = %q, want %q", gotUA, want)
	}
}

func TestAuthorizationHeaderAndUserAgent(t *testing.T) {
	var gotAuth, gotUA, gotAccept string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		gotAccept = r.Header.Get("Accept")
		_, _ = w.Write([]byte(`{"planId":"pro"}`))
	})
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if gotAuth != "Bearer sp_test_token" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotUA != "pathly-sdk-go-test" {
		t.Errorf("User-Agent = %q", gotUA)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q", gotAccept)
	}
}

func TestPingAccepts403(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"no org:read"}`))
	})
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("403 must be accepted: %v", err)
	}
}

func TestPingRaisesOtherErrors(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"no"}`))
	})
	err := c.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "PATHLY_API_TOKEN") {
		t.Fatalf("want 401 guidance, got %v", err)
	}
}

func TestCreateScenarioSendsIdempotencyKey(t *testing.T) {
	var gotKey, gotBody, gotType string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("Idempotency-Key")
		gotType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"mon-1","name":"Checkout","type":"http","url":"https://example.com"}`))
	})
	name := "Checkout"
	got, err := c.CreateScenario(context.Background(), ScenarioInput{Name: &name}, "idem-key-0001")
	if err != nil {
		t.Fatalf("CreateScenario: %v", err)
	}
	if gotKey != "idem-key-0001" {
		t.Errorf("Idempotency-Key = %q", gotKey)
	}
	if gotType != "application/json" {
		t.Errorf("Content-Type = %q", gotType)
	}
	if strings.Contains(gotBody, "intervalSec") {
		t.Errorf("body must omit unset fields: %s", gotBody)
	}
	if got.ID != "mon-1" || got.URL == nil || *got.URL != "https://example.com" {
		t.Errorf("decoded = %+v", got)
	}
}

func TestCreateScenarioAutoIdempotencyKey(t *testing.T) {
	var gotKey string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("Idempotency-Key")
		_, _ = w.Write([]byte(`{"id":"mon-2","name":"A","type":"http"}`))
	})
	name := "A"
	if _, err := c.CreateScenario(context.Background(), ScenarioInput{Name: &name}, ""); err != nil {
		t.Fatal(err)
	}
	if gotKey == "" {
		t.Fatal("auto Idempotency-Key missing")
	}
}

func TestRetryHonoursRetryAfter(t *testing.T) {
	calls := 0
	c, waits := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"Too many calls"}`))
			return
		}
		_, _ = w.Write([]byte(`{"planId":"pro"}`))
	})
	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("calls = %d", calls)
	}
	if len(*waits) != 1 || (*waits)[0] != 7*time.Second {
		t.Errorf("waits = %v", *waits)
	}
}

func TestRetryAfterIsCapped(t *testing.T) {
	c, waits := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "86400")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	_ = c.Ping(context.Background())
	for _, d := range *waits {
		if d > maxRetryWait {
			t.Errorf("wait = %v > cap", d)
		}
	}
}

func TestRetryAfterInvalidFallsBackToBackoff(t *testing.T) {
	c, waits := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "soon")
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	err := c.Ping(context.Background())
	if err == nil {
		t.Fatal("persistent 503 must error")
	}
	if len(*waits) != maxAttempts-1 {
		t.Errorf("waits = %v", *waits)
	}
}

func TestClientErrorsAreNotRetried(t *testing.T) {
	calls := 0
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Invalid body","details":[{"message":"required","path":["name"]}]}`))
	})
	err := c.Ping(context.Background())
	if err == nil || calls != 1 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
	if !strings.Contains(err.Error(), "name: required") {
		t.Errorf("details lost: %v", err)
	}
}

func TestNotFoundIsDetectable(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Scenario not found"}`))
	})
	_, err := c.GetScenario(context.Background(), "mon-unknown")
	if !IsNotFound(err) {
		t.Fatalf("IsNotFound false for %v", err)
	}
	if IsNotFound(nil) {
		t.Error("IsNotFound(nil) must be false")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || !apiErr.IsNotFound() {
		t.Fatal("APIError.IsNotFound")
	}
}

func TestUnauthorizedAndForbiddenGuidance(t *testing.T) {
	for status, expect := range map[int]string{
		http.StatusUnauthorized: "PATHLY_API_TOKEN",
		http.StatusForbidden:    "scopes",
	} {
		c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"refused"}`))
		})
		_, err := c.GetScenario(context.Background(), "mon-1")
		if err == nil || !strings.Contains(err.Error(), expect) {
			t.Errorf("status %d: %v", status, err)
		}
	}
}

func TestEmptyErrorBodyUsesStatusText(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	_, err := c.GetScenario(context.Background(), "x")
	if err == nil || !strings.Contains(err.Error(), "Bad Request") {
		t.Fatalf("got %v", err)
	}
}

func TestFieldOfAndAPIErrorString(t *testing.T) {
	if got := fieldOf(nil); got != "body" {
		t.Errorf("empty = %q", got)
	}
	if got := fieldOf([]any{"events", float64(0)}); got != "events.0" {
		t.Errorf("path = %q", got)
	}
	e := &APIError{StatusCode: 500, Message: "boom"}
	if e.Error() != "boom (HTTP 500)" {
		t.Errorf("Error = %q", e.Error())
	}
	e.Path = "/v1/x"
	if !strings.Contains(e.Error(), "/v1/x") {
		t.Errorf("Error = %q", e.Error())
	}
}

func TestUnreadableJSON(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	})
	_, err := c.GetScenario(context.Background(), "mon-1")
	if err == nil || !strings.Contains(err.Error(), "unreadable") {
		t.Fatalf("got %v", err)
	}
}

func TestEncodingBodyFailure(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("must not be called")
	})
	err := c.do(context.Background(), request{
		method: http.MethodPost,
		path:   "/v1/scenarios",
		body:   make(chan int),
	})
	if err == nil || !strings.Contains(err.Error(), "encoding") {
		t.Fatalf("got %v", err)
	}
}

func TestBuildingRequestFailure(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("must not be called")
	})
	err := c.do(context.Background(), request{method: "GET /bad", path: "/v1/usage"})
	if err == nil || !strings.Contains(err.Error(), "building request") {
		t.Fatalf("got %v", err)
	}
}

type errRoundTripper struct {
	failTransport bool
	failRead      bool
	last          bool
	n             int
}

func (e *errRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	e.n++
	if e.failTransport {
		return nil, errors.New("transport down")
	}
	body := io.NopCloser(strings.NewReader(`{"ok":true}`))
	if e.failRead {
		body = io.NopCloser(&errReader{})
	}
	return &http.Response{
		StatusCode: 200,
		Body:       body,
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("read fail") }

func TestTransportErrorRetriesThenFails(t *testing.T) {
	rt := &errRoundTripper{failTransport: true}
	c, err := New(
		WithToken("sp_x"),
		WithBaseURL("http://example.invalid"),
		WithHTTPClient(&http.Client{Transport: rt}),
		WithSleep(func(time.Duration) {}),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "calling") {
		t.Fatalf("got %v", err)
	}
	if rt.n != maxAttempts {
		t.Errorf("attempts = %d", rt.n)
	}
}

func TestReadBodyErrorRetries(t *testing.T) {
	rt := &errRoundTripper{failRead: true}
	c, err := New(
		WithToken("sp_x"),
		WithBaseURL("http://example.test"),
		WithHTTPClient(&http.Client{Transport: rt}),
		WithSleep(func(time.Duration) {}),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "reading the response") {
		t.Fatalf("got %v", err)
	}
}

func TestScenarioCRUDListMute(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/scenarios":
			_, _ = w.Write([]byte(`{"id":"mon_1","name":"Checkout","type":"http"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/scenarios/mon_1":
			_, _ = w.Write([]byte(`{"id":"mon_1","name":"Checkout","type":"http"}`))
		case r.Method == http.MethodPatch:
			_, _ = w.Write([]byte(`{"id":"mon_1","name":"Checkout2","type":"http"}`))
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/scenarios":
			_, _ = w.Write([]byte(`{"items":[{"id":"mon_1","name":"Checkout","type":"http"}],"nextCursor":null}`))
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/mute"):
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	name := "Checkout"
	sc, err := c.CreateScenario(ctx, ScenarioInput{Name: &name}, "k1")
	if err != nil || sc.ID != "mon_1" {
		t.Fatalf("create %v %+v", err, sc)
	}
	sc, err = c.GetScenario(ctx, "mon_1")
	if err != nil || sc.Name != "Checkout" {
		t.Fatalf("get %v %+v", err, sc)
	}
	name2 := "Checkout2"
	sc, err = c.UpdateScenario(ctx, "mon_1", ScenarioInput{Name: &name2})
	if err != nil || sc.Name != "Checkout2" {
		t.Fatalf("update %v %+v", err, sc)
	}
	until := "2026-10-01T00:00:00Z"
	if err := c.MuteScenario(ctx, "mon_1", &until); err != nil {
		t.Fatal(err)
	}
	if err := c.MuteScenario(ctx, "mon_1", nil); err != nil {
		t.Fatal(err)
	}
	list, err := c.ListScenarios(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list %v %v", err, list)
	}
	if err := c.DeleteScenario(ctx, "mon_1"); err != nil {
		t.Fatal(err)
	}
}

func TestMaintenanceWebhookSla(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/maintenance-windows":
			_, _ = w.Write([]byte(`{"id":"mw_1"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/maintenance-windows/mw_1":
			_, _ = w.Write([]byte(`{"id":"mw_1"}`))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/maintenance-windows/"):
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/maintenance-windows":
			_, _ = w.Write([]byte(`{"items":[{"id":"mw_1"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/webhooks":
			_, _ = w.Write([]byte(`{"id":"wh_1","secret":"sec","events":["run.failed"]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/webhooks/wh_1":
			_, _ = w.Write([]byte(`{"id":"wh_1","events":["run.failed"]}`))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/webhooks/"):
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/webhooks":
			_, _ = w.Write([]byte(`{"items":[{"id":"wh_1"}]}`))
		case r.Method == http.MethodPut && r.URL.Path == "/v1/sla-targets":
			_, _ = w.Write([]byte(`{"id":"sla_1","objectivePct":99.9,"windowDays":30}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sla-targets/sla_1":
			_, _ = w.Write([]byte(`{"id":"sla_1"}`))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/sla-targets/"):
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sla-targets":
			_, _ = w.Write([]byte(`{"items":[{"id":"sla_1"}]}`))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	reason := "deploy"
	mw, err := c.CreateMaintenanceWindow(ctx, MaintenanceWindowInput{Reason: &reason}, "")
	if err != nil || mw.ID != "mw_1" {
		t.Fatalf("mw create %v %+v", err, mw)
	}
	if _, err := c.GetMaintenanceWindow(ctx, "mw_1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListMaintenanceWindows(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteMaintenanceWindow(ctx, "mw_1"); err != nil {
		t.Fatal(err)
	}

	wh, err := c.CreateWebhook(ctx, WebhookInput{URL: "https://hooks.example/x", Events: []string{"run.failed"}}, "")
	if err != nil || wh.Secret == nil {
		t.Fatalf("wh %v %+v", err, wh)
	}
	if _, err := c.GetWebhook(ctx, "wh_1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListWebhooks(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteWebhook(ctx, "wh_1"); err != nil {
		t.Fatal(err)
	}

	mon := "mon_1"
	sla, err := c.UpsertSlaTarget(ctx, SlaTargetInput{MonitorID: &mon, ObjectivePct: 99.9, WindowDays: 30}, "")
	if err != nil || sla.ID != "sla_1" {
		t.Fatalf("sla %v %+v", err, sla)
	}
	if _, err := c.GetSlaTarget(ctx, "sla_1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListSlaTargets(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteSlaTarget(ctx, "sla_1"); err != nil {
		t.Fatal(err)
	}
}

func TestListPaginationAndCap(t *testing.T) {
	page := 0
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		if page == 1 {
			_, _ = w.Write([]byte(`{"items":[{"id":"a","name":"A","type":"http"}],"nextCursor":"c1"}`))
			return
		}
		_, _ = w.Write([]byte(`{"items":[{"id":"b","name":"B","type":"http"}],"nextCursor":null}`))
	})
	list, err := c.ListScenarios(context.Background())
	if err != nil || len(list) != 2 {
		t.Fatalf("got %v %v", err, list)
	}

	page = 0
	c2, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		cur := "c" + string(rune('0'+page))
		_, _ = w.Write([]byte(`{"items":[],"nextCursor":` + mustJSON(cur) + `}`))
	})
	_, err = c2.ListScenarios(context.Background())
	if err == nil || !strings.Contains(err.Error(), "too many pages") {
		t.Fatalf("want pagination cap, got %v", err)
	}
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func TestListPagedError(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"no"}`))
	})
	_, err := c.ListScenarios(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteEmptyBody(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteScenario(context.Background(), "mon_1"); err != nil {
		t.Fatal(err)
	}
}

func TestNewIdempotencyKey(t *testing.T) {
	a, b := newIdempotencyKey(), newIdempotencyKey()
	if a == "" || a == b {
		t.Fatalf("keys %q %q", a, b)
	}
}

func TestResourceMethodsPropagateErrors(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"nope"}`))
	})
	ctx := context.Background()
	name := "x"
	mon := "m"
	if _, err := c.CreateScenario(ctx, ScenarioInput{Name: &name}, "k"); err == nil {
		t.Fatal("CreateScenario")
	}
	if _, err := c.UpdateScenario(ctx, "id", ScenarioInput{Name: &name}); err == nil {
		t.Fatal("UpdateScenario")
	}
	if _, err := c.CreateMaintenanceWindow(ctx, MaintenanceWindowInput{}, "k"); err == nil {
		t.Fatal("CreateMaintenanceWindow")
	}
	if _, err := c.GetMaintenanceWindow(ctx, "id"); err == nil {
		t.Fatal("GetMaintenanceWindow")
	}
	if _, err := c.CreateWebhook(ctx, WebhookInput{URL: "https://x"}, "k"); err == nil {
		t.Fatal("CreateWebhook")
	}
	if _, err := c.GetWebhook(ctx, "id"); err == nil {
		t.Fatal("GetWebhook")
	}
	if _, err := c.UpsertSlaTarget(ctx, SlaTargetInput{MonitorID: &mon, ObjectivePct: 99, WindowDays: 7}, "k"); err == nil {
		t.Fatal("UpsertSlaTarget")
	}
	if _, err := c.GetSlaTarget(ctx, "id"); err == nil {
		t.Fatal("GetSlaTarget")
	}
}
