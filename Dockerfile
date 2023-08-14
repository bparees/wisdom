FROM registry.redhat.io/rhel8/go-toolset:1.18 AS builder
COPY . .
RUN go build -o /tmp/wisdom ./cmd/wisdom

# use UBI instead of scratch as an easy way to get certificates.
FROM registry.redhat.io/ubi8/ubi:latest AS base
COPY --from=builder /tmp/wisdom /wisdom
ENTRYPOINT ["/wisdom"]
EXPOSE 8443
