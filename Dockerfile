FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/elastic-log-lag-linux-amd64 elastic-log-lag

USER nonroot:nonroot

ENTRYPOINT ["/elastic-log-lag"]
