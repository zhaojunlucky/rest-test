package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"math"
	"testing"
)

func TestParseRespJSONArrayBody(t *testing.T) {
	data := `
array: false
validators:
  and:
    - $.entryCount: 1
    - $.entries.length: 1
    - $.entries[0].id: ${ctx.create.resp.id}
    - or:
      - $.entries[0].name: 'test'
      - $.entries[1].name: 'test2'
    - and:
      - $.entries[0].zz: 'xxx'
      - $.entries[1].yy: 'xxx'
script: |
    if (body.entries[0].name != 'xxx') {
      throw new Error("invalid name")
    }`
	var jsonData map[string]any
	err := yaml.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		t.Fatal(err)
	}
	mapWrapper := collection.NewMapWrapper(jsonData)
	jsonBody := RestTestResponseJSONBody{}
	err = jsonBody.Parse(mapWrapper)
	if err != nil {
		t.Error(err)
	}

	assert.False(t, jsonBody.Array)
	assert.Equal(t, math.MinInt, jsonBody.Length)
}
