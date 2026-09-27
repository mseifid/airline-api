# Firefly Airline Backend

A backend service built with **Go**, **Echo v5**, **GORM**, and **PostgreSQL** to handle Firefly Airline Open API.

## Tech Stack

* Go 1.26.8+
* Echo v5
* GORM
* PostgreSQL

## Decisions
- We have two kinds of API, one is `internal` that is specific to airline itself and uses airline apikey and the others are the normal exposed APIs for clients. When we have overlapping logics (like listing countries), I just separated the facing API, but internally I use the same service and method to keep them consistent. They normally are different even in exposing side, but for the sake of simplicity I kept it as is.
- Editing an airplane capacity won't affect the flight, cause if it's going to change it, then some modification logic needed if some flight has had sold more that the new maximum.
- Pagination is not implemented for the sake of simplicity. For a production API it is a must.
- Ticket purchase and cancellation wa done using pessimistic locking scenario.
- We could also have ticket cancellation for airline, but I ignored it for now.

## Things to Improve
- A cron job is needed on production to deactivate the flights when their departure time arrives.
- More sophisticated error handling , and returning domain errors directly to global error handler to be shown in response.
- A seeder could make the first-time-run easier, but because of time shortage I ignored it
- A more sophisticated logging scenario using `slog` could be implemented.
- API Key management would be more production-grade by implementing revoke scenario.

## How to run
1. Clone the repository and install the Go dependencies: `go mod download`
2. Create a .env file based on the `.env.example` file keys.
3. Run application: `go run cmd/main.go`
4. Swagger is available here: `http://localhost:{APP_PORT}/swagger/`