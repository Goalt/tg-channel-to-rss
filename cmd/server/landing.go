package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const datastarCDNURL = "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-RC.7/bundles/datastar.js"

type feedLandingData struct {
	BaseURL        string
	DatastarCDNURL string
	ExampleChannel string
}

var feedLandingTemplate = template.Must(template.New("feed-landing").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>tg-channel-to-rss API</title>
  <meta name="description" content="Interactive landing page for the /feed API with live examples and a Datastar-powered demo." />
  <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
  <script type="module" src="{{ .DatastarCDNURL }}"></script>
  <style>
    :root {
      --bg: #f5efe1;
      --bg-strong: #ead8b5;
      --surface: rgba(255, 250, 241, 0.86);
      --surface-strong: rgba(255, 247, 232, 0.96);
      --line: rgba(45, 43, 39, 0.12);
      --line-strong: rgba(45, 43, 39, 0.22);
      --text: #1f1b16;
      --muted: #5a544d;
      --accent: #0f766e;
      --accent-strong: #0b5f58;
      --accent-soft: rgba(15, 118, 110, 0.12);
      --shadow: 0 24px 70px rgba(44, 30, 9, 0.14);
      --radius-xl: 32px;
      --radius-lg: 24px;
      --radius-md: 18px;
      --radius-sm: 14px;
      --content: 1180px;
      --mono: "IBM Plex Mono", "SFMono-Regular", Consolas, monospace;
      --sans: "Space Grotesk", "Avenir Next", "Segoe UI", sans-serif;
    }

    * { box-sizing: border-box; }

    html {
      scroll-behavior: smooth;
    }

    body {
      margin: 0;
      font-family: var(--sans);
      color: var(--text);
      background:
        radial-gradient(circle at top left, rgba(252, 211, 77, 0.28), transparent 34%),
        radial-gradient(circle at 82% 18%, rgba(20, 184, 166, 0.2), transparent 26%),
        linear-gradient(180deg, #fff7e8 0%, #f4ecde 55%, #efe2cf 100%);
      min-height: 100vh;
    }

    body::before {
      content: "";
      position: fixed;
      inset: 0;
      pointer-events: none;
      background-image:
        linear-gradient(rgba(31, 27, 22, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(31, 27, 22, 0.03) 1px, transparent 1px);
      background-size: 36px 36px;
      mask-image: radial-gradient(circle at center, black 42%, transparent 86%);
    }

    a {
      color: inherit;
      text-decoration: none;
    }

    .shell {
      width: min(calc(100% - 32px), var(--content));
      margin: 0 auto;
      padding: 24px 0 56px;
      position: relative;
      z-index: 1;
    }

    .topbar,
    .hero,
    .grid,
    .foot {
      animation: rise 600ms ease both;
    }

    .topbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 18px;
      padding: 12px 0 20px;
      color: var(--muted);
      font-size: 14px;
    }

    .brand {
      display: inline-flex;
      align-items: center;
      gap: 12px;
      font-weight: 700;
      letter-spacing: 0.03em;
      text-transform: uppercase;
    }

    .brand-logo {
      width: 50px;
      height: 50px;
      display: block;
      border-radius: 16px;
      border: 1px solid rgba(15, 118, 110, 0.12);
      box-shadow: 0 14px 28px rgba(15, 118, 110, 0.12);
      background: rgba(255, 255, 255, 0.92);
      padding: 4px;
    }

    .hero {
      display: grid;
      grid-template-columns: minmax(0, 1.15fr) minmax(320px, 0.85fr);
      gap: 24px;
      align-items: stretch;
    }

    .panel {
      background: var(--surface);
      border: 1px solid var(--line);
      border-radius: var(--radius-xl);
      box-shadow: var(--shadow);
      backdrop-filter: blur(14px);
    }

    .hero-copy {
      padding: 34px;
    }

    .eyebrow {
      display: inline-flex;
      gap: 10px;
      align-items: center;
      padding: 10px 14px;
      border-radius: 999px;
      background: rgba(255, 255, 255, 0.64);
      border: 1px solid var(--line);
      font-size: 12px;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      color: var(--accent-strong);
    }

    h1,
    h2,
    h3,
    p {
      margin: 0;
    }

    h1 {
      margin-top: 18px;
      font-size: clamp(40px, 7vw, 76px);
      line-height: 0.96;
      letter-spacing: -0.05em;
      max-width: 10ch;
    }

    .hero-lead {
      margin-top: 18px;
      max-width: 58ch;
      color: var(--muted);
      font-size: 18px;
      line-height: 1.55;
    }

    .hero-actions,
    .chips,
    .action-row,
    .stat-row,
    .foot {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      align-items: center;
    }

    .hero-actions {
      margin-top: 24px;
    }

    .button,
    .chip {
      border: 1px solid transparent;
      border-radius: 999px;
      transition: transform 180ms ease, background 180ms ease, border-color 180ms ease;
      cursor: pointer;
      font: inherit;
    }

    .button {
      padding: 13px 18px;
      font-weight: 700;
    }

    .button:hover,
    .chip:hover {
      transform: translateY(-1px);
    }

    .button.primary {
      background: linear-gradient(135deg, var(--accent), #14b8a6);
      color: #f7f9f8;
      box-shadow: 0 14px 28px rgba(15, 118, 110, 0.22);
    }

    .button.ghost,
    .chip {
      background: rgba(255, 255, 255, 0.56);
      border-color: var(--line);
      color: var(--text);
    }

    .button.ghost {
      padding-inline: 16px;
    }

    .hero-card {
      padding: 24px;
      display: grid;
      grid-template-rows: auto auto 1fr;
      gap: 18px;
      background:
        linear-gradient(160deg, rgba(255, 255, 255, 0.94), rgba(246, 236, 218, 0.85)),
        var(--surface-strong);
    }

    .mono {
      font-family: var(--mono);
      font-size: 13px;
      line-height: 1.6;
    }

    .hero-card .mono {
      padding: 16px;
      border-radius: var(--radius-md);
      background: #111827;
      color: #f3f4f6;
      overflow-x: auto;
      box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08);
    }

    .stat-row {
      margin-top: 16px;
    }

    .stat {
      flex: 1 1 160px;
      padding: 14px 16px;
      border-radius: var(--radius-md);
      background: rgba(255, 255, 255, 0.56);
      border: 1px solid var(--line);
    }

    .stat strong {
      display: block;
      font-size: 28px;
      line-height: 1;
      margin-bottom: 4px;
    }

    .grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 18px;
      margin-top: 22px;
    }

    .card {
      padding: 24px;
      border-radius: var(--radius-lg);
      background: rgba(255, 252, 246, 0.82);
      border: 1px solid var(--line);
      box-shadow: 0 18px 40px rgba(57, 38, 14, 0.08);
      display: grid;
      gap: 14px;
      min-height: 100%;
    }

    .card h2,
    .card h3 {
      font-size: 22px;
      letter-spacing: -0.03em;
    }

    .card p,
    .label,
    .note,
    .status,
    .foot {
      color: var(--muted);
    }

    .label {
      font-size: 12px;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      font-weight: 700;
    }

    .stack {
      display: grid;
      gap: 12px;
    }

    .surface {
      padding: 16px;
      border-radius: var(--radius-md);
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.66);
    }

    .surface code,
    pre {
      font-family: var(--mono);
      font-size: 13px;
      margin: 0;
      white-space: pre-wrap;
      word-break: break-word;
    }

    .demo {
      margin-top: 22px;
      padding: 26px;
      border-radius: var(--radius-xl);
      background:
        linear-gradient(180deg, rgba(20, 20, 20, 0.95), rgba(34, 34, 34, 0.96)),
        #111827;
      color: #f9fafb;
      border: 1px solid rgba(255, 255, 255, 0.08);
      box-shadow: 0 28px 80px rgba(20, 20, 20, 0.28);
    }

    .demo-grid {
      display: grid;
      grid-template-columns: minmax(280px, 360px) minmax(0, 1fr);
      gap: 18px;
      align-items: start;
    }

    .demo-panel {
      padding: 18px;
      border-radius: var(--radius-lg);
      background: rgba(255, 255, 255, 0.05);
      border: 1px solid rgba(255, 255, 255, 0.08);
    }

    .field {
      display: grid;
      gap: 8px;
      margin-top: 16px;
    }

    .field input {
      width: 100%;
      border-radius: 16px;
      border: 1px solid rgba(255, 255, 255, 0.14);
      background: rgba(255, 255, 255, 0.08);
      color: #f9fafb;
      padding: 14px 16px;
      font: inherit;
      outline: none;
      transition: border-color 180ms ease, transform 180ms ease, background 180ms ease;
    }

    .field input:focus {
      border-color: rgba(45, 212, 191, 0.72);
      background: rgba(255, 255, 255, 0.1);
      transform: translateY(-1px);
    }

    .demo .chip,
    .demo .button.ghost {
      color: #f9fafb;
      background: rgba(255, 255, 255, 0.08);
      border-color: rgba(255, 255, 255, 0.12);
    }

    .status {
      min-height: 22px;
      margin-top: 12px;
      color: #cbd5e1;
    }

    .status.is-error {
      color: #fda4af;
    }

    .status.is-success {
      color: #86efac;
    }

    .demo pre {
      min-height: 380px;
      max-height: 520px;
      overflow: auto;
      padding: 18px;
      border-radius: var(--radius-lg);
      background: rgba(15, 23, 42, 0.88);
      border: 1px solid rgba(255, 255, 255, 0.08);
      color: #e2e8f0;
    }

    .foot {
      justify-content: space-between;
      margin-top: 18px;
      font-size: 14px;
    }

    @keyframes rise {
      from {
        opacity: 0;
        transform: translateY(10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    @media (max-width: 1080px) {
      .hero,
      .demo-grid,
      .grid {
        grid-template-columns: 1fr;
      }
    }

    @media (max-width: 720px) {
      .shell {
        width: min(calc(100% - 20px), var(--content));
        padding-top: 14px;
      }

      .hero-copy,
      .hero-card,
      .card,
      .demo {
        padding: 20px;
      }

      h1 {
        max-width: none;
      }

      .topbar,
      .foot {
        justify-content: flex-start;
      }
    }
  </style>
</head>
<body data-store="{ channel: '{{ .ExampleChannel }}' }">
  <div class="shell">
    <div class="topbar">
      <div class="brand">
      <img class="brand-logo" src="/logo.svg" alt="tg-channel-to-rss logo" />
        <span>tg-channel-to-rss</span>
      </div>
      <div>Telegram channel to JSON feed, with a live Datastar demo.</div>
    </div>

    <section class="hero">
      <div class="panel hero-copy">
        <div class="eyebrow">/feed API landing</div>
        <h1>Turn a public Telegram channel into a JSON feed.</h1>
        <p class="hero-lead">
          Call <span class="mono">/feed/{channel_name}</span>, get structured JSON,
          and test it from this page without leaving the browser.
        </p>
        <div class="hero-actions">
          <a class="button primary" href="#demo">Run live demo</a>
          <a class="button ghost" id="open-feed-link" href="/feed/{{ .ExampleChannel }}" target="_blank" rel="noreferrer">Open current feed</a>
        </div>
        <div class="stat-row">
          <div class="stat">
            <strong>GET</strong>
            <span>Single endpoint for direct JSON consumption.</span>
          </div>
          <div class="stat">
            <strong>60s</strong>
            <span>Public cache window returned by the API.</span>
          </div>
          <div class="stat">
            <strong>Datastar</strong>
            <span>Landing interactivity with declarative data-* bindings.</span>
          </div>
        </div>
      </div>

      <aside class="panel hero-card">
        <div>
          <p class="label">Request shape</p>
          <div class="mono">GET {{ .BaseURL }}/feed/<span data-text="$channel.trim().replace('@', '')"></span></div>
        </div>
        <div>
          <p class="label">What comes back</p>
          <p class="note">A JSON document with channel metadata and parsed Telegram posts.</p>
        </div>
        <div class="mono">{
  "title": "Channel title",
  "link": "https://t.me/s/...",
  "description": "Posts from ...",
  "created": "2026-05-01T10:00:00Z",
  "items": [
    {
      "title": "New post in channel @...",
      "link": "https://t.me/s/.../123",
      "content": "..."
    }
  ]
}</div>
      </aside>
    </section>

    <section class="grid" aria-label="usage examples">
      <article class="card">
        <div class="stack">
          <p class="label">Curl</p>
          <h2>Use it from shell scripts and cron jobs.</h2>
          <p>Pipe the JSON to jq, store it, or fan it out to other systems.</p>
        </div>
        <div class="surface"><code>curl '{{ .BaseURL }}/feed/telegram'</code></div>
      </article>

      <article class="card">
        <div class="stack">
          <p class="label">Frontend</p>
          <h2>Fetch it directly from a browser app.</h2>
          <p>Use relative URLs in development, or an absolute host when deployed.</p>
        </div>
        <div class="surface"><code>fetch('/feed/telegram')
  .then((response) =&gt; response.json())
  .then((feed) =&gt; console.log(feed.items.length))</code></div>
      </article>

      <article class="card">
        <div class="stack">
          <p class="label">Direct URL</p>
          <h2>Open the feed in a browser or pass it to tooling.</h2>
          <p>Public channels only. Sensitive or protected channels may return limited data.</p>
        </div>
        <div class="surface"><code>{{ .BaseURL }}/feed/telegram</code></div>
      </article>
    </section>

    <section class="demo" id="demo">
      <div class="demo-grid">
        <div class="demo-panel">
          <p class="label">Live demo</p>
          <h3 style="margin-top: 10px; font-size: 30px; letter-spacing: -0.04em;">Preview the API response.</h3>
          <p class="note" style="margin-top: 10px; color: #cbd5e1;">
            Enter a public Telegram channel name. The page will call the real endpoint and pretty-print the JSON.
          </p>

          <form data-demo-form>
            <label class="field" for="channel-input">
              <span class="label" style="color: #e2e8f0;">Channel name</span>
              <input
                id="channel-input"
                name="channel"
                value="{{ .ExampleChannel }}"
                data-bind:channel
                spellcheck="false"
                autocomplete="off"
                placeholder="telegram"
              />
            </label>

            <div class="chips" style="margin-top: 14px;">
              <button class="chip" type="button" data-channel-chip="telegram" data-on:click="$channel = 'telegram'">telegram</button>
              <button class="chip" type="button" data-channel-chip="durov" data-on:click="$channel = 'durov'">durov</button>
              <button class="chip" type="button" data-channel-chip="golang" data-on:click="$channel = 'golang'">golang</button>
            </div>

            <div class="surface" style="margin-top: 16px; background: rgba(255, 255, 255, 0.06); border-color: rgba(255, 255, 255, 0.1);">
              <p class="label" style="color: #cbd5e1; margin-bottom: 8px;">Request URL</p>
              <code id="request-url" data-text="'{{ .BaseURL }}/feed/' + $channel.trim().replace('@', '')"></code>
            </div>

            <div class="action-row" style="margin-top: 16px;">
              <button class="button primary" type="submit">Fetch JSON</button>
              <a class="button ghost" id="view-url-link" href="/feed/{{ .ExampleChannel }}" target="_blank" rel="noreferrer">Open in new tab</a>
            </div>
          </form>

          <div class="status" data-demo-status>Ready.</div>
        </div>

        <div class="demo-panel">
          <p class="label">Response</p>
          <pre data-demo-output>{
  "hint": "Submit the demo to fetch a live feed"
}</pre>
        </div>
      </div>

      <div class="foot">
        <span>Datastar powers the reactive bits. The fetch demo itself uses the same <code>/feed/{channel}</code> route your clients call.</span>
        <span class="mono">{{ .BaseURL }}/feed/telegram</span>
      </div>

      <div class="surface" style="margin-top: 16px; background: rgba(255, 255, 255, 0.06); border-color: rgba(255, 255, 255, 0.1); color: #e2e8f0;">
        <p class="label" style="color: #cbd5e1; margin-bottom: 8px;">Service creator</p>
        <p style="margin: 0; line-height: 1.6;">
          This service was developed by
          <a href="https://contractsguard.com" target="_blank" rel="noreferrer" style="color: #99f6e4; font-weight: 700; text-decoration: underline; text-underline-offset: 3px;">
            contractsguard.com
          </a>.
        </p>
      </div>
    </section>
  </div>

  <script>
    document.addEventListener("DOMContentLoaded", function () {
      var form = document.querySelector("[data-demo-form]");
      var input = document.getElementById("channel-input");
      var output = document.querySelector("[data-demo-output]");
      var status = document.querySelector("[data-demo-status]");
      var openLink = document.getElementById("open-feed-link");
      var viewLink = document.getElementById("view-url-link");
      var requestURL = document.getElementById("request-url");
      var chips = document.querySelectorAll("[data-channel-chip]");
      var baseURL = {{ printf "%q" .BaseURL }};

      function normalizeChannel(value) {
        return value.trim().replace(/^@+/, "");
      }

      function buildPath(channel) {
        return "/feed/" + encodeURIComponent(channel);
      }

      function syncLinks() {
        var channel = normalizeChannel(input.value || "");
        var path = channel ? buildPath(channel) : "/feed";
        var absoluteURL = channel ? baseURL + path : baseURL + "/feed";

        if (openLink) {
          openLink.href = path;
        }
        if (viewLink) {
          viewLink.href = path;
        }
        if (requestURL) {
          requestURL.textContent = absoluteURL;
        }
      }

      function setStatus(message, kind) {
        status.textContent = message;
        status.classList.remove("is-error", "is-success");
        if (kind) {
          status.classList.add(kind);
        }
      }

      chips.forEach(function (chip) {
        chip.addEventListener("click", function () {
          input.value = chip.getAttribute("data-channel-chip") || "";
          input.dispatchEvent(new Event("input", { bubbles: true }));
          syncLinks();
        });
      });

      input.addEventListener("input", syncLinks);
      syncLinks();

      form.addEventListener("submit", async function (event) {
        event.preventDefault();

        var channel = normalizeChannel(input.value || "");
        if (!channel) {
          setStatus("Enter a public Telegram channel name first.", "is-error");
          output.textContent = '{\n  "error": "missing channel name"\n}';
          return;
        }

        var path = buildPath(channel);
        setStatus("Fetching " + path + " ...");
        output.textContent = '{\n  "loading": true\n}';

        try {
          var response = await fetch(path, {
            headers: { "Accept": "application/json" }
          });
          var raw = await response.text();
          var pretty = raw;

          try {
            pretty = JSON.stringify(JSON.parse(raw), null, 2);
          } catch (parseError) {
          }

          output.textContent = pretty;
          if (!response.ok) {
            setStatus("Request failed with HTTP " + response.status + ".", "is-error");
            return;
          }

          setStatus("Success. JSON rendered from the live endpoint.", "is-success");
        } catch (error) {
          output.textContent = '{\n  "error": ' + JSON.stringify(String(error && error.message ? error.message : error)) + '\n}';
          setStatus("Network error while calling the API.", "is-error");
        }
      });
    });
  </script>
</body>
</html>
`))

func serveFeedLanding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")

	if err := feedLandingTemplate.Execute(w, feedLandingData{
		BaseURL:        requestBaseURL(r),
		DatastarCDNURL: datastarCDNURL,
		ExampleChannel: "telegram",
	}); err != nil {
		http.Error(w, "Failed to render landing page", http.StatusInternalServerError)
	}
}

func serveLogoSVG(w http.ResponseWriter, r *http.Request) {
	for _, candidate := range []string{"logo.svg", filepath.Join("..", "..", "logo.svg")} {
		content, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Content-Type", "image/svg+xml; charset=UTF-8")
		_, _ = w.Write(content)
		return
	}

	http.NotFound(w, r)
}

func requestBaseURL(r *http.Request) string {
	if override := strings.TrimSpace(os.Getenv("API_BASE_ENV")); override != "" {
		return strings.TrimRight(override, "/")
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwarded != "" {
		scheme = strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}

	host := strings.TrimSpace(r.Host)
	if host == "" {
		host = "localhost"
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}
