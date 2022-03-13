FROM golang:1.17.8-alpine as builder

# system dependencies
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite=3.36.0-r0
WORKDIR /opt/canary/

# database migrations
COPY data data
RUN sqlite3 data/metrics.db '.read data/migrations.sql'

# compile
COPY src/ .
RUN go build

FROM alpine:latest
WORKDIR /opt/canary/
COPY --from=builder /opt/canary/data/metrics.db .
COPY --from=builder /opt/canary/canary .
CMD ["./canary"]