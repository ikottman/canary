FROM golang:1.17.8-alpine as builder

# system dependencies
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite=3.36.0-r0 gcc musl-dev
WORKDIR /opt/canary/

# database migrations
COPY data data
RUN sqlite3 data/measurements.db '.read data/migrations.sql'

# install dependency to speed up builds
COPY server/go.mod server/go.sum ./
RUN go get github.com/mattn/go-sqlite3@v1.14.12
RUN go get github.com/golang-jwt/jwt@v3.2.2

# compile
COPY server/ .
RUN go build

FROM alpine:latest
RUN apk add --update gcc musl-dev
WORKDIR /opt/canary/
COPY --from=builder /opt/canary/data/measurements.db .
COPY --from=builder /opt/canary/canary .
CMD ["./canary"]