package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRestTestResponseDefParse(t *testing.T) {
	data := `
response:
  code: 200
  contentType: application/json
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
	respBodyDef := RestTestResponseDef{}
	err = respBodyDef.Parse(collection.NewMapWrapper(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 200, respBodyDef.Code)
	assert.Equal(t, "application/json", respBodyDef.ContentType)
	assert.NotNil(t, respBodyDef.Body)
	assert.NotNil(t, respBodyDef.Body.BodyValidator)
}
