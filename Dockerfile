FROM golang:1.13

ENV WORKDIR /go/src/microgateway
RUN mkdir ${WORKDIR}
WORKDIR ${WORKDIR}
ADD . . ${WORKDIR}/

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v
RUN go install -v .

CMD ["microgateway"]