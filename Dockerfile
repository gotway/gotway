FROM golang:1.14.4-alpine3.12 AS builder

ENV WORKDIR /go/src/gotway
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY . .

RUN go build -o bin/gotway -v .

FROM alpine:3.12.0

COPY --from=builder /go/src/gotway/bin/gotway /app/gotway

CMD [ "/app/gotway" ]