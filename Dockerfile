FROM golang:1.14.2

ENV WORKDIR /go/src/microgateway
RUN mkdir ${WORKDIR}
WORKDIR ${WORKDIR}
ADD .. ${WORKDIR}/

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN make install

CMD ["microgateway"]