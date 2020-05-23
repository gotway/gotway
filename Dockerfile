FROM golang:1.14.2

ENV WORKDIR /go/src/microgateway
RUN mkdir ${WORKDIR}
WORKDIR ${WORKDIR}
ADD .. ${WORKDIR}/

RUN make install

CMD ["microgateway"]