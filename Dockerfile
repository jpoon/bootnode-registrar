FROM golang:1.9.2-alpine3.6

ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/jpoon/bootnode-registrar

# Package dependencies
RUN apk add --update --no-cache git

# Install dependency tool
RUN go get github.com/golang/dep && go install github.com/golang/dep/cmd/dep

# Copy project files
WORKDIR ${APPROOT}
COPY . ${APPROOT}

# Get dependencies and compile
RUN dep ensure
RUN go build

ENTRYPOINT [ "./bootnode-registrar" ]
EXPOSE 9898
