DROP TABLE IF EXISTS crypto_tx_raw;
DROP SEQUENCE IF EXISTS crypto_tx_raw_seq;

DROP INDEX IF EXISTS txs_account_symbol_type_idx;
DROP INDEX IF EXISTS txs_account_id_event_type_idx;
DROP INDEX IF EXISTS txs_tx_type_tx_time_idx;

DROP TABLE IF EXISTS txs;
DROP SEQUENCE IF EXISTS txs_seq;

DROP DOMAIN IF EXISTS crypto_units;
