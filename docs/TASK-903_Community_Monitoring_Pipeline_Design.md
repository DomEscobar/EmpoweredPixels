# TASK-903: Community Monitoring Pipeline Design

**Document Version:** 1.0  
**Prepared for:** EmpoweredPixels Game Director  
**Prepared by:** Community Analyst (subagent)  
**Date:** 2025-02-09  

---

## Executive Summary

This document outlines a robust, scalable pipeline to aggregate and analyze player sentiment and feedback from Discord, Reddit, and app store reviews (Google Play Store, Apple App Store, Steam). The system will provide real-time insights, trend detection, and automated alerting to support proactive community management and inform product decisions.

**Key Goals:**
- Centralize community feedback from all major platforms
- Provide near-real-time sentiment analysis with ≥85% accuracy
- Detect negative sentiment spikes and emerging issues within 15 minutes
- Generate actionable insights integrated with the existing kanban system
- Maintain full GDPR compliance and data privacy
- Handle platform rate limits and API failures gracefully

**Out of Scope:** Social media sentiment beyond Discord/Reddit (can be added in Phase 4)

---

## 1. Architecture Overview

### 1.1 High-Level System Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                       External Sources                              │
│  ┌─────────┐  ┌──────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │ Discord │  │  Reddit  │  │ App Stores  │  │   Other (RSS)  │  │
│  └────┬────┘  └────┬─────┘  └──────┬──────┘  └────────┬────────┘  │
│       │            │              │                   │            │
│       ▼            ▼              ▼                   ▼            │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │              Data Collection Layer (API/Webhook)           │  │
│  │  ┌──────────┐  ┌──────────┐  ┌─────────────────────┐    │  │
│  │  │ Discord  │  │ Reddit   │  │ App Store Scrapers  │    │  │
│  │  │ Collector│  │ Collector│  │ (TOS-compliant)     │    │  │
│  │  └────┬─────┘  └────┬─────┘  └──────────┬──────────┘    │  │
│  └───────┼─────────────┼──────────────────┼──────────────────┘  │
│          │             │                  │                       │
│          ▼             ▼                  ▼                       │
│  ┌─────────────────────────────────────────────────────────────┐  │
│  │              Message Queue (RabbitMQ / Kafka)              │  │
│  │               (Buffer + Rate Control)                     │  │
│  └─────────────────────────────┬───────────────────────────────┘  │
│                                │                                  │
│          ┌─────────────────────┼─────────────────────┐          │
│          ▼                     ▼                     ▼          │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐  │
│  │   Raw Data   │      │   Raw Data   │      │   Raw Data   │  │
│  │   Storage    │      │   Storage    │      │   Storage    │  │
│  │  (PostgreSQL)│      │  (PostgreSQL)│      │  (PostgreSQL)│  │
│  └──────┬───────┘      └──────┬───────┘      └──────┬───────┘  │
│         │                    │                      │           │
│         └────────────────────┼──────────────────────┘           │
│                              │                                  │
│                              ▼                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Normalization & Enrichment Service           │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │ 1. Parse platform-specific formats                  │ │  │
│  │  │ 2. Extract: author, timestamp, content, metadata    │ │  │
│  │  │ 3. Deduplicate (same user, same content, <5min)     │ │  │
│  │  │ 4. Generate normalized schema                       │ │  │
│  │  │ 5. Anonymize PII (user ID → hash, remove emails)   │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                                ▼                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Sentiment Analysis Engine                   │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │ Option A: Third-party API (e.g., MeaningCloud,      │ │  │
│  │  │          OpenAI, Google Cloud NLP)                  │ │  │
│  │  │ Option B: In-house ML (BERT-based, fine-tuned)      │ │  │
│  │  │ - Cost: High upfront, lower per-message            │ │  │
│  │  │ - Accuracy: 80-85% baseline, 90%+ after fine-tune  │ │  │
│  │  │ - Latency: 300ms - 2s                              │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  │  Output: sentiment_score (-1 to +1), confidence,        │ │  │
│  │          aspect_ratings (gameplay, UI, performance, etc)│ │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                                ▼                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Analytics & Aggregation                     │  │
│  │  • Rolling averages (5min, 1hr, 24hr)                   │  │
│  │  • Topic clustering (keywords, hashtags, mentions)      │  │
│  │  • Spike detection (anomaly detection, >2σ)             │  │
│  │  • Trend forecasting (7-day projection)                 │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                                ▼                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Alerting Engine                             │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │ Triggers:                                            │ │  │
│  │  │ • Negative sentiment spike (>30% increase in 15min)│ │  │
│  │  │ • Keyword: "crash", "bug", "broken"                │ │  │
│  │  │ • Volume spike (>5x normal 1hr rate)               │ │  │
│  │  │ • Single issue clustering (>50 mentions in 1hr)    │ │  │
│  │  └──────────────────────────────────────────────────────┘ │  │
│  │  Actions: Send to DM channels, create kanban cards,     │ │  │
│  │           escalate to director if severity high         │ │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                                ▼                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Dashboard & Reporting (React + Grafana)    │  │
│  │  • Live sentiment gauge                                │  │
│  │  • Platform breakdown with filters                     │  │
│  │  • Top keywords & trending topics                      │  │
│  │  • Historical trend charts                             │  │
│  │  • Drill-down to individual messages                   │  │
│  │  • Export: CSV, PDF, API                               │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                │                                 │
│                                ▼                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Kanban Integration                          │  │
│  │  • Auto-create cards in incubator/backlog               │  │
│  │  • Link to original source (anonymized)                │  │
│  │  • Priority scoring based on sentiment + volume         │  │
│  │  • Update existing cards with new mentions             │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Data Collection Layer

### 2.1 Discord Integration

**Strategy:** Webhook-based collection using Discord API

**Components:**
1. **Guild Message Collector** (for each server)
   - Subscribe to message events via gateway intents
   - Filter channels: `#feedback`, `#bugs`, `#suggestions`, plus configurable list
   - Rate limit: Respect Discord's 50 requests/sec limit per bot

2. **Thread Support**
   - Collect from threads within filtered channels
   - Archive unarchived threads on demand

3. **Polling Fallback**
   - If webhook fails, fallback to REST API polling every 5 minutes
   - Cache last message ID to avoid duplicates

**Technical Implementation:**
```python
# Pseudo-code
class DiscordCollector:
    def __init__(self, bot_token, guild_id, channel_ids):
        self.gateway = DiscordGateway(bot_token)
        self.rest = DiscordREST(bot_token)
        self.channel_ids = channel_ids

    async def start_webhook(self):
        intents = GatewayIntents(
            guild_messages=True,
            message_content=True,
            guilds=True
        )
        await self.gateway.connect(intents)
        self.gateway.on('MESSAGE_CREATE', self.handle_message)

    async def handle_message(self, event):
        if event.channel_id in self.channel_ids:
            message = self.normalize_discord_message(event)
            await queue.publish('raw.discord', message)

    def normalize_discord_message(self, event):
        return {
            'platform': 'discord',
            'source_id': f"discord_{event.channel_id}_{event.id}",
            'author_id': self.anonymize_id(event.author.id),
            'author_name': event.author.username,  # Keep, not PII
            'channel_id': event.channel_id,
            'channel_name': event.channel.name,
            'content': event.content,
            'attachments': [a.url for a in event.attachments],
            'embeds': len(event.embeds),
            'timestamp': event.timestamp.isoformat(),
            'guild_id': event.guild_id,
            'message_type': 'thread' if event.message_type == 'thread' else 'channel',
            'thread_id': event.thread_id if hasattr(event, 'thread_id') else None
        }
```

**Rate Limits & Retry:**
- Gateway: Unlimited (persistent connection)
- REST: 50 req/sec per bot; exponential backoff on 429
- Retry strategy: 3 attempts with jitter, then dead-letter queue

**Error Handling:**
- Gateway disconnect: auto-reconnect with 5s delay
- Invalid auth: alert and stop
- Partial failures: log and continue (skip individual channels)

### 2.2 Reddit Integration

**Strategy:** Official Reddit API with streaming via PRAW

**Components:**
1. **Subreddit Stream Collectors**
   - For each target subreddit (e.g., r/IndieGaming, r/our_game_name)
   - Stream new submissions and comments
   - Filter by flairs: "Feedback", "Bug", "Suggestion" (configurable)

2. **Search Query Collectors**
   - Catch cross-posts not in main subs
   - Search: "game name + bug/crash/feedback"
   - Execute every 30 minutes

**Technical Implementation:**
```python
import praw

class RedditCollector:
    def __init__(self, client_id, client_secret, subreddits):
        self.reddit = praw.Reddit(
            client_id=client_id,
            client_secret=client_secret,
            user_agent='EmpoweredPixels/1.0'
        )
        self.subreddits = subreddits

    async def stream_submissions(self):
        for sub in self.subreddits:
            subreddit = self.reddit.subreddit(sub)
            for submission in subreddit.stream.submissions(skip_existing=True):
                if self.is_relevant(submission):
                    data = self.normalize_reddit_submission(submission)
                    await queue.publish('raw.reddit', data)

    async def stream_comments(self):
        for sub in self.subreddits:
            subreddit = self.reddit.subreddit(sub)
            for comment in subreddit.stream.comments(skip_existing=True):
                data = self.normalize_reddit_comment(comment)
                await queue.publish('raw.reddit', data)

    def normalize_reddit_submission(self, submission):
        return {
            'platform': 'reddit',
            'source_id': f"reddit_{submission.id}",
            'author_id': self.anonymize_id(submission.author.name if submission.author else '[deleted]'),
            'author_name': submission.author.name if submission.author else '[deleted]',
            'title': submission.title,
            'content': submission.selftext,
            'url': submission.url,
            'score': submission.score,
            'upvote_ratio': submission.upvote_ratio,
            'num_comments': submission.num_comments,
            'subreddit': submission.subreddit.display_name,
            'permalink': f"https://reddit.com{submission.permalink}",
            'flair': submission.link_flair_text,
            'created_utc': datetime.fromtimestamp(submission.created_utc).isoformat(),
            'is_self': submission.is_self
        }

    def normalize_reddit_comment(self, comment):
        return {
            'platform': 'reddit',
            'source_id': f"reddit_comment_{comment.id}",
            'author_id': self.anonymize_id(comment.author.name if comment.author else '[deleted]'),
            'parent_id': comment.parent_id,
            'submission_id': comment.submission.id,
            'content': comment.body,
            'score': comment.score,
            'permalink': f"https://reddit.com{comment.permalink}",
            'created_utc': datetime.fromtimestamp(comment.created_utc).isoformat()
        }
```

**Rate Limits:**
- Auth: 60 requests/minute per OAuth client
- Streaming: Uses persistent connections (no rate limit issues)
- Search: 1 request/30s configured; cache results 24h

**TOS Compliance:**
- Reddit API requires adherence to Reddit API Terms
- Must include `User-Agent` string with contact info
- No circumvention of rate limits
- Must display source attribution in dashboard

**Error Handling:**
- 401/403: re-auth with refresh token, alert if persistent
- 429: backoff with jitter, dead-letter after 5 retries
- Connection loss: exponential reconnect (1s, 2s, 4s, 8s...)

### 2.3 App Store Reviews Integration

**Strategy:** TOS-compliant scraping + official APIs where available

**Sources:**
1. **Google Play Store** - Use `google-play-scraper` library (public API, not official)
2. **Apple App Store** - Use `app-store-scraper` library (public API)
3. **Steam Store** - Official Steam Web API for reviews

**Note:** These use public endpoints that don't require scraping private data; they are TOS-compliant as they mimic mobile app behavior.

**Technical Implementation:**
```python
from google_play_scraper import reviews, Sort
from app_store_scraper import AppStore
from steam import SteamAPI

class AppStoreCollector:
    def __init__(self, app_ids):
        self.app_ids = app_ids  # {'google': 'com.game.name', 'apple': 'game-name', 'steam': '123456'}

    async def collect_google_play(self, app_id):
        result, _ = reviews(
            app_id,
            lang='en',
            country='us',
            sort=Sort.NEWEST,
            count=100  # Max per request
        )
        for review in result:
            data = self.normalize_google_review(review, app_id)
            await queue.publish('raw.appstore', data)

    async def collect_apple_app_store(self, app_id):
        app = AppStore(country='us', app_name=app_id)
        app.review(how_many=100)
        for review in app.reviews:
            data = self.normalize_apple_review(review, app_id)
            await queue.publish('raw.appstore', data)

    async def collect_steam_reviews(self, app_id):
        steam = SteamAPI()
        reviews = steam.apps.get_user_reviews(app_id, num_per_page=100)
        for review in reviews:
            data = self.normalize_steam_review(review, app_id)
            await queue.publish('raw.appstore', data)

    def normalize_google_review(self, review, app_id):
        return {
            'platform': 'google_play',
            'source_id': f"gp_{review['reviewId']}",
            'author_id': self.anonymize_id(review['userName']),
            'author_name': review['userName'],
            'content': review['content'],
            'rating': review['score'],
            'thumbs_up': review['thumbsUpCount'],
            'timestamp': review['at'].isoformat(),
            'app_version': review['appVersion'],
            'reply': review.get('replyContent'),
            'reply_timestamp': review.get('repliedAt'),
            'app_id': app_id
        }

# Similar for Apple and Steam
```

**Rate Limiting:**
- Google Play: 5 req/sec (aggressive); 100 reviews/req; request every 10 minutes per app
- Apple: 20 req/sec; 50 reviews/req; request every 10 minutes per app
- Steam: 200 req/5min; 100 reviews/req; request every 15 minutes per app

**Error Handling:**
- 404 (app removed): mark app as inactive, alert
- 429 throttle: backoff 1 hour minimum
- JSON parse errors: log with context, skip batch
- Network errors: 3 retries with exponential backoff

### 2.4 Collection Scheduler

**Orchestration:** Celery Beat or APScheduler

**Schedule:**
- Discord: Real-time (webhook), polling fallback every 5 min if disconnected
- Reddit: Real-time streaming (persistent processes)
- App Stores: Every 10 minutes (staggered to avoid rate limit collisions)
- Health check: Every minute; restart failed collectors

**Process Management:**
- Use systemd or supervisord to keep collector processes alive
- Auto-restart on crash with exponential backoff
- Alert director if collector down >5 minutes

---

## 3. Data Normalization Schema

### 3.1 Unified Data Model

All sources mapped to canonical `Feedback` schema:

```json
{
  "id": "uuid_v4",
  "platform": "discord|reddit|google_play|apple_appstore|steam",
  "source_id": "unique_id_per_platform",
  "source_url": "link_to_original",
  "author_hash": "sha256_of_author_id",
  "author_name": "display_name (not PII)",
  "content_raw": "original_text",
  "content_clean": "stripped_markdown, normalized_unicode",
  "content_language": "en|fr|de|... (auto-detected)",
  "title": "optional (for Reddit/submissions)",
  "rating": 1-5 (where applicable, null for Discord/Reddit),
  "metadata": {
    "channel": "discord_channel_name",
    "subreddit": "subreddit_name",
    "flair": "user_flair_text",
    "app_version": "1.2.3",
    "device": "iOS/Android/PC",
    "attachments": ["url1", "url2"],
    "thumbs_up": 0,
    "score": 0
  },
  "timestamp": "2025-02-09T14:30:00Z",
  "ingested_at": "2025-02-09T14:30:05Z",
  "sentiment": {
    "score": -0.42,
    "confidence": 0.87,
    "label": "negative",
    "aspects": [
      {"aspect": "gameplay", "score": -0.6},
      {"aspect": "performance", "score": -0.8},
      {"aspect": "ui", "score": 0.2}
    ]
  },
  "topics": ["crash", "multiplayer", "level_5"],
  "is_duplicate": false,
  "duplicate_of": null
}
```

### 3.2 Normalization Pipeline

**Step 1: Parse & Extract**
- Platform-specific → dict with raw fields
- Extract author_id, content, timestamp, platform metadata

**Step 2: Content Cleaning**
- Remove markdown/HTML tags
- Normalize Unicode (NFC)
- Strip excessive whitespace
- Truncate to 2000 chars (keep full in separate field if needed)

**Step 3: Anonymization (GDPR)**
- `author_id` → SHA-256 hash with per-app salt (not reversible)
- Remove email addresses, phone numbers, IP addresses if present
- Replace real names in content if detected (NLP PII detection)

**Step 4: Deduplication**
- Bucket: same platform + author_hash + content hash (SHA-1) + ±5min window
- If exists: mark `is_duplicate=true`, store `duplicate_of=existing_id`
- Rationale: Users cross-post same feedback to multiple places

**Step 5: Language Detection**
- Use `langdetect` or `fasttext` library
- Non-English: store `content_language` and process separately
- Optional: translate to English for sentiment (third-party API cost)

**Step 6: Enrichment Placeholders**
- `sentiment`: initially null, filled by analysis engine
- `topics`: initially null, filled by topic clustering

**Storage:**
- PostgreSQL table `feedback` with indexes on:
  - `(platform, source_id)` - unique constraint
  - `timestamp` (for time queries)
  - `sentiment_score` (for sorting)
  - `author_hash` (for aggregating user contributions)
  - `is_duplicate` (filtering)

---

## 4. Sentiment Analysis Approach

### 4.1 Recommendation: Hybrid Strategy (Third-Party API + Custom Fine-Tuning)

**Base Model:** Start with third-party API for rapid deployment (Phase 1) → transition to in-house fine-tuned model for cost optimization (Phase 2-3).

**Third-Party Options Evaluation:**

| Provider | Cost/1K msgs | Accuracy | Latency | GDPR? | Recommendation |
|----------|--------------|----------|---------|-------|----------------|
| OpenAI GPT-4o-mini | ~$0.10 | 90% (zero-shot) | 500ms | No (US) | ❌ High cost, US data |
| Google Cloud NLP | ~$0.02 | 85% | 300ms | Yes (EU) | ✅ Good balance |
| MeaningCloud | ~$0.015 | 82% | 200ms | Yes (EU) | ✅ Best cost/accuracy |
| AWS Comprehend | ~$0.025 | 84% | 400ms | Yes (EU) | ✅ Solid option |
| In-house BERT | $0 (infrastructure) | 80% baseline → 90% fine-tuned | 100ms | ✅ Full control | ⏳ 6-8 weeks dev |

**Selected Path:**
- **Phase 1** (Weeks 1-4): MeaningCloud API (quick start, EU-hosted)
- **Phase 2** (Weeks 5-8): Begin parallel training on labeled dataset
- **Phase 3** (Weeks 9-12): Switch to in-house model, keep API as fallback

### 4.2 Sentiment Analysis Implementation

**A. Third-Party API (Phase 1)**

```python
import requests
from meaningcloud import SentimentRequest

class SentimentAnalyzer:
    def __init__(self, api_key, model='meaningcloud'):
        self.api_key = api_key
        self.model = model
        if model == 'meaningcloud':
            self.endpoint = 'https://api.meaningcloud.com/sentiment-2.1'
        elif model == 'google':
            self.endpoint = 'https://language.googleapis.com/v1/documents:analyzeSentiment'

    async def analyze(self, text: str) -> dict:
        if self.model == 'meaningcloud':
            resp = requests.post(self.endpoint, data={
                'key': self.api_key,
                'txt': text,
                'lang': 'en',
                'model': 'general'
            })
            result = resp.json()
            return {
                'score': result['score_tag'] / 5,  # Convert 0-5 to -1 to 1
                'confidence': result['confidence'],
                'label': result['score_tag'],  # positive, negative, neutral
                'aspects': self.extract_meaningcloud_aspects(result)
            }
```

**Aspect-Based Sentiment:**
- Extract aspects: gameplay, UI, performance, economy, multiplayer, controls
- Use dependency parsing or aspect extraction from API
- Store per-aspect scores to identify pain points

**Handling Non-English:**
- Detect language first
- If supported (English, French, German, Spanish, Portuguese): use API
- If unsupported: default to neutral (0) with low confidence; flag for review

**Cost Estimate (Phase 1):**
- 500K messages/month (estimated)
- $0.015/1K = $7.50/month (very affordable)
- Budget: $50/month max including development

**B. In-House Model (Phase 2-3)**

**Dataset:**
- Label 5,000 samples from existing data (human-labeled)
- Use weak supervision: expand via API+sentiment as pseudo-labels
- Target classes: `positive`, `neutral`, `negative` + aspects

**Model Architecture:**
- Base: `distilbert-base-uncased` (lightweight, fast)
- Fine-tune on EmpoweredPixels feedback corpus
- Aspect extraction: multi-task learning (sentiment + 6 aspect heads)

**Training Pipeline:**
```
Raw feedback (unlabeled)
    ↓
Label 5K (human) + 100K (API pseudo-label)
    ↓
Train/val split 80/20
    ↓
Fine-tune DistilBERT for classification
    ↓
Evaluate: target 85% accuracy, 90% after iteration
    ↓
Export ONNX for fast inference
    ↓
Deploy to Triton Inference Server or FastAPI
```

**Hardware:**
- Training: Google Colab Pro ($10/month) or AWS g4dn.xlarge ($0.526/hr)
- Inference: Small VM (2 vCPU, 4GB RAM) sufficient; batch 32 msgs

**Migration Plan:**
- Week 5: Begin labeling sprint (director + player agents)
- Week 6: Train baseline model, compare to API
- Week 7: A/B test: 50% API, 50% in-house; compare quality
- Week 8: Improve model based on errors, repeat A/B
- Week 9: Switch 100% to in-house; keep API as fallback if quality drops < 85%

---

## 5. Alerting Mechanisms

### 5.1 Alert Triggers

**Alert Types:**

| Trigger | Condition | Severity | Action |
|---------|-----------|----------|--------|
| **Negative Spike** | Neg sentiment > 60% of total in 15min AND volume > 2x baseline | HIGH | DM to @Community Manager, create kanban card |
| **Volume Spike** | Total messages > 5x baseline in 15min | MEDIUM | DM to @Analyst, dashboard highlight |
| **Keyword Alert** | Keywords: "crash", "bug", "broken", "down" in ≥20 msgs/15min | HIGH | Immediate DM to @Releaser & director |
| **Single Issue Cluster** | Same keyword phrase (e.g., "level 5 boss") in ≥50 msgs/1hr | HIGH | Create aggregated card, link all sources |
| **Rating Drop (App Stores)** | Average rating drops >0.5 stars in 24hr | HIGH | DM to @Director, investigate |
| **Critical Bug Report** | Keywords + "can't play", "unplayable", "stuck" | CRITICAL | Page director, create hotfix priority card |
| **Positive Trend** | Pos sentiment > 80% for 24hr (milestone) | LOW | Weekly report highlight (no immediate alert) |

### 5.2 Baseline Calculation

**Baseline Metrics (calculated per platform, rolling 24hr):**
- Average messages per hour
- Average sentiment score
- Average rating (app stores)
- Update baseline every hour (exponential smoothing)

**Example:**
```
baseline_volume = 0.7 * last_baseline + 0.3 * last_24hr_avg
if current_15min_volume > baseline_volume * 5:
    trigger_volume_spike()
```

### 5.3 Alert Delivery Channels

1. **Direct Messages (Primary)**
   - Discord: Send to `#alerts` channel (webhook from backend)
   - Telegram: To director's Telegram (via `message` tool)
   - Slack/Teams: Optional integration

2. **Kanban Auto-Creation**
   - Create card in appropriate column:
     - `incubator` for strategic issues
     - `backlog` for actionable bugs/features
   - Card content:
     - Title: "🚨 [ALERT] Negative sentiment spike on Discord"
     - Description: Summary, stats, top keywords, links (anonymized)
     - Labels: `community`, `urgent`, `bug`, `crash`, etc.
     - Assign to: `Community Analyst` initially

3. **Dashboard Banner**
   - Show red banner across top: "⚠️ Active Alert: Negative sentiment spike detected"
   - Click → opens alert details panel

4. **Escalation Policy**
   - Level 1 (MEDIUM): Notify Community Manager
   - Level 2 (HIGH): Notify Director + Releaser
   - Level 3 (CRITICAL): Page Director via Telegram/phone (if configured)

### 5.4 Alert Throttling & Snoozing

- Cooldown: Same alert type silenced for 2 hours after trigger (unless severity increases)
- Auto-snooze: If sentiment returns to normal for 1 hour, deactivate alert
- Manual snooze: Director can snooze specific alert for custom duration
- Alert history: Stored in `alerts` table for audit

### 5.5 False Positive Prevention

- Require minimum volume before triggering (e.g., >50 msgs/15min)
- Wait 5 minutes, re-evaluate before sending (avoid transient spikes)
- If same alert repeats within 24hr, escalate severity one level

---

## 6. Dashboard & Reporting UI Concepts

### 6.1 Dashboard Layout

**Technology Stack:**
- Frontend: React + TypeScript + Tailwind CSS
- Charts: Recharts or Chart.js
- State: React Query (caching)
- Real-time: WebSocket (Socket.io) or Server-Sent Events

**Layout:**

```
┌─────────────────────────────────────────────────────────────────────┐
│ EMPOWEREDPIXELS COMMUNITY MONITORING              [Last: 2m ago] │
├─────────────┬───────────────┬───────────────┬─────────────────────┤
│             │ Sentiment Gauge│  Volume Spark│   Alerts Panel     │
│             │  [Gauge: 68%] │  [+42% ▲]    │  ⚠️ 2 Active       │
├─────────────┴───────────────┴───────────────┴─────────────────────┤
│ Platform Filters: [✓ Discord] [✓ Reddit] [✓ Google Play]         │
│ Time Range: [1H] [6H] [24H] [7D]                                 │
├─────────────────────────────────────────────────────────────────────┤
│ Top Keywords (Last 24h)         │ Top Topics (Last 24h)          │
│ ┌────────────────────────────┐ │ ┌────────────────────────────┐ │
│ │ crash (142) ▲13            │ │ │ multiplayer (89)          │ │
│ │ bug (98) ▲7                │ │ │ level 5 (67) ▲12         │ │
│ │ server (76) ▼2             │ │ │ economy (54) ▼5          │ │
│ │ update (65)                │ │ │ controls (43)            │ │
│ │ lag (52) ▲9                │ │ │ daily rewards (38) ▲8   │ │
│ └────────────────────────────┘ │ └────────────────────────────┘ │
├─────────────────────────────────────────────────────────────────────┤
│ Sentiment Trend (24h)                                              │
│ [Line chart: score -0.3 → +0.1 (improving)]                       │
├─────────────────────────────────────────────────────────────────────┤
│ Recent Feedback (Live)               [Refresh: 5s]                │
│ ┌───────────────────────────────────────────────────────────────┐ │
│ │ ⚠️ Discord #bugs - u/Player123 (10m ago)                     │ │
│ │    "Game crashes when entering dungeon at level 5"           │ │
│ │    Sentiment: 🔴 Negative | aspects: [crash, gameplay]      │ │
│ │    [View Source] [Create Card]                               │ │
│ ├───────────────────────────────────────────────────────────────┤ │
│ │ ✅ Reddit r/IndieGaming - u/CasualGamer (25m ago)           │ │
│ │    "Loving the new update! Much smoother now."              │ │
│ │    Sentiment: 🟢 Positive | aspects: [performance]          │ │
│ └───────────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

### 6.2 Key Views

**1. Overview Dashboard (Home)**
- Sentiment score gauge (0-100) with trend indicator
- Total volume last 15min vs. baseline
- Alert banner if active
- Top keywords cloud (size = frequency, color = sentiment)
- Top topics list
- Platform breakdown pie chart
- Recent feedback feed (live-updating every 5s)

**2. Platform Drill-Down**
- Click platform → full view filtered to that source
- Tabs: Messages, Ratings (for app stores), Trends
- Channel/Subreddit selector (Discord/Reddit)

**3. Topics Analysis**
- Barchart of topics with sentiment scores
- Click topic → shows all messages related
- Topic trend over time

**4. Historical Reporting**
- Date range picker (default: last 7 days)
- Export: CSV, PDF (weekly summary)
- Compare periods: "This week vs. last week"

**5. Alert Details Modal**
- Trigger reason, affected platforms, severity
- Sample messages (preview)
- Actions: Snooze, Mark Resolved, Create Card
- Timeline of alert history

### 6.3 Interactive Features

**Live Feed Table:**
- Columns: Platform, Channel, Author (anonymized), Content (truncated), Sentiment (color-coded), Timestamp
- Sort by: Sentiment score, timestamp
- Filter: Sentiment range, platform, keyword
- Pagination: Infinite scroll

**Message Actions:**
- `[Create Card]` button: Opens kanban card form pre-filled with:
  - Title: Auto-summarized (first 100 chars)
  - Description: Full message + source link
  - Labels: Based on sentiment + platform + topics
  - Priority: Based on sentiment score (negative = higher)

**Real-Time Updates:**
- WebSocket pushes new messages as they arrive
- Semantic updates: sentiment score changes, new keyword ranking
- Use React Query for optimistic updates

**Performance:**
- Load initial page: <2s
- Live feed: <100ms latency for new messages
- Support 10-20 concurrent analysts

### 6.4 Mobile-Responsive Design

- Dashboard: Responsive grid (CSS Grid)
- Collapsible side panels
- Touch-friendly controls

---

## 7. GDPR Compliance

### 7.1 Data Privacy Principles

**Lawful Basis:** Legitimate interest (monitoring public feedback for product improvement)

**Key Requirements:**
1. **Data Minimization:** Only collect necessary fields
2. **Purpose Limitation:** Use solely for community sentiment analysis
3. **Storage Limitation:** Retain raw data 90 days, aggregated metrics 2 years
4. **Accuracy:** Allow correction requests (remove/rectify)
5. **Transparency:** Document in privacy policy + internal logs
6. **Security:** Encryption at rest, access controls
7. **Rights:** Support erasure, access, portability requests

### 7.2 Implementation

**Anonymization (At Ingestion):**
```python
def anonymize_id(raw_id: str) -> str:
    """One-way hash with app-specific salt (not reversible)"""
    salt = os.getenv('ANON_SALT', 'default-salt')  # Different salt per app if needed
    return hashlib.sha256(f"{raw_id}:{salt}".encode()).hexdigest()[:16]

def remove_pii(text: str) -> str:
    """Remove emails, phone numbers, IPs using regex + NLP"""
    # Regex patterns for email, phone
    text = re.sub(r'\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b', '[EMAIL]', text)
    text = re.sub(r'\b\d{3}[-.]?\d{3}[-.]?\d{4}\b', '[PHONE]', text)
    # spaCy NER for PERSON, DATE, etc if needed
    return text
```

**Data Retention Policy:**
- Raw feedback: 90 days (delete after)
- Normalized with sentiment: 2 years (aggregated insights)
- Aggregated metrics (hourly/daily): Indefinite (no PII)
- Automated cleanup: PostgreSQL `pg_cron` job daily

**Encryption:**
- PostgreSQL: disk encryption (LUKS or cloud provider)
- Sensitive fields (author_hash): optional column-level encryption (pgcrypto)
- In transit: TLS 1.3 for all API/db connections

**Access Control:**
- Database: Read-only for analysts; write only for collector services
- Dashboard: SSO integration (if available) or basic auth + IP whitelist
- Audit log: Record all data access (who queried what and when)

**Right to Erasure Workflow:**
1. User requests deletion via privacy email
2. Verify identity (if possible, limited for anonymous posts)
3. Search `author_hash` across raw and normalized tables
4. Delete or pseudonymize (replace with `[REDACTED]`)
5. Log deletion request and action taken
6. Respond to user within 30 days

**Data Processing Record (GDPR Article 30):**
- Maintain log: purpose, categories of data, retention period, security measures
- Stored in `/docs/GDPR_PROCESSING_RECORD.md`

**Data Breach Response:**
- Detect: Monitor for unauthorized access alerts
- Notify: Within 72 hours to DPA if personal data breach
- Contain: Rotate credentials, patch vulnerabilities

**Transfer Restrictions:**
- Host infrastructure in EU (data residency)
- If using third-party sentiment API: ensure EU data centers (MeaningCloud, Google Cloud EU)
- No transfers to US without SCCs (avoid OpenAI US endpoints)

**Privacy Policy Addition:**
Add section: "Community Feedback Monitoring"
- Explain data sources (Discord, Reddit, app stores)
- State anonymization process
- Provide contact for privacy concerns

---

## 8. Rate Limiting & Error Handling

### 8.1 Distributed Rate Limiting

**Problem:** Multiple collector instances need to share quotas

**Solution:** Token bucket via Redis

```python
import aioredis
from fastapi import HTTPException

class RateLimiter:
    def __init__(self, redis_url: str):
        self.redis = aioredis.from_url(redis_url)

    async def consume(self, key: str, limit: int, period: int) -> bool:
        """
        Token bucket algorithm
        key: e.g., 'reddit:submissions'
        limit: max tokens (requests)
        period: refill period in seconds
        """
        current = await self.redis.get(key)
        if current is None:
            await self.redis.set(key, limit - 1, ex=period)
            return True
        if int(current) > 0:
            await self.redis.decr(key)
            return True
        return False

# Usage
limiter = RateLimiter('redis://localhost:6379')
async def collect_reddit_submissions():
    if not await limiter.consume('reddit:submissions', 60, 60):
        await asyncio.sleep(1)  # Wait for token
        return await collect_reddit_submissions()
    # proceed with API call
```

**Per-Platform Limits (configurable):**
- Discord: 50 req/sec (REST), streaming unlimited
- Reddit: 60 req/min via API, streaming unlimited
- Google Play: 5 req/sec → limit to 4
- Apple: 20 req/sec → limit to 15

### 8.2 Retry Strategy

**Exponential Backoff with Jitter:**

```python
async def with_retry(func, max_retries=3, base_delay=1.0):
    for attempt in range(max_retries):
        try:
            return await func()
        except (RateLimitError, NetworkError) as e:
            if attempt == max_retries - 1:
                raise
            delay = base_delay * (2 ** attempt) + random.uniform(0, 0.1)
            await asyncio.sleep(delay)
```

**Retryable Errors:**
- 429 Too Many Requests
- 5xx server errors
- Network timeouts (connect/read)
- 502/503/504 (upstream issues)

**Non-Retryable:**
- 4xx (except 429): 400 bad request, 401 unauthorized, 403 forbidden
- Authentication failures: alert and stop collector

### 8.3 Error Handling Pipeline

**Critical Failures (stop collector):**
- Invalid API credentials (401)
- Database connection lost (cannot reconnect after 3 attempts)
- Disk full / OOM

**Non-Critical (log + skip):**
- Single message parse error (log, continue to next)
- Temporary network hiccup (retry)
- Third-party API rate limit hit (throttle and retry later)

**Dead-Letter Queue (DLQ):**
- Messages that fail processing after 3 retries → store in DLQ (S3 or PostgreSQL)
- Include: raw_message, error, timestamp, retry_count
- Daily manual review or automated re-processing script

**Circuit Breaker:**
- If platform API fails >50% of requests over 5min → open circuit
- Stop attempting for 15min, alert director
- After 15min, try single request to test if recovered

**Monitoring:**
- Metrics: `collector_errors_total{platform,error_type}`
- Dashboard panel: Collector health status (green/yellow/red)
- Email/Slack alerts on collector failure >5min

---

## 9. Integration with Existing Kanban/Feedback System

### 9.1 Kanban Board Structure (from AGENTS.md)

**Columns:** incubator → backlog → in_progress → review → done

**Agent Ecosystem:**
- **Community Analyst (this agent)** → monitors external, publishes trends to `community` channel
- **Director** → reads all agent messages, assigns tasks
- **Feature Forge (coder)** → implements features from `incubator`/`backlog`
- **Player Agents** → validate features, provide feedback loops

### 9.2 Integration Mechanisms

**A. Message Bus Publication**

```python
import json
from datetime import datetime

class KanbanPublisher:
    def __init__(self, channel: str):
        self.channel = channel  # e.g., 'community', 'feedback'

    async def publish_trend_report(self):
        report = {
            'from': 'analyst',
            'to': 'director',
            'type': 'trend_report',
            'priority': 'medium',
            'timestamp': datetime.utcnow().isoformat(),
            'content': {
                'period': 'weekly',
                'sentiment_score_change': '+0.12',
                'top_requests': [
                    {'feature': 'multiplayer matchmaking', 'mentions': 142, 'sentiment': 'negative'},
                    {'feature': 'daily rewards scaling', 'mentions': 89, 'sentiment': 'neutral'}
                ],
                'negative_spikes': [
                    {'trigger': 'crash on level 5', 'mentions': 67, 'platforms': ['discord', 'reddit']}
                ],
                'recommendations': [
                    'Prioritize crash fix in hotpatch',
                    'Consider scaling daily rewards in next update'
                ]
            }
        }
        await self.channel.publish(json.dumps(report))

    async def publish_urgent_card_request(self, alert: dict):
        card = {
            'from': 'analyst',
            'to': 'director',
            'type': 'card_request',
            'priority': 'high',
            'content': {
                'title': alert['title'],
                'description': alert['summary'],
                'column': 'incubator',  # or 'backlog' if actionable now
                'labels': ['community', alert['severity']],
                'source_messages': alert['sample_message_ids'],  # IDs from feedback table
                'stats': {
                    'mention_count': alert['volume'],
                    'sentiment_score': alert['avg_sentiment'],
                    'platforms': alert['platforms']
                }
            }
        }
        await self.channel.publish(json.dumps(card))
```

**B. Direct Kanban API (if available)**

If EmpoweredPixels has a GitHub Projects or Jira API:

```python
import requests

class GitHubKanban:
    def __init__(self, token, repo, project_id):
        self.headers = {'Authorization': f'Bearer {token}'}
        self.project_id = project_id
        self.base_url = f'https://api.github.com/repos/{repo}/projects'

    async def create_card(self, column_name: str, card: dict):
        # 1. Get project columns
        project = requests.get(f'{self.base_url}/{self.project_id}', headers=self.headers).json()
        columns = {c['name']: c['id'] for c in project['columns']}
        column_id = columns[column_name]

        # 2. Create card
        resp = requests.post(
            f'{self.base_url}/columns/{column_id}/cards',
            headers=self.headers,
            json={'content': card['description'], 'text': card['title']}
        )
        return resp.json()['id']
```

**C. Weekly Summary Report**

Director receives weekly digest via `message` tool or Telegram:

```
📊 **COMMUNITY INSIGHTS WEEKLY** (2025-W06)

**Sentiment Score:** 0.42 (↑0.12 vs last week) 🟢

**Top 3 Requests:**
1. Multiplayer matchmaking (142 mentions, 🔴 negative)
2. Daily rewards scaling (89 mentions, 🟡 neutral)
3. Tutorial improvements (67 mentions, 🟢 positive)

**Alert Summary:**
• 3 negative spikes detected (all resolved)
• 1 critical crash issue (#1234) - patched

**Action Items:**
→ Prioritize matchmaking QoLA in sprint
→ Consider daily rewards rework for v1.5

📈 Full dashboard: [link]
```

---

## 10. Implementation Phases

### Phase 1: Foundation & Rapid Collection (Weeks 1-4)

**Goals:**
- Build data collection for Discord and Reddit (app stores Week 3-4)
- Set up message queue and PostgreSQL storage
- Implement basic normalization pipeline
- Deploy third-party sentiment API (MeaningCloud)
- MVP dashboard with live feed

**Deliverables:**
1. **Collector Services** (Discord, Reddit)
   - Deployable Docker containers
   - Systemd units for auto-start
   - Health check endpoints

2. **Message Queue** (RabbitMQ or Redis Streams)
   - Exchanges: `raw.discord`, `raw.reddit`, `raw.appstore`
   - Dead-letter queues configured

3. **Storage Layer**
   - PostgreSQL schema: `feedback` table + indexes
   - Connection pool (PgBouncer if needed)

4. **Normalization Service**
   - Consume from raw queues
   - Apply cleaning + anonymization
   - Write to normalized `feedback` table
   - Call sentiment API and enrich

5. **Sentiment API Integration**
   - MeaningCloud account + API key
   - Rate-limited wrapper
   - Fallback to neutral if API fails

6. **Dashboard MVP** (React + FastAPI backend)
   - Live feedback feed (last 100 messages)
   - Sentiment gauge (simple pie chart)
   - Platform filter
   - Auto-refresh every 30s

**Success Criteria:**
- Collecting ≥90% of messages from test Discord/Reddit
- Ingestion latency <60s from source to dashboard
- Sentiment analysis running with >80% non-error rate
- Dashboard accessible to director and analyst

---

### Phase 2: Alerting & Kanban Integration (Weeks 5-8)

**Goals:**
- Implement robust alerting system with thresholds
- Integrate with kanban auto-card creation
- Add app store collection (all three platforms)
- Begin in-house sentiment model training
- Enhance dashboard with charts and analytics

**Deliverables:**
1. **Alerting Engine**
   - Baseline calculation (rolling averages)
   - Trigger detection logic
   - DM integration (Discord + Telegram)
   - Alert history storage

2. **Kanban Integration**
   - GitHub Projects API or manual card templates
   - Auto-create cards from alerts
   - Priority scoring algorithm

3. **App Store Collectors**
   - Google Play, Apple App Store, Steam
   - Respect rate limits, stagger requests
   - Normalize ratings and reviews

4. **Sentiment Model Training** (parallel track)
   - Label 5,000 samples (director + player agents assist)
   - Train DistilBERT model
   - Validate against API accuracy

5. **Dashboard Enhancements**
   - Top keywords and topics panels
   - Sentiment trend charts (24h, 7d)
   - Volume sparkline
   - Alert panel with dismiss
   - Export functionality (CSV)

6. **Data Retention & GDPR**
   - Implement cleanup jobs (90-day retention)
   - Anonymization verification
   - Access logging
   - Privacy policy update

**Success Criteria:**
- Alerts trigger within 15min of spike
- Kanban cards auto-created with correct data
- All 3 app stores collecting successfully
- In-house model achieves ≥80% accuracy on holdout set
- Dashboard fully featured and used by analyst daily

---

### Phase 3: Scaling & In-House Transition (Weeks 9-12)

**Goals:**
- Switch to in-house sentiment model (with API fallback)
- Optimize performance and reliability
- Add advanced analytics: topic clustering, forecasting
- Mobile-responsive dashboard
- Complete GDPR compliance documentation

**Deliverables:**
1. **In-House Model Deployment**
   - Serve via FastAPI + ONNX Runtime
   - A/B test 50/50 with API
   - Switch 100% after validation
   - Monitor latency <200ms

2. **Advanced Analytics**
   - Topic clustering (k-means on embeddings or BERTopic)
   - Trend forecasting (ARIMA or Prophet)
   - Aspect-based sentiment reporting

3. **Dashboard**
   - Mobile-responsive layout
   - Drill-down views (platform-specific, topic-specific)
   - Historical comparisons
   - User preferences (saved filters)

4. **Performance Tuning**
   - Collector scaling: 1 instance per platform with autoscaling based on queue depth
   - Database connection pooling optimization
   - Caching layer (Redis) for frequent queries

5. **Documentation & Compliance**
   - Complete GDPR processing record
   - Data breach response plan
   - Auditor access logs for 1 year
   - Privacy notice published

**Success Criteria:**
- In-house model operates with >85% accuracy, <200ms latency
- System processes 100K messages/day without backlog
- Dashboard used by director weekly for decision-making
- GDPR compliance verified (internal audit)

---

### Phase 4: Optimization & Expansion (Weeks 13-16)

**Goals:**
- Add new community sources (Twitter/X, YouTube comments, Steam forums)
- Implement predictive alerting (forecast next 24h sentiment)
- Fine-tune alert thresholds based on false positive analysis
- Multi-language support expansion

**Deliverables:**
1. **Additional Sources** (stretch)
   - Twitter API v2 (filter: mentions of game)
   - YouTube comments (YouTube Data API)
   - Steam community discussions

2. **Predictive Features**
   - 24-hour sentiment forecast model
   - Early warning: detect trends before spike (rising slope detection)
   - Volume prediction (seasonal patterns)

3. **Fine-Tuning**
   - Review 3 months of alerts: reduce false positives by 50%
   - Adjust thresholds per platform (Discord vs. app stores differ)
   - Provide user-configurable alert rules (analyst dashboard)

4. **Multi-Language Expansion**
   - Detect languages: en, fr, de, es, pt, ja, zh
   - Train or acquire models for non-English sentiment
   - Dashboard language switcher

5. **API for External Tools**
   - REST API: `/api/v1/sentiment?platform=discord&since=2025-02-01`
   - Rate-limited (100 req/min)
   - Documentation (OpenAPI/Swagger)

**Success Criteria:**
- System stable with <1% error rate
- Alert false positive rate <20%
- Director relying on dashboard for weekly sync
- API consumed by at least one external tool (e.g., Data Studio)

---

## 11. Risk Assessment & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Platform API changes** | Medium | High | Monitor API changelogs; version pinning; fallback scrapers |
| **Rate limit exceeded** | High | Medium | Token bucket + backoff; stagger requests; increase quotas |
| **Sentiment analysis inaccuracy** | Medium | High | Human review queue; adjustable thresholds; fallback to API |
| **GDPR violation** | Low | Critical | Anonymize at ingress; legal review; data residency |
| **Scalability bottleneck** | Medium | Medium | Load testing at 2x expected volume; autoscaling; monitoring |
| **Collector downtime** | Medium | Medium | Health checks + auto-restart; dead-man alerts |
| **Data loss (DB corruption)** | Low | High | Daily backups; point-in-time recovery; read replicas |
| **False alert fatigue** | High | Medium | Cooldown periods; tuning phase 3; manual adjustment |
| **Missing data (app store scraping blocked)** | Medium | Medium | TOS-compliant methods; rotate user-agents; alternative sources |
| **Integration failure (kanban broken)** | Low | Medium | Fallback: manual card templates; monitor API health |

---

## 12. Cost Estimate

### Infrastructure (Monthly, EU region)

| Component | Spec | Cost (€) |
|-----------|------|----------|
| **Database** | PostgreSQL 2 vCPU, 8GB, 200GB | €40 |
| **Message Queue** | RabbitMQ (managed) or Redis | €25 |
| **Collector VMs** | 3 × 2 vCPU, 4GB (Discord, Reddit, App Stores) | €90 |
| **Analytics VM** | 4 vCPU, 16GB (normalization + sentiment) | €80 |
| **Dashboard VM** | 2 vCPU, 4GB (React + FastAPI) | €40 |
| **Storage** | 500GB block storage (backups) | €20 |
| **Sentiment API** | MeaningCloud (500K msgs/mo) | €8 |
| **Monitoring** | Grafana Cloud (free tier) | €0 |
| **Backup DR** | Cross-region snapshot | €15 |
| **Total** | | **€318/month** |

### Development Costs (Time)

| Phase | Tasks | Engineer-Weeks |
|-------|-------|---------------|
| Phase 1 | Collectors, DB, normalization, dashboard MVP | 4 |
| Phase 2 | Alerting, kanban, app stores, model training start | 4 |
| Phase 3 | Model deployment, scaling, GDPR docs | 3 |
| Phase 4 | New sources, forecasting, API | 2 |
| **Total** | | **13 weeks** (1 person) or parallelized to 6-8 weeks with 2 engineers |

*Note: Phase overlap possible; director review gates at end of each phase.*

---

## 13. Success Metrics (KPIs)

**Technical KPIs:**
- Data completeness: ≥95% of public feedback captured from monitored sources
- Ingestion latency: P95 < 60s from source → dashboard
- System uptime: ≥99.5% (excluding scheduled maintenance)
- Sentiment accuracy: ≥85% (validated against human labels monthly)

**Product KPIs:**
- Alert precision: ≥80% (avoiding false alarms)
- Time-to-detect critical issue: ≤30 minutes from first mention
- Kanban conversion: ≥70% of high-priority alerts result in cards created
- Analyst adoption: Director checks dashboard ≥3x/week

**Business KPIs:**
- Reduction in unresolved critical issues (by count)
- Player satisfaction score (CSAT) improvement trend (correlated)
- Response time to community-reported bugs (reduce by 50%)

---

## 14. Next Steps for Director Review

1. **Approve architecture** → Proceed to Phase 1 implementation
2. **Provide API credentials:**
   - Discord bot token and guild/channel IDs to monitor
   - Reddit client ID/secret and target subreddits
   - App store app IDs (Google, Apple, Steam)
   - MeaningCloud API key (or alternative choice)
3. **Kanban system details:**
   - GitHub Projects, Jira, or other? Provide access token/project ID
4. **Privacy policy review:**
   - Legal team to review GDPR section
5. **Resource allocation:**
   - Assign 1-2 engineers (backend + frontend) for Phase 1
   - Budget approval: €350/month infra + €50 sentiment API buffer
6. **Success criteria agreement:**
   - Define acceptable thresholds for each KPI
   - Review cadence: weekly sprint demos, monthly director checkpoint

---

## Appendix A: Entity Relationship Diagram (Simplified)

```
feedback
├─ id (UUID PK)
├─ platform
├─ source_id (unique per platform)
├─ author_hash
├─ content_raw, content_clean
├─ metadata (JSONB)
├─ timestamp
├─ ingested_at
├─ sentiment_score (-1 to 1)
├─ sentiment_confidence
├─ aspects (JSONB)
├─ topics (text[] array)
└─ (indexes: platform+source_id, timestamp, author_hash)

alerts
├─ id (UUID PK)
├─ alert_type
├─ severity
├─ triggered_at
├─ resolved_at
├─ sample_message_ids (UUID[])

kanban_cards
├─ id (UUID PK)
├─ external_id (GitHub/Jira card ID)
├─ feedback_source_ids (UUID[] linked)
└─ status

daily_aggregates (materialized view or table)
├─ date
├─ platform
├─ message_count
├─ avg_sentiment
├─ top_keywords (JSONB)
└─ top_topics (JSONB)
```

---

## Appendix B: Configuration File Structure

`config/collectors.yaml`:
```yaml
discord:
  bot_token: "${DISCORD_BOT_TOKEN}"
  guilds:
    - id: "123456789"
      channels: ["feedback", "bugs", "suggestions"]
  reconnect_delay: 5

reddit:
  client_id: "${REDDIT_CLIENT_ID}"
  client_secret: "${REDDIT_CLIENT_SECRET}"
  subreddits: ["IndieGaming", "OurGameName", "gaming"]
  search_queries: ["Our Game crash", "Our Game feedback"]
  polling_interval: 300

appstores:
  google_play:
    apps: ["com.ourgame.name"]
    interval: 600
  apple_appstore:
    apps: ["Our-Game-Name"]
    interval: 600
  steam:
    apps: [123456]
    interval: 900

queue:
  type: rabbitmq
  url: "${RABBITMQ_URL}"
  raw_exchange: "community.raw"
  dlq: "community.dlq"

database:
  url: "${POSTGRES_URL}"
  pool_size: 10

sentiment:
  provider: "meaningcloud"  # or "google", "inhouse"
  api_key: "${SENTIMENT_API_KEY}"
  inhouse_url: "http://localhost:8000/predict"  # if provider=inhouse

alerts:
  slack_webhook: "${SLACK_WEBHOOK}"
  telegram_token: "${TELEGRAM_TOKEN}"
  telegram_chat_id: "${TELEGRAM_CHAT_ID}"
  thresholds:
    negative_spike_sentiment: 0.6
    negative_spike_volume_multiplier: 2.0
    volume_spike_multiplier: 5.0
    keyword_mentions: 20

kanban:
  type: "github"
  github_token: "${GITHUB_TOKEN}"
  repo: "org/repo"
  project_id: 1
```

---

## Appendix C: Monitoring & Observability

**Metrics (Prometheus-compatible):**
```
# Collector metrics
collector_messages_total{platform="discord"}
collector_errors_total{platform="reddit", error="rate_limit"}
collector_latency_seconds{platform="google_play"}

# Queue depth
queue_messages{exchange="raw.discord", state="ready"}

# Sentiment analysis
sentiment_predictions_total
sentiment_errors_total
sentiment_latency_seconds

# Database
db_connections_active
db_queries_total{type="insert"}

# Alerts
alerts_triggered_total{severity="high", type="negative_spike"}
alerts_active

# Dashboard
dashboard_page_views_total
dashboard_api_latency_seconds
```

**Logging (structured JSON):**
```json
{
  "timestamp": "2025-02-09T14:30:00Z",
  "level": "INFO",
  "service": "discord-collector",
  "message": "Collected 23 messages from channel #feedback",
  "channel_id": "123456",
  "message_count": 23,
  "ingestion_latency_ms": 1205
}
```

**Alerting (on monitoring system):**
- Collector down >5min → page on-call
- Queue depth >10K → scale up collectors
- Sentiment error rate >5% → switch to fallback
- DB connections >80% pool → scale DB

---

**Document End**
