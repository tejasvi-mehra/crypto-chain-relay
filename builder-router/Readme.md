# wallet

## Setup Instructions

1. Create env file - `.env.local` from data in `.env.example`
 
2. Build the application binary: `make build`

3. Run binary: `make run-local`

API-spec: event-spec/1.0.0/rest-api

Most of the code on this is mock - this service is more to showcase some basic common validations we need to do before we initiate a withdrawal along with how the routing should work based on what withdrawal tx it is. Mock services for ETH/BTC/Fiat show how the specific builder should work in this case - ETH builder is more complete and works if used along with the code for wallet in the repo