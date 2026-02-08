-- Create community schema
CREATE SCHEMA IF NOT EXISTS community;

-- Raw posts from external sources (Discord, Reddit, App Stores)
CREATE TABLE IF NOT EXISTS community.raw_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source VARCHAR(50) NOT NULL, -- 'discord', 'reddit', 'app_store_ios', 'app_store_android'
    external_id VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    content TEXT NOT NULL,
    posted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    collected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    UNIQUE (source, external_id)
);

-- Processed posts with sentiment and topics
CREATE TABLE IF NOT EXISTS community.posts (
    id UUID PRIMARY KEY REFERENCES community.raw_posts(id) ON DELETE CASCADE,
    sentiment_score INTEGER NOT NULL, -- -100 to 100
    sentiment_label VARCHAR(20) NOT NULL, -- 'negative', 'neutral', 'positive'
    topics TEXT[] DEFAULT '{}',
    feature_requests TEXT[] DEFAULT '{}',
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Weekly aggregated reports
CREATE TABLE IF NOT EXISTS community.weekly_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_start DATE NOT NULL UNIQUE,
    total_posts INTEGER NOT NULL,
    unique_authors INTEGER NOT NULL,
    sentiment_distribution JSONB NOT NULL, -- {positive: X, neutral: Y, negative: Z}
    top_topics JSONB NOT NULL, -- [{topic: 'balance', count: 10}, ...]
    top_feature_requests JSONB NOT NULL, -- [{request: 'more gold', count: 5}, ...]
    raw_summary TEXT, -- Human-readable summary
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_community_raw_posts_source_posted ON community.raw_posts(source, posted_at);
CREATE INDEX IF NOT EXISTS idx_community_raw_posts_collected_at ON community.raw_posts(collected_at);
CREATE INDEX IF NOT EXISTS idx_community_posts_sentiment ON community.posts(sentiment_label);
CREATE INDEX IF NOT EXISTS idx_community_posts_topics ON community.posts USING GIN (topics);
CREATE INDEX IF NOT EXISTS idx_community_posts_feature_requests ON community.posts USING GIN (feature_requests);
