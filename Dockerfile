FROM golang:alpine AS builder
RUN apk update && \
    apk add --no-cache curl && \
    apk add --no-cache git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/ 
RUN dep ensure && go build -o /go/bin/go-p2p
CMD ["./main"]

FROM scratch
COPY --from=builder /go/bin/go-p2p /go/bin/go-p2p
EXPOSE 3000
ENTRYPOINT ["/go/bin/go-p2p"]