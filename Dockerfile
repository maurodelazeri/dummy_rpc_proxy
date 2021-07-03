FROM golang:1.16.4
WORKDIR /go/dummy_rpc_proxy
COPY . .
RUN go build -o dummy_rpc_proxy
ENTRYPOINT ["./dummy_rpc_proxy"]
