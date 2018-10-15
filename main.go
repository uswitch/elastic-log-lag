package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/kubernetes/staging/src/k8s.io/sample-controller/pkg/signals"

	"github.com/olivere/elastic"
)

type query struct {
	Index      string `json:"index"`
	TimeField  string `json:"timeField"`
	QueryKey   string `json:"queryKey"`
	QueryValue string `json:"queryValue"`
	TimeLayout string `json:"timeLayout"`
}

type options struct {
	configFile  string
	elasticURL  string
	elasticPass string
	elasticUser string
}

func main() {

	var opts options

	kingpin.Flag("config-file", "path to config file").Required().StringVar(&opts.configFile)
	kingpin.Flag("elastic-url", "elasticsearch url").Required().StringVar(&opts.elasticURL)
	kingpin.Parse()

	opts.elasticUser = os.Getenv("ELASTIC_USER")
	opts.elasticPass = os.Getenv("ELASTIC_PASSWORD")

	queries := make([]query, 0)

	file, err := ioutil.ReadFile(opts.configFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	json.Unmarshal(file, &queries)

	client, err := elastic.NewClient(
		elastic.SetURL(opts.elasticURL),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(opts.elasticUser, opts.elasticPass),
	)
	if err != nil {
		log.Fatalf("error creating elastic client: %v", err)
	}

	stopCh := signals.SetupSignalHandler()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricServer := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: mux,
	}

	go func() {
		if err := metricServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve metric server: %v", err)
		}
	}()

	for _, query := range queries {
		q := newQuerier(query, client)
		go q.Run(ctx)
	}

	<-stopCh
	log.Infof("shutting down")
	cancel()

}
