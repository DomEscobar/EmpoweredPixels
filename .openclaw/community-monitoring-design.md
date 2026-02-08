# Community Monitoring Architecture — Design Doc

**Project:** EmpoweredPixels Community Analytics
**Author:** Community Analyst (analyst agent)
**Date:** 2025-02-08
**Status:** Draft for Review

---

## 1. Overview

This system monitors community sentiment across Discord, Reddit, and app stores (Apple App Store, Google Play Store). It aggregates feedback, identifies trending topics and feature requests, and produces a weekly report delivered to the internal feedback channel.

**Goals:**
- Detect emerging bugs and pain points early
- Identify high-demand features from player requests
- Track player satisfaction trends over time
- Inform product roadmap and prioritization

**Non-goals:**
- Real-time response to community issues (that's for a community manager)
- Individual user support or moderation
- Deep NLP analysis (use lightweight, interpretable methods)

---

## 2. System Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   Discord       │    │    Reddit        │    │  App Stores      │
│   (Bot/API)     │    │   (PRAW/API)     │    │  (AppBot/API)    │
└────────┬────────┘    └────────┬─────────┘    └────────┬─────────┘
         │                      │                       │
         └───────────┬──────────┴───────────┬───────────┘
                     ▼                      ▼
         ┌─────────────────────────────────────────────┐
         │        Collector Service (Go)              │
         │  • Fetch new posts/comments/reviews        │
         │  • Normalize into common schema            │
         │  • Filter relevant content (game-related)  │
         │  • Basic deduplication                     │
         └────────────────┬────────────────────────────┘
                          │
                          ▼
         ┌─────────────────────────────────────────────┐
         │         Storage Layer                       │
         │  ┌─────────────────────────────────────┐   │
         │  │ PostgreSQL:                        │   │
         │  │  • raw_community_posts             │   │
         │  │  • processed_sentiment             │   │
         │  │  • weekly_reports                  │   │
         │  └─────────────────────────────────────┘   │
         └────────────────┬────────────────────────────┘
                          │
                          ▼
         ┌─────────────────────────────────────────────┐
         │     Analyzer Service (Go)                  │
         │  • Sentiment analysis (keyword + ML API)   │
         │  • Topic clustering                       │
         │  • Feature request extraction             │
         │  • Trend detection (vs prior week)        │
         └────────────────┬────────────────────────────┘
                          │
                          ▼
         ┌─────────────────────────────────────────────┐
         │    Reporter Service (Go)                  │
         │  • Generate weekly summary (Telegram)     │
         │  • Save JSON report                       │
         │  • Update dashboard data (optional)       │
         └─────────────────────────────────────────────┘
```

**Services:** All built as Go binaries within existing `backend/` module. New package: `internal/community/`.

**Deployment:** Run as cron-triggered jobs via OpenClaw scheduler or systemd timers. Each service can run independently but orchestrated weekly.

---

## 3. Data Sources & Authentication

### 3.1 Discord

**Method:** Discord Gateway API with a Bot token.
- Create a bot in Discord Developer Portal
- Invite bot to relevant guilds (Discord servers): EmpoweredPixels official, community-run servers
- Required scopes: `bot`, `applications.commands`
- Required permissions: `Read Messages/View Channels`, `Send Messages`, `Read Message History`, `Embed Links`

**Endpoints:**
- `GET /guilds/{guild.id}/channels` → list text channels to monitor
- `GET /channels/{channel.id}/messages` → fetch messages (requires `Read Message History`)
- Webhook alternative: if bot not feasible, use Discord webhook to push community posts to our system (requires trusted community managers to set up)

**Rate limits:** ~50 requests/second per bot; message fetch limited to last 100 messages per request. Need incremental fetching via message IDs.

**Configuration file:** `.openclaw/secrets/discord.json`
```json
{
  "bot_token": "MT...",
  "guild_ids": ["123456789012345678"],
  "monitored_channel_ids": ["987654321098765432"],
  "last_message_id": "previous_run_checkpoint"
}
```

### 3.2 Reddit

**Method:** Reddit API via PRAW (Python) or pure HTTP from Go. Prefer Go for consistency.

- Create Reddit app at https://www.reddit.com/prefs/apps
- Choose "script" type
- Credentials: client_id, client_secret, username, password (for OAuth2)

**Target subreddits:**
- r/EtherealIron (official subreddit, if exists)
- r/IndieGaming, r/gaming, r/MMORPG (if game posts appear)
- r/boostforgames (or relevant platforms)

**Endpoints:**
- `GET /r/{subreddit}/new` → new posts
- `GET /comments/{post_id}` → comments on posts
- Search endpoint: `GET /search?q=empoweredpixels|etherealiron` → catch mentions

**Rate limits:** ~60 requests/minute for script auth.

**Configuration file:** `.openclaw/secrets/reddit.json`
```json
{
  "client_id": "abc123...",
  "client_secret": "xyz789...",
  "username": "empoweredpixels_bot",
  "password": "vault_password",
  "subreddits": ["EtherealIron", "IndieGaming"],
  "last_post_ids": {}
}
```

### 3.3 App Stores

**Apple App Store:**
- Use App Store Connect API (JWT auth) to fetch reviews.
- Requires: private key, key ID, issuer ID, vendor number (Team ID)
- Limited to 20 requests/minute, reviews available for app versions.
- Alternative: third-party services like Appbot.io or AppSkeletons? Prefer direct Apple API if feasible, else use Appbot API key if budget allows.

**Google Play Store:**
- Use Google Play Developer API (OAuth2 service account)
- Requires: service account key JSON, package name
- Can list reviews and reply; rate limits reasonable (~1000/day)

**Simpler alternative:** Use a review-monitoring SaaS (e.g., Appbot, AppFollow) that provides API access for all stores. Costs ~$50-100/month but simplifies integration. Recommended for MVP.

**Configuration file:** `.openclaw/secrets/appstores.json`
```json
{
  "apple": {
    "key_id": "ABC123DEFG",
    "issuer_id": "987654321-1234567890",
    "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEv...",
    "app_id": "1234567890"
  },
  "google": {
    "service_account_key": "/path/to/google-play-key.json",
    "package_name": "com.empoweredpixels.etherealiron"
  },
  "appbot": {
    "api_key": "appbot_...",
    "store_apps": ["ios:123456", "android:com.empoweredpixels.etherealiron"]
  }
}
```

---

## 4. Storage Strategy

### 4.1 Database Schema

Add to existing PostgreSQL database (same as game DB or separate analytics DB? Same DB is okay, separate schema `community` for isolation).

```sql
-- Raw collected posts (immutable)
CREATE TABLE community.raw_posts (
    id BIGSERIAL PRIMARY KEY,
    source VARCHAR(50) NOT NULL, -- 'discord', 'reddit', 'app_store_ios', 'app_store_android'
    external_id VARCHAR(255) NOT NULL, -- platform-specific ID (message.id, review.id)
    url TEXT, -- link to original (if public)
    author_id VARCHAR(255),
    author_name VARCHAR(255),
    content TEXT NOT NULL,
    metadata JSONB, -- platform-specific fields: channel, guild, subreddit, rating, version, etc.
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, external_id)
);

-- Normalized sentiment and analysis
CREATE TABLE community.posts (
    id BIGSERIAL PRIMARY KEY,
    raw_post_id BIGINT REFERENCES community.raw_posts(id),
    source VARCHAR(50) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255),
    content TEXT NOT NULL,
    platform_metadata JSONB, -- minimal platform fields
    sentiment_score SMALLINT, -- -100 to +100 (negative to positive)
    sentiment_label VARCHAR(20), -- 'negative', 'neutral', 'positive'
    topics TEXT[], -- array of detected topics: 'bug', 'feature-request', 'combat', 'graphics', 'ui', 'economy', etc.
    feature_requests TEXT[], -- extracted feature request phrases (or empty)
    is_feature_request BOOLEAN GENERATED ALWAYS AS (array_length(feature_requests, 1) > 0) STORED,
    posted_at TIMESTAMPTZ NOT NULL, -- original post timestamp
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, external_id)
);

-- Aggregated weekly reports
CREATE TABLE community.weekly_reports (
    id BIGSERIAL PRIMARY KEY,
    week_start DATE NOT NULL, -- Monday date
    generated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    total_posts INTEGER NOT NULL,
    sentiment_distribution JSONB, -- {negative: n, neutral: n, positive: n}
    top_topics JSONB, -- {topic: count, ...}
    feature_requests JSONB, -- {request: count, ...} top 10
    top_negative_posts JSONB, -- array of {source, external_id, author, content, sentiment}
    raw_summary TEXT, -- human-readable summary (Telegram message)
    UNIQUE(week_start)
);

-- Indexes for performance
CREATE INDEX idx_posts_source_collected ON community.posts(source, collected_at DESC);
CREATE INDEX idx_posts_sentiment ON community.posts(sentiment_score);
CREATE INDEX idx_posts_topics ON community.posts USING GIN(topics);
CREATE INDEX idx_raw_posts_unique ON community.raw_posts(source, external_id);
```

### 4.2 Checkpointing

To avoid re-fetching old posts, each collector stores its last fetch state:
- Option A: In `community.posts` by querying max `collected_at` or `external_id` per source.
- Option B: Dedicated table `community.collector_state` with columns `source`, `last_external_id`, `last_run_at`.

Simpler: use a file-based checkpoint per source (e.g., `var/community/checkpoint-discord.json`) since cron runs on single host. But database is more robust.

**Recommendation:** File-based checkpoints in `/.openclaw/community/checkpoints/` with simple JSON:
```json
{
  "discord": {"last_message_id": "123456789012345678", "last_run": "2025-02-08T07:00:00Z"},
  "reddit": {"last_post_ids": {"all": "abcdef", "subreddit1": "123abc"}},
  "app_store_ios": {"last_review_id": "999888777"},
  "app_store_android": {"last_review_id": "111222333"}
}
```

---

## 5. Collector Service

**Bin:** `backend/cmd/community-collector`
**Logic:**
1. Load config and credentials from `.openclaw/secrets/`
2. Load checkpoint
3. For each source:
   - Call platform API with `since` parameter (message ID > last, or timestamp > last_run - 7d)
   - Paginate until no more new items (limit to 7 days to avoid huge backfills)
   - Normalize each item into `CommunityPost` struct (source, external_id, author, content, posted_at, metadata)
   - Store in `community.raw_posts` (do NOT update if duplicate)
   - Convert to `community.posts` with placeholder sentiment=NULL, topics=[] (will be filled by analyzer)
4. Update checkpoint with latest IDs/timestamps
5. Log summary: fetched X posts, stored Y new.

**Error handling:**
- Retry with exponential backoff on network errors (max 3 retries)
- If API returns 401/403 → log critical error, alert (credentials expired)
- If rate limited → wait until reset (parse `X-RateLimit-Reset` headers or known bounds)
- Partial failure: continue with other sources; write partial checkpoint anyway to avoid re-fetching same data.

---

## 6. Analyzer Service

**Bin:** `backend/cmd/community-analyzer`
**Trigger:** After collector completes (can be chained or same binary with subcommand).

**Sentiment Analysis Approach:**

**Hybrid:**
- Phase 1 (MVP): Rule-based with keyword scoring.
  - Build lexicon of game-specific terms with polarities:
    - Positive: "love", "awesome", "great", "addictive", "balanced", "smooth", "fun"
    - Negative: "broken", "bug", "crash", "lag", "unfair", "grindy", "pay-to-win", "OP", "nerf"
    - Neutral/contextual: "update", "patch", "release", "server", "maintenance"
  - Algorithm: VADER-like (valence aware). Count positive/negative words, weight intensifiers ("very", "extremely"), handle negations ("not good").
  - Score range: -100 (extremely negative) to +100 (extremely positive).
  - Label mapping: score < -20 → negative; -20 to +20 → neutral; > +20 → positive. (Adjustable thresholds.)

- Phase 2 (Optional): Call external ML API (e.g., OpenAI, Cohere, or open-source hosted) for more accuracy. Keep as config flag. Cost consideration: ~$0.002 per post; weekly volume may be 100-500 posts → manageable.

**Feature Request Extraction:**
- Pattern matching:
  - "I wish...", "It would be great if...", "Should add...", "Need more...", "Hope they..." 
  - "Please implement", "Missing", "Would like to see"
- Keyword topics also help: if a negative post mentions "combat" and "balance", it's likely about combat balance, not a feature request.
- Store extracted phrases in `feature_requests` array (could be long; limit to 3 per post?).

**Topic Clustering:**
- Use keyword tags based on game domain knowledge:
  - combat, balance, pvp, pve, progression, economy, ui/ux, graphics, performance, matchmaking, guilds, events, shop, monetization, quality-of-life, bugs, crashes, server
- Simple approach: scan content for topic keywords (case-insensitive, plural/singular variants). Multiple topics allowed.
- More advanced: LDA or embedding clustering (overkill for MVP).

**Trend Detection:**
- Compare current week's counts vs prior 4 weeks average.
- Flag topics with >50% increase (spike) or sustained high volume.
- Identify "top negative posts" by lowest sentiment score (manual review by PM).

**Output:** Update `community.posts` rows with `sentiment_score`, `sentiment_label`, `topics`, `feature_requests`.

**Weekly Aggregation:**
- Compute:
  - Total posts, unique authors
  - Sentiment distribution (counts/percentages)
  - Top 5 topics by volume
  - Top 10 feature requests (by frequency)
  - Most negative posts (lowest 3 sentiment scores, length > 20 chars)
  - Trend: compare sentiment average vs last week
  - New sources or channels that appeared
- Generate human-readable markdown summary (Telegram-friendly).
- Insert into `community.weekly_reports` and also write to file: `/.openclaw/community/reports/weekly-YYYY-WW.json` and `.txt`.

---

## 7. Reporter Service

**Bin:** `backend/cmd/community-reporter` (or part of analyzer)

**Trigger:** Weekly cron (Monday 08:00 Europe/Berlin).

**Channels:**
- **Primary:** Telegram message to feedback channel (from config: `reporting.telegram_chat_id`).
  - Send formatted text with emojis and links to full report.
  - If >4000 chars, split or send file? Telegram supports up to 4096; summarize and send link to JSON.
- **Secondary:** Save JSON to file and/or post to internal dashboard (future).
- **Optional:** Message to PM (direct) if critical issues (spike in crashes, major负面 sentiment).

**Message template:**
```
📊 Community Weekly Report — Week of Feb 3, 2026

🗞️ Overview
• Total posts: 142 (↑12% vs last week)
• Sentiment: 😊 18% | 😐 52% | 😠 30%  (last week: 25% negative)
• Active channels: Discord (89), Reddit (38), App Store (15)

🔥 Top Topics
1. combat balance (24)
2. UI/UX (18)
3. server lag (12)
4. feature-request: new fighter types (11)
5. shop prices (9)

💡 Feature Requests (Top 5)
1. "add guild system" (8 mentions)
2. "more solo content" (6)
3. "fix matchmaking" (5)
4. "improve inventory management" (4)
5. "add spectator mode" (3)

⚠️ Negative Highlights
• Server lag spikes during peak EU hours (Discord #reports)
• Combat balance concerns: Fireball overpowered (Reddit)
• UI bug: tooltips not dismissing (App Store 2★ review)

📈 Trend: Sentiment improving after last patch. Monitor server stability.

Full report: file:///.openclaw/community/reports/weekly-2025-W06.json
```

---

## 8. Deployment Plan

### 8.1 Scheduling with OpenClaw

Use `cron` tool to schedule weekly runs:
- `community-collector` runs daily (e.g., 02:00) to keep data fresh.
- `community-analyzer` + `community-reporter` run Mondays 08:00.

Cron job configuration (to be added to `.openclaw/config.json` agents or standalone):

```json
{
  "jobs": [
    {
      "name": "Community Collector Daily",
      "schedule": { "kind": "cron", "expr": "0 2 * * *", "tz": "Europe/Berlin" },
      "payload": { "kind": "agentTurn", "message": "Run community collector" },
      "sessionTarget": "isolated",
      "enabled": true
    },
    {
      "name": "Community Weekly Report",
      "schedule": { "kind": "cron", "expr": "0 8 * * 1", "tz": "Europe/Berlin" },
      "payload": { "kind": "agentTurn", "message": "Run community analyzer and reporter" },
      "sessionTarget": "isolated",
      "enabled": true
    }
  ]
}
```

Alternatively, use systemd timers on the host. But OpenClaw cron integrates with agent model; we can spawn an `analyst` agent to run the commands directly.

### 8.2 Container/Serverless

Not necessary initially. The Go binaries run as ephemeral processes triggered by cron. They can be built as part of the existing `backend` Makefile.

### 8.3 Secrets Management

- Store in `.openclaw/secrets/` as JSON files (not committed).
- Add `.openclaw/secrets/*` to `.gitignore`.
- Document required files in setup script.
- Alternatively, use environment variables injected by OpenClaw gateway config (preferred for production). But file-based is simpler for MVP.

---

## 9. Error Handling & Observability

**Logging:**
- Use structured logging (JSON) to stdout.
- Log levels: info (normal), warning (recoverable), error (failed but continued), fatal (cannot proceed).
- Log file rotation via systemd journal or Docker if containerized (not yet).

**Metrics:** (optional)
- Counters: posts collected, sentiment scored, errors by type.
- Could write to `metrics` table or Prometheus endpoint later.

**Alerts:** If collector fails for >2 consecutive runs, send Telegram message to PM/diagnostic channel.

**Idempotency:** Collector checks `UNIQUE(source, external_id)` constraint; safe to retry. Analyzer updates `posts` table, can re-run on same data.

---

## 10. Implementation Tasks (for Coder/Foundry)

Once design approved, break into:

1. `TASK-910`: DB schema creation — add tables `community.raw_posts`, `community.posts`, `community.weekly_reports` plus indexes.
2. `TASK-911`: Collector service — Discord integration (bot token, message fetching)
3. `TASK-912`: Collector service — Reddit integration (OAuth2, subreddit fetch)
4. `TASK-913`: Collector service — App Store integration (choose Apple API vs Appbot, implement at least one for MVP)
5. `TASK-914`: Analyzer service — sentiment lexicon + scoring algorithm
6. `TASK-915`: Analyzer service — topic tagging and feature request extraction
7. `TASK-916`: Analyzer service — weekly aggregation and report generation
8. `TASK-917`: Reporter service — Telegram delivery
9. `TASK-918:** Cron job setup and secrets configuration
10. `TASK-919:** End-to-end test with mock data and dry-run report

---

## 11. Future Enhancements

- Dashboard UI (Grafana-like) for real-time sentiment trends
- Integration with kanban: auto-create tasks from high-priority feature requests
- Sentiment model training on game-specific corpus
- Multi-language support (global community)
- Alerting on sudden negative sentiment spikes
- More sophisticated topic modeling (BERT embeddings)

---

## 12. Questions / Open Issues

1. **App Store API decision:** Apple native API vs third-party (Appbot)? Recommendation: start with Appbot for simplicity and unified API across iOS/Android, then replace later if needed.
2. **Sentiment model:** Rule-based okay for MVP? Yes, but will miss nuance. Keep plugin architecture to swap in ML later.
3. **Storage:** Use existing game DB? Yes, but isolated schema prevents accidental cross-writes.
4. **Cron host:** Where will collectors run? On the same server hosting OpenClaw gateway? Yes, simplest.
5. **Data retention:** How long to keep raw posts? Suggest 90 days for `raw_posts`, indefinitely for aggregated `weekly_reports`. Add cleanup job later.

---

**Next steps:** Director review and approval. Upon approval, tasks TASK-910–TASK-919 will be created and assigned to coder/foundry.
