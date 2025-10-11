CREATE DATABASE wallet; 

CREATE SEQUENCE IF NOT EXISTS wallet_seq;

CREATE TABLE IF NOT EXISTS wallet (
    id INTEGER PRIMARY KEY DEFAULT nextval('wallet_seq'),
    account_id TEXT NOT NULL,
    seed_id UUID NOT NULL,
    purpose INTEGER NOT NULL,
    coin_type INTEGER NOT NULL,
    account_index INTEGER NOT NULL DEFAULT 0,
    change_index INTEGER NOT NULL DEFAULT 0,
    address_index INTEGER NOT NULL DEFAULT 0,
    public_key TEXT NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL  DEFAULT NOW(),
    UNIQUE (seed_id, purpose, coin_type, account_index, change_index, address_index)
);

CREATE INDEX IF NOT EXISTS wallet_account_id_idx ON wallet (account_id);
CREATE INDEX IF NOT EXISTS wallet_address_idx ON wallet (address);
