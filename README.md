# gotway 

[![CI](https://github.com/gotway/gotway/actions/workflows/ci.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/ci.yml)
[![Release](https://github.com/gotway/gotway/actions/workflows/release.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/release.yml)
[![Deploy](https://github.com/gotway/gotway/actions/workflows/deploy.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gotway/gotway)](https://goreportcard.com/report/github.com/gotway/gotway)
[![Go Reference](https://pkg.go.dev/badge/github.com/gotway/gotway.svg)](https://pkg.go.dev/github.com/gotway/gotway)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gotway)](https://artifacthub.io/packages/search?repo=gotway)

☸️ Cloud native API Gateway powered with in-redis cache.

- API composition: expose your services to the internet using a single endpoint
- Cloud native: configure routing and cache using [Kubernetes CRDs](./manifests/examples/catalog.yml)
- In-memory cache using redis 
- Cache invalidation using tags
- Health checking
- Management [REST API](#management-rest-api-)
- ~6MB [Docker image](https://hub.docker.com/r/gotwaygateway/gotway/tags) available for multiple architectures
- [Helm chart](https://github.com/gotway/charts)

#### Installation 🌱

```bash
helm repo add gotway https://charts.gotway.duckdns.org
```
```bash
helm install gotway gotway/gotway
```

#### Quickstart ⚡

Let's register the [catalog](https://github.com/gotway/service-examples/tree/main/cmd/catalog) service into Gotway by creating an `IngressHTTP` CRD:

```bash
kubectl apply -f ./manifests/examples/catalog.yml 
``` 
```yaml
apiVersion: gotway.io/v1alpha1
kind: IngressHTTP
metadata:
  name: catalog
spec:
  match:
    host: catalog.gotway.duckdns.org:4433
  service:
    name: catalog
    url: http://gotway-catalog
    healthPath: /health
  cache:
    ttl: 30
    statuses:
      - 200
      - 404
    tags:
      - "catalog"
      - "products"
```

We are now able to route requests through Gotway, let's create a product:

```bash 
curl -k --request POST 'https://catalog.gotway.duckdns.org:4433/product' \
--header 'Content-Type: application/json' \
--data-raw '{
	"name": "sneakers",
	"price": 69000,
	"color": "white",
	"size": "42"
}'
```
```json
{
    "id": 911902081
}
``` 
```bash
curl -k --request GET 'https://catalog.gotway.duckdns.org:4433/products'
```
```json
{
    "products": [
        {
            "id": 911902081,
            "name": "sneakers",
            "price": 69000,
            "color": "white",
            "size": "42"
        }
    ],
    "totalCount": 1
}
``` 

This response has a TTL of 30 seconds, let's invalidate the cache for the catalog service by providing one of its tags:

```bash
curl -k --request DELETE 'https://gotway.duckdns.org:4433/api/cache' \
--header 'Content-Type: application/json' \
--data-raw '{
    "tags": ["catalog"]
}'
``` 


#### Management REST API ⚡

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/9776-3e976745-8b33-46c1-bfe6-d7211722d809?action=collection%2Ffork&collection-url=entityId%3D9776-3e976745-8b33-46c1-bfe6-d7211722d809%26entityType%3Dcollection%26workspaceId%3D10c73242-ad78-405e-b364-b37e56fbb5d3#?env%5BGotway%20Local%5D=W3sia2V5IjoidXJsIiwidmFsdWUiOiJodHRwczovL2dvdHdheToxMTAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoidXJsQ2F0YWxvZyIsInZhbHVlIjoiaHR0cHM6Ly9jYXRhbG9nOjExMDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxTdG9jayIsInZhbHVlIjoiaHR0cHM6Ly9zdG9jazoxMTAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoiaG9zdENhdGFsb2ciLCJ2YWx1ZSI6ImNhdGFsb2c6MTEwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6Imhvc3RTdG9jayIsInZhbHVlIjoic3RvY2s6MTEwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InBhdGhQcmVmaXhSb3V0ZSIsInZhbHVlIjoiL3JvdXRlLlJvdXRlIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJpbnRlcm5hbFVybENhdGFsb2ciLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTIwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6ImludGVybmFsVXJsU3RvY2siLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTMwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6ImludGVybmFsVXJsUm91dGUiLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MTQwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InByb2R1Y3RJZCIsInZhbHVlIjoiMTIzNCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoicHJvZHVjdElkMiIsInZhbHVlIjoiNDU2IiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJwcm9kdWN0SWQzIiwidmFsdWUiOiI3ODkiLCJlbmFibGVkIjp0cnVlfV0=)
