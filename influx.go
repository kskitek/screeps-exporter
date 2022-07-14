package main

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func initInfluxClient(url, token string) (influxdb2.Client, error) {
	opts := influxdb2.DefaultOptions()
	retryTimeout := time.Second * 5
	opts.SetMaxRetryTime(uint(retryTimeout.Milliseconds()))
	requestTimeout := time.Second * 5
	opts.SetHTTPRequestTimeout(uint(requestTimeout))
	warningLevel := uint(1)
	opts.SetLogLevel(warningLevel)
	opts.SetPrecision(time.Second)

	client := influxdb2.NewClientWithOptions(url, token, opts)
	_, _ = client.Health(context.Background())
	return client, nil
}

func reportsIntoPoints(reports *reports, shard string) ([]*write.Point, error) {
	var points []*write.Point

	// TODO add time to report
	time := reports.Time
	for k, v := range reports.Global {
		points = append(points, influxdb2.NewPoint(
			k,
			map[string]string{"shard": shard},
			v.(map[string]interface{}),
			time,
		))
		fmt.Println(k, v)
	}

	for controllerName, controllerMem := range reports.Controllers {
		for room, mem := range controllerMem {
			tags := map[string]string{
				"shard": shard,
				"room":  room,
			}
			fields := mem
			point := influxdb2.NewPoint(
				controllerName,
				tags,
				fields,
				time,
			)
			points = append(points, point)
		}
	}

	return points, nil
}

func writeToInflux(points []*write.Point, client influxdb2.Client, org, bucket string) error {
	api := client.WriteAPI(org, bucket)
	for _, p := range points {
		api.WritePoint(p)
	}

	api.Flush()

	// let's start with getting the first error - do all or nothing
	select {
	case err := <-api.Errors():
		return err
	default:
		return nil
	}
}
