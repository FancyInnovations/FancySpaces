CREATE TABLE IF NOT EXISTS fancyspaces.version_downloads
(
    space_id      LowCardinality(String),
    version_id    LowCardinality(String),
    downloaded_at DateTime CODEC (DoubleDelta, LZ4),
    ip_hash       String,
    user_agent    String
)
ENGINE MergeTree
PARTITION BY toYYYYMM(downloaded_at)
ORDER BY (space_id, version_id, downloaded_at);