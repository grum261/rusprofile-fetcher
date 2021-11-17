FROM golang:1.17.3-alpine3.14 AS builder

WORKDIR /build/

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \ 
    rusprofile-fetcher/cmd/gateway

#-

FROM alpine:3.14 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/

ENV PATH=/api/bin/:$PATH

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/gateway ./bin/gateway
COPY --from=builder /build/internal/rpc_server/ ./internal/rpc_server

EXPOSE 8000

CMD [ "gateway" ]