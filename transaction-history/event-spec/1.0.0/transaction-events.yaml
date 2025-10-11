openapi: 3.0.0
info:
  title: |
    TOPIC: transaction-events
  version: 1.0.0
servers: []
paths: {}
components:
  schemas:
    TxEventData:
      type: object
      required:
        - specversion
        - source
        - datacontenttype
        - eventType
        - traceID
        - status
        - transactionType
        - from
        - to
        - accountID
        - amountDecimal
        - symbol
        - transaction_time
        - createdAt
        - updatedAt
      properties:
        specversion:
          type: string
          enum: ['1.0']
          description: versioning for kafka spec
        source:
          type: string
          description: sending app name
        datacontenttype:
          type: string
          enum: ['application/json']
        eventType: 
          type: string
          enum: ['crypto','fiat']
        traceID:
          description: 'ULID, created when new TX is detected'
          type: string
          example: '01G2VVPZXK8J9MQB8TWM11Z5GQ'
        status:
          $ref: '#/components/schemas/Status'
        transactionType:
          $ref: '#/components/schemas/TransactionType'
        from: 
          type: string
          description: 'Identifier for from (address or account)'
          example: '0x6e685a45db4d97ba160fa067cb81b40dfed47245'
        to: 
          type: string
          description: 'Identifier for to (address or account)'
          example: '0x6e685a45db4d97ba160fa067cb81b40dfed47245'
        accountID:
          type: string
          description: 'Internal Identifier for user account related to this tx'
          example: '237a0364-38f9-11ef-86b2-578335f3facf'
        amountDecimal:
          type: string
          description: 'Amount for UI, for crypto withdrawal this is the sum of withdrawal'
          example: 2.5
        symbol:
          type: string
          description: "symbol used to identify whats being transferred - ETH,USD,BTC,USDT etc"
          example: 'ETH'
        transactionTime:
          type: string
          example: '2022-04-21T21:35:59.999Z'
        createdAt:
          type: string
          example: '2022-04-21T21:35:59.999Z'
        updatedAt:
          type: string
          example: '2022-04-21T21:35:59.999Z'
        # optional fields
        cryptoSpecificInfo: 
          $ref: '#/components/schemas/CryptoSpecificInfo'
        fiatSpecificInfo: 
          $ref: '#/components/schemas/FiatSpecificInfo'

    Status:
      type: string
      enum: [
        detected,
        completed,
        failed,
        orphaned
      ]

    TransactionType:
      type: string
      enum: [
        deposit,
        withdrawal,
        network_fee
      ]

    CryptoSpecificInfo:
      type: object
      required:
        - blockchaintype
        - chainType
        - chainID
        - blockNo
        - blockTimestamp
        - fromAddress
        - toAddress
        - amountInfo
        - txHash
        - transactionType
      properties:
        blockchaintype: 
          type: string
          description: "blockchaintype if its crypto"
          example: "evm,bitcoin"
        chainID:
          type: string
          description: chain id of the chain
          example: '1'
        block:
          type: integer
          description: 'Block number in Decimal format'
          example: 27537537
        blockTimestamp:
          type: string
          description: 'Block time from the chain'
          example: '2022-05-17T22:00:00.000Z'
        fromAddress:
          type: string
          description: 'Address (Optional)'
          example: '0x6e685a45db4d97ba160fa067cb81b40dfed47245'
        toAddress:
          type: string
          description: 'Address (Optional)'
          example: '0x6e685a45db4d97ba160fa067cb81b40dfed47245'
        amountInfo:
          $ref: '#/components/schemas/TransferAmountInfo'
        txHash: 
          type: string
          description: 'transaction hash'
          example: '0x7c6bd55ce2986c33a35a061660ae5b8a41262df0129c1e36c32da2a06289fc0c'

    TransferAmountInfo:
      type: object
      required:
        - ticker
        - valueAtomic
        - valueDecimal
        - decimals
        - standard
      properties:
        ticker:
          type: string
          example: 'ETH'
        valueAtomic:
          type: string
          description: 'Atomic Amount in DEC'
          example: '2500000000000000000'
        valueDecimal:
          type: string
          description: 'Amount for UI'
          example: '2.5'
        decimals:
          type: integer
          description: 'Decimal of the asset'
          example: 18
        contractAddress: # Optional
          type: string
          description: 'Contract Address'
          example: '0x6e685a45db4d97ba160fa067cb81b40dfed47245'
        standard:
          type: string
          description: 'spec of asset transferred, native or token standard.'
          example: 'NATIVE, ERC-20'

    FiatSpecificInfo:
      type: object
      description: placeholder for any fiat specific fields we want to add based on previous implementation


