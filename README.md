# Elastic-log-lag

Elastic-log-lag queries a given index and finds the newest document, it will then find the time difference between that document's timestamp and the current time.
It exposes this as a prometheus metric of `log_lag_seconds` with a label of `index: index-name`. It also exposes the metric as a histogram of `log_lag_histogram_seconds` which is useful if you want to calculate percentiles/quantiles etc.

The rationale behind this was to allow us to easily see if we are having a problem somewhere in our logging pipeline, for example loadbalancer logs should be generated many times in a second, so if the time difference starts to become minutes we will know there's a problem somewhere in the pipeline.

It can also be used to find out how long ago a specific event happened, for example querying an index for the last occurence of `job.status: completed`.


## Flags
Elastic-log-lag takes the following flags
```
--config-file=CONFIG-FILE  path to config file
--elastic-url=ELASTIC-URL  elasticsearch url
```

It also takes `ELASTIC_USER` and `ELASTIC_PASSWORD` as `ENV` vars if you need to use basic auth when talking to Elasticsarch.
## Configuration


* `index`: name of index to query
* `queryKey`: key to query
* `queryValue`: value to find
* `timeField`: field used when sorting documents by age
* `timeLayout`: the format of your time field ([golang style](https://golang.org/src/time/format.go))

You provide this config in a json file like so:

```json
[
  {
    "index": "index-a-*",
    "timeField": "json.time",
    "queryKey": "kubernetes.container_name",
    "queryValue": "nginx",
    "timeLayout": "2006-01-02T15:04:05+00:00"
  },
  {
    "index": "index-b-*",
    "timeField": "time_local",
    "queryKey": "application",
    "queryValue": "foo",
    "timeLayout": "2006-01-02T15:04:05.000Z"
  }
]

```
