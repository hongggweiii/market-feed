-- Create trades table
-- ClickHouse constraints NON-NULL by default
CREATE TABLE IF NOT EXISTS trades (
    event_type String,
    event_time Int64,
    symbol LowCardinality(String),
    trade_id Int64,
    price Decimal(18, 8),
    quantity Float64,
    trade_time DateTime64(3),
    is_market_maker Bool
) ENGINE = MergeTree()
ORDER BY (symbol, trade_time, trade_id);