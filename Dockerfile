#
# Build go project
#
FROM golang:1.13-alpine as go-builder

WORKDIR /forward-proxy

COPY . .

RUN apk add -u -t build-tools curl git && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy main.go && \
    apk del build-tools && \
    rm -rf /var/cache/apk/*

#
# Runtime container
#
FROM alpine:latest  

WORKDIR /app

RUN apk --no-cache add ca-certificates bash curl

COPY --from=go-builder /forward-proxy/proxy .

ENTRYPOINT ["./proxy"]
