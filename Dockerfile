# Build proxy service
FROM debian:9.6-slim as builder

ARG VERSION

RUN mkdir /build
WORKDIR /build

RUN apt-get update && apt-get install -y curl && \
    curl -L "https://github.com/CyberGRX/api-connector-bulk/releases/download/v${VERSION}/api-connector-bulk-${VERSION}.linux" && \
    curl -L "https://github.com/CyberGRX/api-connector-bulk/releases/download/v${VERSION}/api-connector-bulk-${VERSION}.linux.sha" && \
    cat api-connector-bulk-$VERSION.linux.sha && \
    sha256sum --check api-connector-bulk-$VERSION.linux.sha && \
    mv api-connector-bulk-$VERSION.linux api-connector-bulk


# Pull base image.
FROM scratch

ENV GIN_MODE="release"

RUN mkdir /workingdir
WORKDIR /workingdir

COPY --from=builder /build/api-connector-bulk /workingdir/api-connector-bulk
RUN chmod +x /workingdir/api-connector-bulk

# API Connector always runs on port 8080
EXPOSE 8080

CMD ["/workingdir/api-connector-bulk"]