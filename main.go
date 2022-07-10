package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

func getMemory(token, shard string) (*memory, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodGet, "https://screeps.com/api/user/memory", nil)
	q := req.URL.Query()
	q.Add("_token", token)
	q.Add("shard", shard)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error when getting memory: %d", resp.StatusCode)
	}

	return readBody(resp.Body)
}

func readBody(r io.Reader) (*memory, error) {
	var body response
	err := json.NewDecoder(r).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode body: %w", err)
	}

	// compress/zlib
	data, err := base64.StdEncoding.DecodeString(body.Data[3:])
	if err != nil {
		return nil, fmt.Errorf("unable to base64 decode memory: %w", err)
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("unable to create gzip reader: %w", err)
	}

	b, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("unable to read gziped data: %w", err)
	}

	var mem memory
	err = json.Unmarshal(b, &mem)
	if err != nil {
		return nil, fmt.Errorf("unable to read json memory: %w", err)
	}

	return &mem, nil
}

func reportsIntoPoints(reports reports, shard string) (interface{}, error) {
	fmt.Println(shard, reports.Global)

	for room, controllerMem := range reports.Rooms {
		for controllerName, mem := range controllerMem {
			fmt.Println(shard, room, controllerName, mem)
		}
	}
	return nil, nil
}

type response struct {
	Data string `json:"data"`
	OK   int    `json:"ok"`
}

type room = string
type controllerMem = map[string]interface{}

type memory struct {
	Reports reports `json:"reports"`
}

type reports struct {
	Global map[string]interface{} `json:"global"`
	Rooms  map[room]controllerMem `json:"rooms"`
}