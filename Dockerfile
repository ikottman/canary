FROM golang:1.17.8-alpine as builder

# system dependencies
RUN apk update && apk upgrade
RUN apk add --no-cache gcc musl-dev
WORKDIR /opt/canary/

# install dependency to speed up builds
COPY server/go.mod server/go.sum ./
RUN go get github.com/mattn/go-sqlite3@v1.14.12
RUN go get github.com/golang-jwt/jwt@v3.2.2

# compile
COPY server/ .
RUN go build

FROM alpine:latest
ENV CANARY_PUBLIC_KEY LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUE4ZDI2VU9vTm9ZTGJoY2g5MFVYbwpWTzh2WnFVd2tCSVVjM0pmRFMrY2VSN2hsMkZraG00QmhPZXBnbjZpMEdmRXRDWUNhTVArS08vMHJwQ2J3SjN2Cmt0c1F0TEdoYWE4NTRkMDM0dVZDQmhRZmdlZkpXLzZaQkEyc2RJTjdKK1Z1NFduL2kyY2dQR1dPTHUyN1Zxbm8KamF3SWdxOGQ0Z1hWY05jWWZESmdqUXM2V1V6dHZ4RVU2SG94TnFDT2liZTduaVQwMUdNUzVKWVRqd2dIcnVkQgpVeml3MXAyd1pOL1lYQ3VFWmJSbzhsR1VnWDVJQk9kUGh1TVBnbGdwdzROSElWaHE1a2dGdWI2WW5wZmNNcnBZCjZmOTZpa21oQ051VWhubWJsblhscm0vSnkwbnVaNFNweFJudmcrUmxJaGpjTHIrQnNpdDN4Rzh1R08rU2F6bjgKY1BZNHpHd05Nby9oTFhkM3dISkFQUitTdHgyc3hyVXJZdkpxWjgrZ2xoYVRycVVWQjJQd2ptMkZMQm53RXd2Tgo0Z2NNRnY4V2Y4T1lzdThYbjJZMndPWEpyQjd4ZVRNMWJFYVl0R1lqSjgvMDF3aVdnMUZoR09TNXNXdURWZVZMClRHVWJ4dTQ5SHo0enpRMWVNaklNUE5pYzNKMmNFMUpBVm4xU0Fpd3RmcDhwbGU2SkxoVTMzUFFhWEFrMlEyY0gKdHVlSVNVWXc5bEoybFdrOFpRaTY1bHFvQXZMR2p6RXBGNUlNVFM2VlZlUlBjY0hrUkhCSlJOWTlUT0xSVmVoUApmamRLalhEQTMxRUY2OWFqV1AvR0Z6NWJLYTJiR0NDWE5Sbmx2NTN2RW1TNkd0Ynk2OGFjeHliQkttOWdFaVhqCjdCWmRpOE9RWnlXT0NCWVRxSm5HbUdNQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQoKCg==
RUN apk add --update gcc musl-dev sqlite=3.36.0-r0 tzdata
WORKDIR /opt/canary/
COPY --from=builder /opt/canary/migrations.sh /opt/canary/migrations.sh
COPY --from=builder /opt/canary/migrations.sql /opt/canary/migrations.sql
COPY --from=builder /opt/canary/canary /opt/canary/canary
CMD ["/bin/sh", "/opt/canary/migrations.sh"]