package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Goalt/tg-channel-to-rss/internal/app"
)

func TestNewServerHandlerServesFeedLanding(t *testing.T) {
	handler := newServerHandler(app.NewService(http.DefaultClient), nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req.Host = "example.com"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("unexpected content type: %q", got)
	}
	body := rec.Body.String()
	checks := []string{
		"tg-channel-to-rss API",
		"Datastar",
		"http://example.com/feed/telegram",
		"data-demo-form",
	}
	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Fatalf("expected landing page to contain %q", check)
		}
	}
}

func TestNewServerHandlerServesFeedLandingWithForwardedProto(t *testing.T) {
	handler := newServerHandler(app.NewService(http.DefaultClient), nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/feed/", nil)
	req.Host = "api.example.com"
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "https://api.example.com/feed/telegram") {
		t.Fatalf("expected https absolute URL in landing page, got %q", rec.Body.String())
	}
}

func TestNewServerHandlerServesFeedLandingWithAPIBaseEnv(t *testing.T) {
	t.Setenv("API_BASE_ENV", "https://public.example.com/base/")

	handler := newServerHandler(app.NewService(http.DefaultClient), nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/feed", nil)
	req.Host = "internal.example.local"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "https://public.example.com/base/feed/telegram") {
		t.Fatalf("expected API_BASE_ENV URL in landing page, got %q", rec.Body.String())
	}
}

func TestNewServerHandlerServesFeedJSON(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/s/testch1" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`
<html>
  <head>
    <title>Test Channel</title>
    <meta property="og:description" content="Test Description" />
  </head>
  <body>
    <div class="tgme_widget_message_bubble">
      <a class="tgme_widget_message_date" href="https://t.me/testch1/123">date</a>
      <time class="time" datetime="2026-04-21T10:00:00+00:00"></time>
      <div class="tgme_widget_message_text">Visit https://example.com and hi</div>
    </div>
  </body>
</html>`))
	}))
	defer upstream.Close()

	svc := app.NewService(upstream.Client())
	svc.BaseURL = upstream.URL
	svc.Now = func() time.Time { return time.Date(2026, 4, 21, 15, 0, 0, 0, time.UTC) }

	handler := newServerHandler(svc, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/feed/testch1", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("unexpected content type: %q", got)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"title":"Test Channel"`) {
		t.Fatalf("expected JSON feed body, got %q", body)
	}
}
