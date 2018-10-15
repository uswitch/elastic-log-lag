FROM scratch

ADD bin/elastic-log-lag elastic-log-lag

ENTRYPOINT ["/elastic-log-lag"]
