package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_globalReportToPoints(t *testing.T) {
	in := reports{
		Global: map[string]interface{}{
			"gcl": 2,
			"gpl": 0,
		},
	}

	result, _ := reportsIntoPoints(in, "shard1")

	assert.Len(t, result, 1)
	assert.Equal(t, "globals", result[0].Name())
	assert.Len(t, result[0].TagList(), 1)
	assert.Equal(t, "shard", result[0].TagList()[0].Key)
	assert.Equal(t, "shard1", result[0].TagList()[0].Value)

	assert.Len(t, result[0].FieldList(), 2)
	assert.Equal(t, "gcl", result[0].FieldList()[0].Key)
	assert.Equal(t, int64(2), result[0].FieldList()[0].Value)
	assert.Equal(t, "gpl", result[0].FieldList()[1].Key)
	assert.Equal(t, int64(0), result[0].FieldList()[1].Value)
}

func Test_controllersReportToPoints(t *testing.T) {
	in := reports{
		Rooms: map[room]controllerMem{
			"room1": {
				"harvesting": map[string]interface{}{
					"energy": 4.20,
				},
				"spawning": map[string]interface{}{},
			},
		},
	}

	result, _ := reportsIntoPoints(in, "shard1")

	assert.Len(t, result, 2)
	assert.Equal(t, "harvesting", result[0].Name())
	assert.Len(t, result[0].TagList(), 2)
	// are tags fixed order? it is a map...
	assert.Equal(t, "room", result[0].TagList()[0].Key)
	assert.Equal(t, "room1", result[0].TagList()[0].Value)
	assert.Equal(t, "shard", result[0].TagList()[1].Key)
	assert.Equal(t, "shard1", result[0].TagList()[1].Value)

	assert.Len(t, result[0].FieldList(), 1)
	assert.Equal(t, "energy", result[0].FieldList()[0].Key)
	assert.Equal(t, 4.20, result[0].FieldList()[0].Value)
}
