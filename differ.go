package main

import (
	"context"
	"time"

	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type querier struct {
	Query  query
	Client *elastic.Client
	Lag    prometheus.Gauge
}

func newQuerier(q query, client *elastic.Client) querier {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "log_lag_seconds",
			Help:        "Amount of time between the most recent log and now in seconds",
			ConstLabels: prometheus.Labels{"index": q.Index},
		},
	)
	prometheus.MustRegister(gauge)

	return querier{
		Query:  q,
		Client: client,
		Lag:    gauge,
	}
}

func (q querier) Run(ctx context.Context) {

	ticker := time.Tick(time.Second * 60)
	q.getTimeDiff()

	go func() {
		for {
			select {
			case <-ticker:
				q.getTimeDiff()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (q querier) getTimeDiff() {
	termQuery := elastic.NewTermQuery(q.Query.QueryKey, q.Query.QueryValue)
	result, err := q.Client.Search().
		Index(q.Query.Index).
		Sort(q.Query.TimeField, false).
		Query(termQuery).
		From(0).Size(1).
		Do(context.Background())

	if err != nil {
		log.Errorf("error querying elasticsearch: %v", err)
		return
	}

	if len(result.Hits.Hits) == 0 {
		log.Errorf("no results found when querying index: %v", q.Query.Index)
		return
	}

	raw := string(*result.Hits.Hits[0].Source)
	timeRaw := gjson.Get(raw, q.Query.TimeField)
	timestamp, err := time.Parse(q.Query.TimeLayout, timeRaw.String())

	if err != nil {
		log.Errorf("error parsing timestamp: %v", err)
		return
	}

	timeDiff := time.Now().Sub(timestamp)
	log.Infof("time diff for index: %v is: %v\n", q.Query.Index, timeDiff)
	q.Lag.Set(timeDiff.Seconds())

}
