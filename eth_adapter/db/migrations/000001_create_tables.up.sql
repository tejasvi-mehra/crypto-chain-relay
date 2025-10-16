-- addresess table 
CREATE SEQUENCE IF NOT EXISTS addresses_seq;
CREATE TABLE IF NOT EXISTS addresses (
  id INTEGER PRIMARY KEY DEFAULT nextval('addresses_seq'::regclass),
  address TEXT NOT NULL UNIQUE,
  account_id TEXT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS addresses_address_idx ON addresses (address);
CREATE INDEX IF NOT EXISTS account_id_address_idx ON addresses (account_id);

-- blocks data table
CREATE SEQUENCE IF NOT EXISTS block_data_partitioned_seq;

CREATE TABLE IF NOT EXISTS block_data (
  id INTEGER NOT NULL DEFAULT nextval('block_data_partitioned_seq'::regclass),
  block_no INTEGER,
  block_hash TEXT,
  parent_hash TEXT,
  block_timestamp TIMESTAMP DEFAULT now(),
  process_time_ms INTEGER,
  confirmed BOOLEAN,
  retry_count INTEGER,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE INDEX IF NOT EXISTS block_data_block_no_idx ON block_data (block_no, created_at);
CREATE INDEX IF NOT EXISTS block_data_block_hash_idx ON block_data (block_hash, created_at);
CREATE INDEX IF NOT EXISTS block_data_parent_hash_idx ON block_data (parent_hash, created_at);
CREATE INDEX IF NOT EXISTS block_data_confirmed_idx ON block_data (confirmed, created_at);

-- Sample partitions for August and September 2025
CREATE TABLE IF NOT EXISTS block_data_2025_10 PARTITION OF block_data 
  FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE IF NOT EXISTS block_data_2025_11 PARTITION OF block_data 
  FOR VALUES FROM ('2025-1-01') TO ('2025-12-01');

-- transactions table 

CREATE SEQUENCE IF NOT EXISTS chain_txs_seq;

CREATE TABLE IF NOT EXISTS chain_txs (
  id INTEGER PRIMARY KEY DEFAULT nextval('chain_txs_seq'),
  block_no INTEGER,
  block_time TIMESTAMP,
  tx_type TEXT, -- withdraw or deposit
  tx_hash TEXT,
  trace_id TEXT UNIQUE,
  gas_limit TEXT,
  gas_price TEXT,
  amount TEXT,
  ticker TEXT,
  spec TEXT, -- one of ERC-20 or NATIVE
  account_id TEXT,
  action_address TEXT,
  from_address TEXT,
  to_address TEXT,
  amount_info JSONB, -- value in atomic/decimal/USD
  status TEXT, -- detected/completed/failed/orphan
  kafka_event_data JSONB,
  created_at  TIMESTAMP DEFAULT now(),
  updated_at  TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS chain_txs_action_address_idx ON chain_txs (action_address);
CREATE INDEX IF NOT EXISTS chain_txs_tx_hash_idx ON chain_txs (tx_hash);
CREATE INDEX IF NOT EXISTS chain_txs_block_no_idx ON chain_txs (block_no);
CREATE INDEX IF NOT EXISTS chain_txs_trace_id_idx ON chain_txs (trace_id);
