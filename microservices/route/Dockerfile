FROM golang:1.14.4-alpine3.12 AS builder

ENV WORKDIR /go/src/route
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY . .

RUN go build -o bin/route -v .

FROM alpine:3.12.0

COPY --from=builder /go/src/route/bin/route /app/route

CMD ["/app/route"]