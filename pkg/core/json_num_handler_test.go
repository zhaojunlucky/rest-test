package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConvertObj(t *testing.T) {
	jsonStr := `{"entryCount":2,"entries":[{"id":1,"createdAt":"2024-08-24T20:06:44.151871+08:00","updatedAt":"2024-08-24T20:06:44.151871+08:00","web":"https://github.com","api":"https://api.github.com/v3","name":"github"},{"id":2,"createdAt":"2024-08-24T20:07:41.286469+08:00","updatedAt":"2024-08-24T20:07:41.286469+08:00","web":"https://github1.com","api":"https://api.github.com/v4","name":"test2"}]}`

	jsonDecoder := json.NewDecoder(strings.NewReader(jsonStr))
	jsonDecoder.UseNumber()

	var obj map[string]any
	err := jsonDecoder.Decode(&obj)
	if err != nil {
		t.Fatal(err)
	}

	newObj := ConvertObj(obj)
	t.Log(newObj)
}
