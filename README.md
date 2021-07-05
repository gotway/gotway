# gotway 

[![CI](https://github.com/gotway/gotway/actions/workflows/ci.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/ci.yml)
[![Release](https://github.com/gotway/gotway/actions/workflows/release.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/release.yml)
[![Deploy](https://github.com/gotway/gotway/actions/workflows/deploy.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gotway/gotway)](https://goreportcard.com/report/github.com/gotway/gotway)
[![Go Reference](https://pkg.go.dev/badge/github.com/gotway/gotway.svg)](https://pkg.go.dev/github.com/gotway/gotway)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gotway)](https://artifacthub.io/packages/search?repo=gotway)

Simple HTTP API Gateway powered with in-redis cache 🚀

- API composition: expose your services to the internet using a single endpoint
- Configurable cache in redis 
- Cache invalidation using tags
- Cache invalidation specifying the URL path
- Health checking
- Management REST API 

- ~10MB [Docker image](https://hub.docker.com/r/gotwaygateway/gotway/tags) available for multiple architectures

#### Installation 🌱

```bash
$ helm repo add gotway https://charts.gotway.duckdns.org
$ helm install gotway gotway/gotway
```

#### Quickstart ⚡

We will register [catalog](https://github.com/gotway/service-examples/tree/main/cmd/catalog) as an example:

```bash
curl --request POST 'https://api.gotway.duckdns.org/api/service' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "catalog",
    "match": {
        "host": "catalog.gotway.duckdns.org"
    },
    "backend": {
        "url": "http://catalog:80"
    },
    "cache": {
        "ttl": 30,
        "statuses": [200, 404],
        "tags": ["catalog", "products"]
     }
}'
```

After executing that command, our service will be available at
[https://catalog.gotway.duckdns.org](https://catalog.gotway.duckdns.org). The following endpoints will be routed through gotway:

- **GET** https://catalog.gotway.duckdns.org/products
- **POST** https://catalog.gotway.duckdns.org/product
- **GET** https://catalog.gotway.duckdns.org/product/1
- **DELETE** https://catalog.gotway.duckdns.org/product/1
- **PUT** https://catalog.gotway.duckdns.org/product/1

#### Management REST API ⚡

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/2e80e5165001548d7d43#?env%5BGotway%20Local%5D=W3sia2V5IjoidXJsIiwidmFsdWUiOiJodHRwczovL2xvY2FsaG9zdDo4MDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxDYXRhbG9nIiwidmFsdWUiOiJodHRwOi8vbG9jYWxob3N0OjkwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InVybFJvdXRlIiwidmFsdWUiOiJodHRwOi8vbG9jYWxob3N0OjExMDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxTdG9jayIsInZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDoxMDAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoicHJvZHVjdElkIiwidmFsdWUiOiIxMjM0IiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJwcm9kdWN0SWQyIiwidmFsdWUiOiI0NTYiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InByb2R1Y3RJZDMiLCJ2YWx1ZSI6Ijc4OSIsImVuYWJsZWQiOnRydWV9XQ==)