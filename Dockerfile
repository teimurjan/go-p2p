FROM golang:alpine
RUN apk update && \
    apk add --no-cache curl && \
    apk add --no-cache git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . $GOPATH/src/github.com/teimurjan/go-p2p
WORKDIR $GOPATH/src/github.com/teimurjan/go-p2p 
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/go-p2p
EXPOSE 3000/udp
ENTRYPOINT ["/go/bin/go-p2p"]