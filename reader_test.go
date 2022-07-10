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
	in := `{
    "reports": {
      "global": {
        "gcl": 8,
        "gpl": 1,
        "bucket": 8000,
        "limit": 20,
        "used": 0.8
      },
      "rooms": {
        "W12N02": {
          "spawning": {
            "queue": 1
          },
          "harversing": {
            "creeps": 4,
            "gathered": 100
          },
          "room": {
            "lvl": 7,
            "maxEnergy": 300,
            "energy": 200
          }
        }
      }
    }
  }
  `

	mem, err := readBody(toScreepsGzipedBody(in))

	assert.NoError(t, err)
	assert.NotNil(t, mem)
	assert.Equal(t, 1, len(mem.Reports.Rooms))
	assert.Equal(t, 3, len(mem.Reports.Rooms["W12N02"]))
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
