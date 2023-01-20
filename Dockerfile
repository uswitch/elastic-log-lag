FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/elastic-log-lag elastic-log-lag

USER nonroot:nonroot

ENTRYPOINT ["/elastic-log-lag"]
