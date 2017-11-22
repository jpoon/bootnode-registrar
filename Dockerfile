FROM golang:1.9.2

ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/jpoon/bootnode-registrar

ADD . ${APPROOT}
WORKDIR ${APPROOT}

RUN go get -u github.com/golang/dep/... \
    && dep ensure 
RUN make

CMD ["./bootnode-registrar"]
