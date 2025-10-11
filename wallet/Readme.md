# wallet

## Setup Instructions

1. Ensure a PostgreSQL database is running locally on your machine.

2. Update the database user and password in the `Makefile` to match your PostgreSQL credentials.

3. Initialize the database schema: `make db-init`

4. Run database migrations to set up tables and schema: `make db-migrate-up-local`

5. Create env file - `.env.local` from data in `.env.example`
 
6. Build the application binary: `make build`

7. Run binary: `make run-local`

API-spec: event-spec/1.0.0/rest-api
