CREATE DOMAIN crypto_units AS NUMERIC(50,18);

CREATE SEQUENCE IF NOT EXISTS txs_seq;

CREATE TABLE IF NOT EXISTS txs (
    id INTEGER PRIMARY KEY DEFAULT nextval('txs_seq'),
    tx_type TEXT, -- withdraw or deposit
    trace_id TEXT UNIQUE,
    amount_decimal crypto_units NOT NULL,
    symbol TEXT, -- ticker for crypto, currency for fiat
    account_id TEXT NOT NULL,
    "from" TEXT, -- for crypto this is the from address, for fiat its the sender account identifier 
    "to" TEXT,
    status TEXT,
    tx_time TIMESTAMP DEFAULT NOW(),
    event_type TEXT, -- one of crypto or fiat
    fiat_specific_info jsonb DEFAULT '{}',
    crypto_specific_info jsonb DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- for balance
CREATE INDEX IF NOT EXISTS txs_account_symbol_type_idx ON txs (account_id, symbol, tx_type);
CREATE INDEX IF NOT EXISTS txs_account_id_event_type_idx ON txs (account_id, event_type);
CREATE INDEX IF NOT EXISTS txs_tx_type_tx_time_idx ON txs (tx_type, tx_time DESC);


-- For debug info if need to check what event came
CREATE SEQUENCE IF NOT EXISTS crypto_tx_raw_seq;
CREATE TABLE IF NOT EXISTS crypto_tx_raw (
    id INTEGER PRIMARY KEY DEFAULT nextval('crypto_tx_raw_seq'),
    trace_id TEXT,
    event_data jsonb,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);