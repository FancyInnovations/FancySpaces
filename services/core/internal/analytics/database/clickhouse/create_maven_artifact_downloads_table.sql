CREATE TABLE IF NOT EXISTS fancyspaces.maven_artifact_downloads
(
    space_id        LowCardinality(String),
    repository_name LowCardinality(String),
    group_id        LowCardinality(String),
    artifact_id     LowCardinality(String),
    version         LowCardinality(String),
    downloaded_at   DateTime CODEC (DoubleDelta, LZ4),
    ip_hash         String
)
ENGINE MergeTree
PARTITION BY toYYYYMM(downloaded_at)
ORDER BY (space_id, repository_name, downloaded_at);