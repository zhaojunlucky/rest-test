package validator

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"reflect"
	"strings"
	"testing"
)

func TestJSONPathLanguage_Compile(t *testing.T) {

	lang := NewJSONPathLanguage()
	err := lang.Compile(`contain("UNIQUE constraint failed:", $.errorMessages[0].message)`)
	if err != nil {
		t.Fatal(err)
	}
	str := "{\"a\":123,\"b\":12.3}"

	var parsed map[string]interface{}
	d := json.NewDecoder(strings.NewReader(str))
	d.UseNumber()
	fmt.Println(d.Decode(&parsed))

	get, err := jsonpath.Get("$.a", parsed)
	if err != nil {
		t.Fatal(err)
	}
	valType := reflect.ValueOf(get)
	fmt.Println("@@@@")
	fmt.Println(valType.Type().String())
	fmt.Println("-------")
}
