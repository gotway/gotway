FROM golang:1.14.4-alpine3.12 AS builder

ENV WORKDIR /go/src/catalog
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY . .

RUN go build -o bin/catalog -v .

FROM alpine:3.12.0

COPY --from=builder /go/src/catalog/bin/catalog /app/catalog

CMD ["/app/catalog"]