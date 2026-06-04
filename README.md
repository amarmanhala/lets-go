# Countries API

A minimal Go backend with no web framework. It exposes one endpoint using the standard library and stores country data in PostgreSQL.

The app reads these database environment variables:

```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
PORT=8080
```

## Run

```sh
DB_HOST=localhost DB_USER=admin DB_PASSWORD=password DB_NAME=mydb go run .
```

## Docker

```sh
docker build -t lets-go:v1 .
docker network create lets-go-net
docker network connect lets-go-net postgres
docker run --rm \
  --name lets-go-api \
  --network lets-go-net \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=admin \
  -e DB_PASSWORD=password \
  -e DB_NAME=mydb \
  lets-go:v1
```

## Endpoint

```sh
curl http://localhost:8080/api/countries
```

Returns 50 mocked country records:

```json
[
  {
    "name": "Canada",
    "population": 40097761,
    "flag_colors": ["red", "white"],
    "language": "English"
  }
]
```
# lets-go
