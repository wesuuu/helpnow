-- Marketing Metrics Schema for ClickHouse

-- Impressions
CREATE TABLE IF NOT EXISTS impressions (
    organization_id UInt64,
    campaign_id UInt64,
    count UInt64,
    meta String, -- JSON
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, campaign_id, created_at);

-- Cost Per Acquisition (CPA)
CREATE TABLE IF NOT EXISTS cpa (
    organization_id UInt64,
    campaign_id UInt64,
    cost Float64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, campaign_id, created_at);

-- Click Rate
CREATE TABLE IF NOT EXISTS click_rate (
    organization_id UInt64,
    campaign_id UInt64,
    clicks UInt64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, campaign_id, created_at);

-- Open Rate
CREATE TABLE IF NOT EXISTS open_rate (
    organization_id UInt64,
    campaign_id UInt64,
    opens UInt64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, campaign_id, created_at);
