openapi: 3.0.0
info:
  title: Transaction API
  version: 1.0.0
servers:
  - url: https://transaction-history:8080
security:
  - BearerAuth: []

paths:
  /transactions:
    post:
      summary: Get top N transactions for an account
      description: Retrieves top N transactions for given accountId excluding network_fee transactions, ordered by txTime descending.
      operationId: getTransactions
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - accountId
                - count
              properties:
                accountId:
                  type: string
                  description: Account identifier
                  example: "237a0364-38f9-11ef-86b2-578335f3facf"
                count:
                  type: integer
                  format: int32
                  description: Number of transactions to retrieve
                  example: 10
      responses:
        '200':
          description: A list of transactions
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                      example: 1234
                    txType:
                      type: string
                      example: withdrawal
                    traceId:
                      type: string
                      example: "01G2VVPZXK8J9MQB8TWM11Z5GQ"
                    amountDecimal:
                      type: string
                      example: "2.5"
                    symbol:
                      type: string
                      example: ETH
                    accountId:
                      type: string
                      example: "237a0364-38f9-11ef-86b2-578335f3facf"
                    from:
                      type: string
                      example: "0x6e685a45db4d97ba160fa067cb81b40dfed47245"
                    to:
                      type: string
                      example: "0x7aab1234cd56789ef00123456789abcde1234567"
                    status:
                      type: string
                      example: completed
                    txTime:
                      type: string
                      format: date-time
                      example: "2022-04-21T21:35:59.999Z"
                    eventType:
                      type: string
                      example: crypto
                    cryptoSpecificInfo:
                      description: additional info to display for crypto tx details
                      $ref: "#/components/schemas/CryptoSpecificInfo"
                    fiatSpecificInfo:
                      description: additional info to display for fiat tx details
                      $ref: "#/components/schemas/FiatSpecificInfo"
        '400':
          description: Invalid request parameters
        '500':
          description: Internal server error
  /balances:
    post:
      summary: Get balances for an account
      description: Retrieves balances grouped by symbol, transaction type and event type for the given accountId.
      operationId: getBalances
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - accountId
              properties:
                accountId:
                  type: string
                  description: Account identifier
                  example: "237a0364-38f9-11ef-86b2-578335f3facf"
      responses:
        '200':
          description: List of balances grouped by symbol, txType, and eventType
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    symbol:
                      type: string
                      example: ETH
                    txType:
                      type: string
                      enum: ['withdraw', 'deposit', 'network_fee']
                      example: deposit
                    eventType:
                      type: string
                      enum: ['crypto', 'fiat']
                      example: crypto
                    amountDecimal:
                      type: number
                      format: double
                      example: 123.456
        '400':
          description: Invalid request parameters
        '500':
          description: Internal server error


components:
  schemas:
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

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: 'JWT bearer token in header - "Authorization: Bearer TOKEN"'