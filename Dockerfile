FROM golang:1.17.8-alpine as builder
WORKDIR /opt/canary/
COPY src/ .
RUN go build

FROM alpine:latest
WORKDIR /opt/canary/
COPY --from=builder /opt/canary/canary .
CMD ["./canary"]