# microgateway

A simple, lightweight and blazingly fast Microgateway written in Go

## building and running with docker

```
docker build -t go-microgateway .
docker run -it --name go-microgateway -p 8080:8080 go-microgateway
```