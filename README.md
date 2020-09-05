# microgateway
[![Build Status](https://travis-ci.org/gosmo-devs/microgateway.svg)](https://travis-ci.org/gosmo-devs/microgateway)

A simple, lightweight and blazingly fast API microgateway 🚀

# Features ⚡

- API composition and dynamic routing for **REST** and **gRPC** microservices
- Service discovery by registering in microgateway API
- Service health checking 
- ~10MB [Docker image](https://hub.docker.com/repository/registry-1.docker.io/gosmogolang/microgateway/tags?page=1) available for multiple architectures

[Upcoming features](https://github.com/gosmo-devs/microgateway/milestone/1) 🚧
- Centralized auth via OAuth2
- Cache management
- Rate limiting
- Microgateway API enhancements

# Service discovery 🔭

Services can be discovered in runtime by registering them in the microgateway API.

## REST

We will register [catalog](./microservices/catalog) as an example:

```
curl --request POST 'https://<microgateway>/api/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "type": "rest",
    "url": "http://<catalog>",
    "path": "catalog"
}'
```

After executing that command, our service will be available at
`https://<microgateway>/<path>`. The following endpoints will be routed through microgateway:

- `GET https://<microgateway>/catalog/products`
- `POST https://<microgateway>/catalog/product`
- `GET https://<microgateway>/catalog/product/<id>`
- `DELETE https://<microgateway>/catalog/product/<id>`
- `PUT https://<microgateway>/catalog/product/<id>`

## gRPC

We will register [route](./microservices/route) as an example:

```
curl --request POST 'https://<microgateway>/api/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "type": "grpc",
    "url": "http://<route>:<port>",
    "path": "route.Route"
}'
```

Where `route.Route` represents the package and the service name of your gRPC service. Another example could be:

[grpc.health.v1.Health](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

This will be defined in your [.proto](./microservices/route/pb/route.proto) file and will be used as path in the dynamic routing.

In this case, the RPC methods routed through microgateway will be:

- `https://<microgateway>/route.Route/GetFeature`
- `https://<microgateway>/route.Route/ListFeatures`
- `https://<microgateway>/route.Route/RecordRoute`
- `https://<microgateway>/route.Route/RouteChat`

For testing them, we have a [gRPC go client](./microservices/route/client/client.go):
```
$ cd microservices/route
$ make cli
```

# Service health checking 🚑

Microgateway will make a health probe to check that our services are responding. In other case, a `502 Bad Gateway` will be returned.

## REST

By default, the health probe will be done by requesting `http://<microservice>/health`. However, it is posible to use a custom path by specifying `healthPath` when registering.

An example of REST health endpoint is available [here](./microservices/catalog/api/api.go).

## gRPC

By default, the standard [gRPC health checking protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) is used. However, it is posible to use another one by specifying `healthPath` when registering.

An example of gRPC health checking protocol implementation can be found [here](./microservices/route/server/server.go).

# Useful scripts 🔨

### Run dev environment

```
$ ./scripts/run-dev.sh
```

Run services in your machine using [tmux](https://github.com/tmux/tmux/wiki). Additional services, like Redis, are defined in [docker-compose.dev.yml](./docker-compose.dev.yml).

### Build Docker images locally

```
$ ./scripts/build.sh
```

### Run services using local images

```
$ docker-compose up -d
```
This can be useful for testing images before pushing them to DockerHub.

### Release images for multiple architectures

```
$ ./scripts/release.sh
```
This script will be executed by TravisCI when a tag is pushed.

### Deploy to a Kubernetes cluster

```
$ ./scripts/deploy.sh <environment>
```

# Docker images 🐳

|Service|Image|
|-------|-----|
|Microgateway|[gosmogolang/microgateway](https://hub.docker.com/r/gosmogolang/microgateway)|
|Catalog|[gosmogolang/catalog](https://hub.docker.com/r/gosmogolang/catalog)|
|Stock|[gosmogolang/stock](https://hub.docker.com/r/gosmogolang/stock)|
|Route|[gosmogolang/route](https://hub.docker.com/r/gosmogolang/route)|
