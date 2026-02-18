CREATE TABLE user_positions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(4) NOT NULL, -- 'long'/'short'
    size DECIMAL(20,8) NOT NULL,
    entry_price DECIMAL(20,8) NOT NULL,
    leverage INT NOT NULL DEFAULT 1,
    margin DECIMAL(20,8) NOT NULL,
    unrealized_pnl DECIMAL(20,8) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_user_symbol (user_id, symbol)
);

CREATE TABLE funding_rates (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    rate DECIMAL(10,8) NOT NULL, -- e.g., 0.0001 (0.01%)
    timestamp TIMESTAMP DEFAULT NOW(),
    UNIQUE(symbol, date_trunc('hour', timestamp))
);

