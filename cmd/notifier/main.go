package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Goalt/tg-channel-to-rss/internal/app"
	"github.com/Goalt/tg-channel-to-rss/internal/notifier"
)

func main() {
	channels := splitList(os.Getenv("TG_CHANNELS"))
	webhooks := splitList(os.Getenv("WEBHOOKS"))

	if len(channels) == 0 {
		log.Fatalf("TG_CHANNELS must be a non-empty comma-separated list of Telegram channel names")
	}
	if len(webhooks) == 0 {
		log.Fatalf("WEBHOOKS must be a non-empty comma-separated list of webhook URLs")
	}

	interval, err := time.ParseDuration(envOrDefault("POLL_INTERVAL", "5m"))
	if err != nil {
		log.Fatalf("invalid POLL_INTERVAL: %v", err)
	}

	svc := app.NewService(&http.Client{Timeout: 30 * time.Second})
	n := notifier.New(notifier.Config{
		Channels:    channels,
		Webhooks:    webhooks,
		Interval:    interval,
		HTTPTimeout: 30 * time.Second,
	}, svc, nil, nil)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("tg-notifier: polling %d channel(s) every %s, dispatching to %d webhook(s)", len(channels), interval, len(webhooks))
	if err := n.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("notifier stopped: %v", err)
	}
}

func splitList(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envOrDefault(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
