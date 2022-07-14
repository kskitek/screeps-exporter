package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readEmptyMemory(t *testing.T) {
	in := `{
  }
  `
	mem, err := readBody(toScreepsGzipedBody(in))

	assert.NoError(t, err)
	assert.NotNil(t, mem)
}

func Test_readSampleMemory(t *testing.T) {
	in := `
{"reports": {
  "time": "2022-07-12T21:34:58.205Z",
  "global": {
    "cpu": {
      "bucket": 5284,
      "used": 0.4894122999999979,
      "limit": 20
    },
    "gcl": {
      "level": 7,
      "progress": 2477721.6011481136
    },
    "gpl": {
      "level": 0,
      "progress": 0
    },
    "pixels": {
      "generated": 14387
    }
  },
  "harvesting": {
    "W37N49": {
      "harvesters": 0,
      "energy": 231,
      "capacity": 300
    }
  },
  "spawning": {
    "W37N49": {
      "queueLength": 2
    }
  }
}}`

	mem, err := readBody(toScreepsGzipedBody(in))

	assert.NoError(t, err)
	assert.NotNil(t, mem)
	assert.Len(t, mem.Global, 4)
	assert.Len(t, mem.Global["cpu"], 3)
	cpu := mem.Global["cpu"].(map[string]interface{})
	assert.Equal(t, 20.0, cpu["limit"])
	assert.Len(t, mem.Controllers, 2)
	assert.Len(t, mem.Controllers["harvesting"], 1)
	assert.Len(t, mem.Controllers["harvesting"]["W37N49"], 3)
	assert.Equal(t, 231.0, mem.Controllers["harvesting"]["W37N49"]["energy"])
	assert.Len(t, mem.Controllers["spawning"], 1)
	assert.Len(t, mem.Controllers["spawning"]["W37N49"], 1)

	assert.Equal(t, 1, len(mem.Controllers["harvesting"]))
}

// toScreepsGzipedBody provides input in screeps API format.
// Data field in screeps response body is in format:
// `"data": "gz:base64(gzip(memory.json))"`
// where `base64` and `gzip` are functions.
func toScreepsGzipedBody(data string) io.Reader {
	buff := bytes.NewBuffer([]byte{})
	gz := gzip.NewWriter(buff)
	_, _ = gz.Write([]byte(data))
	gz.Close()
	based := base64.StdEncoding.EncodeToString(buff.Bytes())
	in := fmt.Sprintf(`{
    "data": "gz:%s",
    "ok": 1
  }`, based)

	return bytes.NewBufferString(in)
}
