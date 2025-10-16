# trial-project-ct


This repo contains 4 projects in total - 
- builder-router
- eth_adapter
- transaction-history
- wallet
Each project contains its own Makefile for setup, event-spec for api-specs, postman collection for sample calls and db migrations for any SQL tables 

## builder-router
This service would be responsible for routing the withdrawal tx to either eth-builer, btc-builder or fiat-builder. The API specs are for builder-router and chain-specific routers are under event/spec. This service uses mocks for transaction-history and implements a sample mock eth-builder and btc-builder which should be built separately. transction-history is used to get the balance to check if amount is actually withdrawable. There is also JWK validation to match accountID in the request with JWT subject to ensure only user who holds the account can make withdrawls. Gas/nonce etc is mocked for now but should be fetched from the node if not provided. This service calls the subservices eth-builder or btc-builder to do withdrawal. Those services build unsigned tx, call wallet to sign and then broadcast the response 

## eth_adapter
This service connects to the ETH Node to detect transactions from the blocks, mark them as completed or orphan after certain blocks have passed. The publisher/db are mocked and just print the json output in console. The publisher output can be used to add new txs in transction-history. This actually connects to the node for now and gets info for certain txs for 2 addresses as defined in address map in main. Block no can be changed to get a different block id need to 

## transaction-history
This service keeps the transactions as detected by eth_adapter, btc_adapter or fiat_adapter - same examples for txs are under event-spec/1.0.0/example-txs.json. This service also returns the tx history and balance for an account. This service should also have JWT but is not added in this project - similar implementation as to what is in builder-router to ensure we only return txs for the account with valid JWT

## wallet
This service helps in generating a new unqiue address for an account + chainID using HDPaths. It also is responsible for signing txs. It has a subservice called KMS but its abstracted so it can be moved to a more secure location instead. Each accountID generates a particular seed which is used to generate a fixed address for that accountID. During signing as well we need to only provide accountID and unsignedTx as build through builder-router to sign the tx. For eth it expects unsignedTx, for btc it expects psbt