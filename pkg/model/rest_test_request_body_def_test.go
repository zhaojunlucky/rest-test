package model

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestRequestBodyNil(t *testing.T) {
	d := &RestTestRequestBodyDef{}
	err := d.Parse(nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, d.HasBody())
}

func TestRequestBodyString(t *testing.T) {
	d := &RestTestRequestBodyDef{}
	err := d.Parse("hello")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, d.HasBody())
	assert.Equal(t, "hello", d.Body)
}

func TestRequestBodyInt(t *testing.T) {
	d := &RestTestRequestBodyDef{}
	err := d.Parse(9)
	if err == nil {
		t.Fatal("expect error")
	}
	assert.False(t, d.HasBody())
}

func TestRequestBodyMap(t *testing.T) {
	data := `
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
	d := &RestTestRequestBodyDef{}
	err := d.Parse(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, d.HasBody())
	assert.Zero(t, len(d.Body))
	assert.Equal(t, "create_gh_server2.json", d.File)
	assert.True(t, len(d.Script) > 0)
}
