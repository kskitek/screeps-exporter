package main

import (
	"fmt"
)

func reportsIntoPoints(reports reports, shard string) (interface{}, error) {
	fmt.Println(shard, reports.Global)

	for room, controllerMem := range reports.Rooms {
		for controllerName, mem := range controllerMem {
			fmt.Println(shard, room, controllerName, mem)
		}
	}
	return nil, nil
}
