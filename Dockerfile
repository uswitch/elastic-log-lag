FROM scratch

ADD elastic-log-lag eslag

ENTRYPOINT ["/eslag"]
