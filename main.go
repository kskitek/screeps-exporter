package main

import (
	"log"
	"os"

	"github.com/influxdata/influxdb-client-go/v2"
)

var client influxdb2.Client

// TODO this is a lazy script.. not a production grade program..
// please improve it someday... ;)

func main() {
	token := os.Getenv("TOKEN")
	shard := os.Getenv("SHARD")

	influxURL := os.Getenv("INFLUX_URL")
	influxToken := os.Getenv("INFLUX_TOKEN")
	influxOrg := os.Getenv("INFLUX_ORG")
	influxBucket := os.Getenv("INFLUX_BUCKET")
	c, err := initInfluxClient(influxURL, influxToken)
	if err != nil {
		log.Panic(err)
	}
	client = c

	log.Printf("Reading memory from shard: %s\n", shard)
	reports, err := getMemory(token, shard)
	if err != nil {
		panic(err)
	}
	log.Printf("Got memory from shard: %s\n", shard)

	points, err := reportsIntoPoints(reports, shard)

	writeToInflux(points, client, influxOrg, influxBucket)
	log.Printf("Written %d points from shard %s\n", len(points), shard)
}
