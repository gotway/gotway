# microgateway

A simple, lightweight and blazingly fast Microgateway written in Go

## Development

Install [golang/dep](https://github.com/golang/dep) for dependency management and then run:

```
$ dep ensure
$ go run .
```

## Production

```
$ docker build -t microgateway .
$ docker run -it --name microgateway -p 8080:8080 -d microgateway
```