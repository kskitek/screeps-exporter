package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type room = string
type controllerName = string
type controllerMem = map[room]map[string]interface{}

type memory struct {
	Reports map[string]interface{} `json:"reports"`
}

type reports struct {
	Time        time.Time
	Global      map[string]interface{}
	Controllers map[controllerName]controllerMem
}

func getMemory(token, shard string) (*reports, error) {
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

func readBody(r io.Reader) (*reports, error) {
	var body response
	err := json.NewDecoder(r).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode body: %w", err)
	}

	if len(body.Data) < 3 {
		return nil, fmt.Errorf("returned body is malformed. Body: %s", body.Data)
	}

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

	return memoryToReports(mem)
}

func memoryToReports(mem memory) (*reports, error) {
	// TODO think about better reports structure.. this is really nasty...
	reports := &reports{}

	if timeString, ok := mem.Reports["time"].(string); ok {
		t, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			return nil, err
		}
		reports.Time = t
	}
	if global, ok := mem.Reports["global"]; ok {
		reports.Global = global.(map[string]interface{})
	}

	reports.Controllers = make(map[controllerName]controllerMem)
	for controller, v := range mem.Reports {
		if controller == "time" || controller == "global" {
			continue
		}

		for room, mem := range v.(map[string]interface{}) {
			if reports.Controllers[controller] == nil {
				reports.Controllers[controller] = make(controllerMem)
			}
			reports.Controllers[controller][room] = mem.(map[string]interface{})
		}
	}

	return reports, nil
}

type response struct {
	Data string `json:"data"`
	OK   int    `json:"ok"`
}
