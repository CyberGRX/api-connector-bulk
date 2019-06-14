# Build proxy service
FROM golang:alpine AS builder

ARG VERSION

RUN mkdir /build
WORKDIR /build

# Create appuser.
RUN adduser -D -g '' grx_appuser

RUN apk update && apk add --no-cache curl ca-certificates tzdata && \
    update-ca-certificates && \
    echo "Fetching https://github.com/CyberGRX/api-connector-bulk/releases/download/v${VERSION}/api-connector-bulk-${VERSION}" && \
    curl -L -O "https://github.com/CyberGRX/api-connector-bulk/releases/download/v${VERSION}/api-connector-bulk-${VERSION}.linux" && \
    curl -L -O "https://github.com/CyberGRX/api-connector-bulk/releases/download/v${VERSION}/api-connector-bulk-${VERSION}.linux.sha" && \
    sha256sum -c "api-connector-bulk-${VERSION}.linux.sha" && \
    mv "api-connector-bulk-${VERSION}.linux" api-connector-bulk && \
    chmod +x api-connector-bulk


# Pull base image.
FROM scratch

ENV GIN_MODE="release"
ENV HOST=""
ENV PORT="8080"

WORKDIR /workingdir

# Copy user, timezone info and SSL configuration from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable
COPY --from=builder /build/api-connector-bulk /workingdir/api-connector-bulk

# Use an unprivileged user.
USER grx_appuser

# API Connector always runs on port 8080
EXPOSE 8080

CMD ["/workingdir/api-connector-bulk"]