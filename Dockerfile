FROM golang:1.9.2-alpine3.6 as gopher

ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/jpoon/bootnode-registrar

# Package dependencies
RUN apk add --update --no-cache git gcc libc-dev

# Install dependency tool
RUN go get github.com/golang/dep && go install github.com/golang/dep/cmd/dep

# Copy project files
WORKDIR ${APPROOT}
COPY . ${APPROOT}

# Get dependencies and compile
RUN dep ensure
RUN go build

FROM alpine:latest as final

WORKDIR /work
COPY --from=gopher ["/go/src/github.com/jpoon/bootnode-registrar/bootnode-registrar", "/work/bootnode-registrar"]

ENTRYPOINT [ "./bootnode-registrar" ]
EXPOSE 9898
