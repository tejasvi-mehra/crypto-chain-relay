INSERT INTO txs (
    tx_type,
    trace_id,
    amount_decimal,
    symbol,
    account_id,
    "from",
    "to",
    status,
    tx_time,
    event_type,
    fiat_specific_info,
    created_at,
    updated_at
) VALUES (
    'deposit',
    '01G3ABCDXYZ9876543210EFGHJK',
    '1000.00'::crypto_units,
    'USD',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    'bank-account-123456789',
    'internal-account-987654321',
    'completed',
    '2025-10-01T10:15:30.000Z'::timestamp,
    'fiat',
    '{
        "paymentProvider": "stripe",
        "referenceNumber": "STR123456789",
        "currencyCode": "USD",
        "paymentMethod": "ACH"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'deposit',
    '01G3ABCDXYZ9876543210EFGHJK-1',
    '2000.00'::crypto_units,
    'USD',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    'bank-account-123456789',
    'internal-account-987654321',
    'completed',
    '2025-10-01T10:15:30.000Z'::timestamp,
    'fiat',
    '{
        "paymentProvider": "stripe",
        "referenceNumber": "STR123456789",
        "currencyCode": "USD",
        "paymentMethod": "ACH"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'withdrawal',
    '01G3ABCDXYZ9876543210EFGHJK-2',
    '500.00'::crypto_units,
    'USD',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    'bank-account-123456789',
    'internal-account-987654321',
    'completed',
    '2025-10-01T10:15:30.000Z'::timestamp,
    'fiat',
    '{
        "paymentProvider": "stripe",
        "referenceNumber": "STR123456789",
        "currencyCode": "USD",
        "paymentMethod": "ACH"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'deposit',
    '01G3ABCDXYZ9876543210EFGHJK-3',
    '1000.00'::crypto_units,
    'USD',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    'bank-account-123456789',
    'internal-account-987654321',
    'completed',
    '2025-10-01T10:15:30.000Z'::timestamp,
    'fiat',
    '{
        "paymentProvider": "stripe",
        "referenceNumber": "STR123456789",
        "currencyCode": "USD",
        "paymentMethod": "ACH"
    }'::jsonb,
    NOW(),
    NOW()
);

INSERT INTO txs (
    tx_type,
    trace_id,
    amount_decimal,
    symbol,
    account_id,
    "from",
    "to",
    status,
    tx_time,
    event_type,
    crypto_specific_info,
    created_at,
    updated_at
) VALUES (
    'deposit',  -- from "transactionType": "withdrawal"
    '01G2VVPZXK8J9MQB8TWM11Z5GR',  -- traceID
    '110.75'::crypto_units,  -- amountDecimal
    'USDT',  -- symbol
    '237a0364-38f9-11ef-86b2-578335f3facf',  -- accountID
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',  -- from
    '0x7aab1234cd56789ef00123456789abcde1234567',  -- to
    'completed',  -- status
    '2022-04-21T21:35:59.999Z'::timestamp,  -- transactionTime
    'crypto',  -- eventType
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "USDT",
            "valueAtomic": "150750000",
            "valueDecimal": "150.75",
            "decimals": 6,
            "contractAddress": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
            "standard": "ERC-20"
        },
        "txHash": "0xa1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef",
        "transactionType": "withdrawal"
    }'::jsonb,  -- cryptoSpecificInfo
    NOW(),
    NOW()
),(
    'deposit',  -- from "transactionType": "withdrawal"
    '01G2VVPZXK8J9MQB8TWM11Z5GR-1',  -- traceID
    '250.75'::crypto_units,  -- amountDecimal
    'USDT',  -- symbol
    '237a0364-38f9-11ef-86b2-578335f3facf',  -- accountID
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',  -- from
    '0x7aab1234cd56789ef00123456789abcde1234567',  -- to
    'completed',  -- status
    '2022-04-21T21:35:59.999Z'::timestamp,  -- transactionTime
    'crypto',  -- eventType
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "USDT",
            "valueAtomic": "150750000",
            "valueDecimal": "150.75",
            "decimals": 6,
            "contractAddress": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
            "standard": "ERC-20"
        },
        "txHash": "0xa1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef",
        "transactionType": "withdrawal"
    }'::jsonb,  -- cryptoSpecificInfo
    NOW(),
    NOW()
),(
    'withdrawal',  -- from "transactionType": "withdrawal"
    '01G2VVPZXK8J9MQB8TWM11Z5GR-2',  -- traceID
    '50.75'::crypto_units,  -- amountDecimal
    'USDT',  -- symbol
    '237a0364-38f9-11ef-86b2-578335f3facf',  -- accountID
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',  -- from
    '0x7aab1234cd56789ef00123456789abcde1234567',  -- to
    'completed',  -- status
    '2022-04-21T21:35:59.999Z'::timestamp,  -- transactionTime
    'crypto',  -- eventType
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "USDT",
            "valueAtomic": "150750000",
            "valueDecimal": "150.75",
            "decimals": 6,
            "contractAddress": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
            "standard": "ERC-20"
        },
        "txHash": "0xa1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef",
        "transactionType": "withdrawal"
    }'::jsonb,  -- cryptoSpecificInfo
    NOW(),
    NOW()
);

INSERT INTO txs (
    tx_type,
    trace_id,
    amount_decimal,
    symbol,
    account_id,
    "from",
    "to",
    status,
    tx_time,
    event_type,
    crypto_specific_info,
    created_at,
    updated_at
) VALUES (
    'deposit',
    '01G2WW8JXPB7YQF4W8E3FF9YTA',
    '0.12345678'::crypto_units,
    'BTC',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
    '1BoatSLRHtKNngkdXEeobR76b53LETtpyT',
    'completed',
    '2025-10-01T15:30:22.000Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "bitcoin",
        "chainID": "0",
        "block": 781234,
        "blockTimestamp": "2025-10-01T15:00:00.000Z",
        "fromAddress": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        "toAddress": "1BoatSLRHtKNngkdXEeobR76b53LETtpyT",
        "amountInfo": {
            "ticker": "BTC",
            "valueAtomic": "12345678",
            "valueDecimal": "0.12345678",
            "decimals": 8,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "e2a1d9b779a1dd4a3eb8f2bedaa9b31cbe9a765ee0f5b0a5264f2e7771c2f4ab",
        "transactionType": "deposit"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'deposit',
    '01G2WW8JXPB7YQF4W8E3FF9YTA-1',
    '0.5'::crypto_units,
    'BTC',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
    '1BoatSLRHtKNngkdXEeobR76b53LETtpyT',
    'completed',
    '2025-10-01T15:30:22.000Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "bitcoin",
        "chainID": "0",
        "block": 781234,
        "blockTimestamp": "2025-10-01T15:00:00.000Z",
        "fromAddress": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        "toAddress": "1BoatSLRHtKNngkdXEeobR76b53LETtpyT",
        "amountInfo": {
            "ticker": "BTC",
            "valueAtomic": "12345678",
            "valueDecimal": "0.12345678",
            "decimals": 8,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "e2a1d9b779a1dd4a3eb8f2bedaa9b31cbe9a765ee0f5b0a5264f2e7771c2f4ab",
        "transactionType": "deposit"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'withdrawal',
    '01G2WW8JXPB7YQF4W8E3FF9YTA-2',
    '0.05'::crypto_units,
    'BTC',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa',
    '1BoatSLRHtKNngkdXEeobR76b53LETtpyT',
    'completed',
    '2025-10-01T15:30:22.000Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "bitcoin",
        "chainID": "0",
        "block": 781234,
        "blockTimestamp": "2025-10-01T15:00:00.000Z",
        "fromAddress": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        "toAddress": "1BoatSLRHtKNngkdXEeobR76b53LETtpyT",
        "amountInfo": {
            "ticker": "BTC",
            "valueAtomic": "12345678",
            "valueDecimal": "0.12345678",
            "decimals": 8,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "e2a1d9b779a1dd4a3eb8f2bedaa9b31cbe9a765ee0f5b0a5264f2e7771c2f4ab",
        "transactionType": "deposit"
    }'::jsonb,
    NOW(),
    NOW()
);

INSERT INTO txs (
    tx_type,
    trace_id,
    amount_decimal,
    symbol,
    account_id,
    "from",
    "to",
    status,
    tx_time,
    event_type,
    crypto_specific_info,
    created_at,
    updated_at
) VALUES (
    'deposit',
    '01G2VVPZXK8J9MQB8TWM11Z5GQ',
    '2.5'::crypto_units,
    'ETH',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',
    '0x7aab1234cd56789ef00123456789abcde1234567',
    'completed',
    '2022-04-21T21:35:59.999Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "ETH",
            "valueAtomic": "2500000000000000000",
            "valueDecimal": "2.5",
            "decimals": 18,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "0x7c6bd55ce2986c33a35a061660ae5b8a41262df0129c1e36c32da2a06289fc0c",
        "transactionType": "withdrawal"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'withdraw',
    '01G2VVPZXK8J9MQB8TWM11Z5GQ-1',
    '1.5'::crypto_units,
    'ETH',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',
    '0x7aab1234cd56789ef00123456789abcde1234567',
    'completed',
    '2022-04-21T21:35:59.999Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "ETH",
            "valueAtomic": "2500000000000000000",
            "valueDecimal": "2.5",
            "decimals": 18,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "0x7c6bd55ce2986c33a35a061660ae5b8a41262df0129c1e36c32da2a06289fc0c",
        "transactionType": "withdrawal"
    }'::jsonb,
    NOW(),
    NOW()
),(
    'deposit',
    '01G2VVPZXK8J9MQB8TWM11Z5GQ-2',
    '12.5'::crypto_units,
    'ETH',
    '237a0364-38f9-11ef-86b2-578335f3facf',
    '0x6e685a45db4d97ba160fa067cb81b40dfed47245',
    '0x7aab1234cd56789ef00123456789abcde1234567',
    'completed',
    '2022-04-21T21:35:59.999Z'::timestamp,
    'crypto',
    '{
        "blockchaintype": "evm",
        "chainID": "1",
        "block": 27537537,
        "blockTimestamp": "2022-05-17T22:00:00.000Z",
        "fromAddress": "0x6e685a45db4d97ba160fa067cb81b40dfed47245",
        "toAddress": "0x7aab1234cd56789ef00123456789abcde1234567",
        "amountInfo": {
            "ticker": "ETH",
            "valueAtomic": "2500000000000000000",
            "valueDecimal": "2.5",
            "decimals": 18,
            "contractAddress": "",
            "standard": "NATIVE"
        },
        "txHash": "0x7c6bd55ce2986c33a35a061660ae5b8a41262df0129c1e36c32da2a06289fc0c",
        "transactionType": "withdrawal"
    }'::jsonb,
    NOW(),
    NOW()
);