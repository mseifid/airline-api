# Firefly Airline Backend

A backend service built with **Go**, **Echo v5**, **GORM**, and **PostgreSQL** to handle Firefly Airline Open API.

## Tech Stack

* Go 1.25+
* Echo v5
* GORM
* PostgreSQL
* `godotenv` — loads environment variables from `.env`
* `envconfig` — maps environment variables to typed Go configuration structs

## Decisions
- We have two kinds of API, one is `internal` that is specific to airline itself and uses airline apikey and the others are the normal exposed APIs for clients. When we have overlapping logics (like listing countries), I just separated the facing API, but internally I use the same service and method to keep them consistent. They normally are different even in exposing side, but for the sake of simplicity I kept it as is.
- Editing an airplane capacity won't affect the flight, cause if it's going to change it, then some modification logic needed if some flight has had sold more that the new maximum.
- Pagination is not implemented for the sake of simplicity. For a production API it is a must.


## Status

Current implementation:

* [x] Environment configuration
* [x] GORM connection
* [x] Echo v5 HTTP server
* [x] Database migrations
* [ ] Models
* [ ] Repositories
* [ ] Use cases / services
* [ ] HTTP handlers
* [ ] Swagger documentation
* [ ] Tests
