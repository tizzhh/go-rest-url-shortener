# Url-shortener

A go-based url-shortener with authentication.

## Usage 💡

- clone the repository
- `make`
- `./url-shortener`

## Requirements

- Go 1.22.5
- config in the format of local.yaml in config/ directory. Path to config should be set in CONFIG_PATH env variable

## Used packages / tools / stack

- REST.
- Sqlite3.
- Testing with [Mockery](https://github.com/vektra/mockery).
- The usage of slog as the centralized logger.
- [Chi](https://github.com/go-chi/chi) for routing.


## Endpoints

| Name        | HTTP Method | Route          |
|-------------|-------------|----------------|
| Create a new alias | POST | /url |
| Delete an alias | DELETE | /url/{alias} |
| Get a redirect from alias | GET | /{alias}

##  Database design

#### url

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | INT      | ✅        | ✅           |
| alias          | TEXT      | ✅        |             |
| url         | TEXT      | ✅        |             |
