package main

import (
	"log"
	"net/http"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/kelseyhightower/envconfig"
)

var client influxdb2.Client

// TODO this is a lazy script.. not a production grade program..
// please improve it someday... ;)

type screepsConfig struct {
	Token string
	Shard string
}

type influxConfig struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

func main() {
	var screepsConfig screepsConfig
	err := envconfig.Process("", &screepsConfig)
	if err != nil {
		panic(err)
	}

	var influxConfig influxConfig
	err = envconfig.Process("INFLUX", &influxConfig)
	if err != nil {
		panic(err)
	}

	c, err := initInfluxClient(influxConfig.URL, influxConfig.Token)
	if err != nil {
		log.Panic(err)
	}
	client = c

	handler := handler{
		influx: c,
		sc:     screepsConfig,
		ic:     influxConfig,
	}

	http.HandleFunc("/", handler.handleEvent)
	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

type handler struct {
	influx influxdb2.Client
	sc     screepsConfig
	ic     influxConfig
}

func (h handler) handleEvent(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Reading memory from shard: %s\n", h.sc.Shard)
	reports, err := getMemory(h.sc.Token, h.sc.Shard)
	if err != nil {
		panic(err)
	}
	log.Printf("Got memory from shard: %s\n", h.sc.Shard)

	points, err := reportsIntoPoints(reports, h.sc.Shard)

	writeToInflux(points, client, h.ic.Org, h.ic.Bucket)
	log.Printf("Written %d points from shard %s\n", len(points), h.sc.Shard)
}
