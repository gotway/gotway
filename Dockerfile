FROM golang:1.14.4-alpine3.12 AS builder

ENV WORKDIR /go/src/microgateway
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY . .

RUN go build -o bin/microgateway -v .

FROM alpine:3.12.0

COPY --from=builder /go/src/microgateway/bin/microgateway /app/microgateway

CMD [ "/app/microgateway" ]