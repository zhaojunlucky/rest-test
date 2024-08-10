package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRestTestRequestDef_Parse(t *testing.T) {
	data := `
request:
  url: github
  method: POST
  body:
    file: create_gh_server2.json
    script: |
      let data = JSON.parse(body)
      data.name = 'test2'
      return JSON.stringify(data)
`
	var jsonData map[string]any
	if err := yaml.Unmarshal([]byte(data), &jsonData); err != nil {
		t.Fatal(err)
	}
	reqDef := RestTestRequestDef{}
	err := reqDef.Parse(collection.NewMapWrapper(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "github", reqDef.URL)
	assert.Equal(t, "POST", reqDef.Method)
	assert.NotNil(t, reqDef.Body)
	assert.True(t, reqDef.Body.HasBody())
	assert.Zero(t, len(reqDef.Body.Body))
	assert.Equal(t, "create_gh_server2.json", reqDef.Body.File)
	assert.True(t, len(reqDef.Body.Script) > 0)
}
