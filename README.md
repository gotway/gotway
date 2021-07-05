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

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/9776-3e976745-8b33-46c1-bfe6-d7211722d809?action=collection%2Ffork&collection-url=entityId%3D9776-3e976745-8b33-46c1-bfe6-d7211722d809%26entityType%3Dcollection%26workspaceId%3D10c73242-ad78-405e-b364-b37e56fbb5d3#?env%5BGotway%20Local%5D=W3sia2V5IjoidXJsIiwidmFsdWUiOiJodHRwczovL2dvdHdheToxMTAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoidXJsQ2F0YWxvZyIsInZhbHVlIjoiaHR0cHM6Ly9jYXRhbG9nOjExMDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxTdG9jayIsInZhbHVlIjoiaHR0cHM6Ly9zdG9jazoxMTAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoiaG9zdENhdGFsb2ciLCJ2YWx1ZSI6ImNhdGFsb2c6MTEwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6Imhvc3RTdG9jayIsInZhbHVlIjoic3RvY2s6MTEwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InBhdGhQcmVmaXhSb3V0ZSIsInZhbHVlIjoiL3JvdXRlLlJvdXRlIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJpbnRlcm5hbFVybENhdGFsb2ciLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTIwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6ImludGVybmFsVXJsU3RvY2siLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTMwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6ImludGVybmFsVXJsUm91dGUiLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTQwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InByb2R1Y3RJZCIsInZhbHVlIjoiMTIzNCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoicHJvZHVjdElkMiIsInZhbHVlIjoiNDU2IiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJwcm9kdWN0SWQzIiwidmFsdWUiOiI3ODkiLCJlbmFibGVkIjp0cnVlfV0=)