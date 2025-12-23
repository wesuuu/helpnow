-- Marketing Metrics Schema for ClickHouse

-- Impressions
CREATE TABLE IF NOT EXISTS impressions (
    organization_id UInt64,
    meta String, -- JSON
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, created_at);

-- Cost Per Acquisition (CPA)
CREATE TABLE IF NOT EXISTS cpa (
    organization_id UInt64,
    cost Float64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, created_at);

-- Clicks (formerly click_rate)
CREATE TABLE IF NOT EXISTS click_rate (
    organization_id UInt64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, created_at);

-- Opens (formerly open_rate)
CREATE TABLE IF NOT EXISTS open_rate (
    organization_id UInt64,
    meta String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (organization_id, created_at);
