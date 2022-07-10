package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")
	shard := os.Getenv("SHARD")

	memory, err := getMemory(token, shard)
	if err != nil {
		panic(err)
	}

	fmt.Println(memory)
	points, err := reportsIntoPoints(memory.Reports, shard)
	fmt.Println(points, err)
}
