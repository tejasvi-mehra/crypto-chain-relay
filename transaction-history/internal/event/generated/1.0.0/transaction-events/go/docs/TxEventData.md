# TxEventData

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Specversion** | **string** | versioning for kafka spec | [default to null]
**Source** | **string** | sending app name | [default to null]
**Datacontenttype** | **string** |  | [default to null]
**EventType** | **string** |  | [default to null]
**TraceID** | **string** | ULID, created when new TX is detected | [default to null]
**Status** | [***Status**](Status.md) |  | [default to null]
**TransactionType** | [***TransactionType**](TransactionType.md) |  | [default to null]
**From** | **string** | Identifier for from (address or account) | [default to null]
**To** | **string** | Identifier for to (address or account) | [default to null]
**AccountID** | **string** | Internal Identifier for user account related to this tx | [default to null]
**AmountDecimal** | **string** | Amount for UI, for crypto withdrawal this is the sum of withdrawal | [default to null]
**Symbol** | **string** | symbol used to identify whats being transferred - ETH,USD,BTC,USDT etc | [default to null]
**TransactionTime** | **string** |  | [optional] [default to null]
**CreatedAt** | **string** |  | [default to null]
**UpdatedAt** | **string** |  | [default to null]
**CryptoSpecificInfo** | [***CryptoSpecificInfo**](CryptoSpecificInfo.md) |  | [optional] [default to null]
**FiatSpecificInfo** | [***FiatSpecificInfo**](FiatSpecificInfo.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

