# gotway
[![CI](https://github.com/gotway/gotway/actions/workflows/ci.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/ci.yml)
[![Release](https://github.com/gotway/gotway/actions/workflows/release.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/release.yml)
[![Deploy](https://github.com/gotway/gotway/actions/workflows/deploy.yml/badge.svg)](https://github.com/gotway/gotway/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gotway/gotway)](https://goreportcard.com/report/github.com/gotway/gotway)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gotway)](https://artifacthub.io/packages/search?repo=gotway)

A simple, lightweight and blazingly fast API gateway 🚀

- API composition and dynamic routing
- Support for **REST** and **gRPC** microservices
- Configuration and object management via **gotway REST API**
- **Discover services** dynamically in runtime by registering them in gotway API
- **Health checking** to make sure everything is up and running
- **Cache** your service responses temporarily in gotway's Redis for improving your API response time
- **Cache invalidation** using tags and paths via gotway API
- ~10MB [Docker image](https://hub.docker.com/r/gotwaygateway/gotway/tags) available for multiple architectures

---

- [Installation 🌱](#installation-)
- [Roadmap 🛣️](https://github.com/gotway/gotway/milestone/1)
- [Features ⚡](#features-)
    - [Service discovery 🔭](#service-discovery-)
    - [Health checking 🚑](#health-checking-)
    - [Cache 💾](#cache-)

---

## Installation 🌱

###### Environment variables
Set up this [env variables](./config/config.go) for configuring your gotway instance.


###### Install from source

```bash
$ docker-compose -f docker-compose.redis.yml up -d
$ make run
```

###### Docker

```bash
$ docker-compose -f docker-compose.redis.yml -f docker-compose.yml up -d
```

###### Kubernetes + Helm
```bash
$ helm repo add gotway https://charts.gotway.duckdns.org
$ helm install gotway gotway/gotway
```

## Features ⚡

#### Service discovery 🔭

Services can be discovered in runtime by registering them in the gotway API.

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/2e80e5165001548d7d43#?env%5BGotway%20Local%5D=W3sia2V5IjoidXJsIiwidmFsdWUiOiJodHRwczovL2xvY2FsaG9zdDo4MDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxDYXRhbG9nIiwidmFsdWUiOiJodHRwOi8vbG9jYWxob3N0OjkwMDAiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InVybFJvdXRlIiwidmFsdWUiOiJodHRwOi8vbG9jYWxob3N0OjExMDAwIiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJ1cmxTdG9jayIsInZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDoxMDAwMCIsImVuYWJsZWQiOnRydWV9LHsia2V5IjoicHJvZHVjdElkIiwidmFsdWUiOiIxMjM0IiwiZW5hYmxlZCI6dHJ1ZX0seyJrZXkiOiJwcm9kdWN0SWQyIiwidmFsdWUiOiI0NTYiLCJlbmFibGVkIjp0cnVlfSx7ImtleSI6InByb2R1Y3RJZDMiLCJ2YWx1ZSI6Ijc4OSIsImVuYWJsZWQiOnRydWV9XQ==)

###### REST

We will register [catalog](./cmd/catalog) as an example:

```bash
curl --request POST 'https://<gotway>/api/service' \
--header 'Content-Type: application/json' \
--data-raw '{
    "type": "rest",
    "url": "http://<catalog>",
    "path": "catalog"
}'
```

After executing that command, our service will be available at
`https://<gotway>/<path>`. The following endpoints will be routed through gotway:

- `GET https://<gotway>/catalog/products`
- `POST https://<gotway>/catalog/product`
- `GET https://<gotway>/catalog/product/<id>`
- `DELETE https://<gotway>/catalog/product/<id>`
- `PUT https://<gotway>/catalog/product/<id>`

###### gRPC

We will register [route](./cmd/route) as an example:

```bash
curl --request POST 'https://<gotway>/api/service' \
--header 'Content-Type: application/json' \
--data-raw '{
    "type": "grpc",
    "url": "http://<route>:<port>",
    "path": "route.Route"
}'
```

Where `route.Route` represents the package and the service name of your gRPC service. Another example could be:

[grpc.health.v1.Health](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

This will be defined in your [.proto](./cmd/route/pb/route.proto) file and will be used as path in the dynamic routing.

In this case, the RPC methods routed through gotway will be:

- `https://<gotway>/route.Route/GetFeature`
- `https://<gotway>/route.Route/ListFeatures`
- `https://<gotway>/route.Route/RecordRoute`
- `https://<gotway>/route.Route/RouteChat`

For testing them, we have a [gRPC go client](./cmd/route/client/client.go):
```bash
$ cd cmd/route
$ make cli
```

#### Health checking 🚑

gotway will make a health probe to check that our services are responding. In other case, a `502 Bad Gateway` will be returned.

###### REST

By default, the health probe will be done by requesting `http://<microservice>/health`. However, it is posible to use a custom path by specifying `healthPath` when registering.

An example of REST health endpoint is available [here](./cmd/catalog/api/api.go).

###### gRPC

By default, the standard [gRPC health checking protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) is used. However, it is posible to use another one by specifying `healthPath` when registering.

An example of gRPC health checking protocol implementation can be found [here](./cmd/route/server/server.go).

#### Cache 💾

Store microservice responses temporarily in gotway for improving your API response time. You will need to specify the cache policy when registring your service:

```bash
curl --request POST 'https://<gotway>/api/service' \
--header 'Content-Type: application/json' \
--data-raw '{
    "type": "rest",
    "url": "http://<catalog>",
    "path": "catalog",
    "cache": {
        "ttl": 30,
        "statuses": [200, 404],
        "tags": ["catalog", "products"]
     }
}'
```
- `ttl`: Time to live of the cache
- `statuses`: HTTP cacheable statuses
- `tags`: Used for invalidation

###### Override TTL from microservice

Set `Cache-Control: s-maxage=<seconds>` header from your microservice response to override service default TTL.

###### Override Tags from microservice

Set `X-Cache-Tags: <tag>` custom headers from your microservice response to override service default tags.

###### Cache invalidation

Any cache tagged with `<tag>` can be invalidated with:

```bash
curl --request POST 'https://<gotway>/api/cache' \
--header 'Content-Type: application/json' \
--data-raw '{
    "tags": ["<tag>"]
}'
```
You can also provide a path to be invalidated:
```bash
curl --request POST 'https://<gotway>/api/cache' \
--header 'Content-Type: application/json' \
--data-raw '{
{
    "paths": [
        {
            "servicePath": "catalog",
            "path": "/products?offset=0&limit=10"
        }
    ]
}'
```
