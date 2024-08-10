package model

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRestTestResponseBodyDef_Parse(t *testing.T) {
	data := `
body:
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
	respBodyDef := RestTestResponseBodyDef{}

	err = respBodyDef.Parse(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, JSON, respBodyDef.Type)
	assert.NotNil(t, respBodyDef.BodyValidator)

}
