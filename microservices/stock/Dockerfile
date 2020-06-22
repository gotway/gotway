FROM golang:1.14.4-alpine3.12 AS builder

ENV WORKDIR /go/src/stock
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY . .

RUN go build -o bin/stock -v .

FROM alpine:3.12.0

COPY --from=builder /go/src/stock/bin/stock /app/stock

CMD ["/app/stock"]