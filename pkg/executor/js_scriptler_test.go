package executor

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/env"
	"reflect"
	"testing"
)

func TestJSScriptlerExpand(t *testing.T) {
	envs := env.NewReadWriteEnv(env.NewOSEnv(), map[string]string{
		"APP_ID":   "1",
		"APP_NAME": "rest-tester",
	})
	js, err := NewJSScriptler(envs, NewTestSuiteCaseContext())
	if err != nil {
		t.Fatal(err)
	}
	actual, err := js.Expand("hello ${env.APP_ID}")
	if err != nil {
		t.Fatal(err)
	}
	if actual != "hello 1" {
		t.Fatalf("expect 'hello 1', but got %s", actual)
	}

	dyMap := map[string]string{
		"Name": "rest-tester ${env.APP_NAME}",
		"ID":   "id-${env.APP_ID}",
	}

	actualMap, err := js.ExpandMap(dyMap)

	if err != nil {
		t.Fatal(err)
	}
	if actualMap["Name"] != "rest-tester rest-tester" {
		t.Fatalf("expect 'rest-tester rest-tester', but got %s", actualMap["Name"])
	}
	if actualMap["ID"] != "id-1" {
		t.Fatalf("expect 'id-1', but got %s", actualMap["ID"])
	}
}

func TestJSScriptler_RunScript(t *testing.T) {
	envs := env.NewReadWriteEnv(env.NewOSEnv(), map[string]string{
		"APP_ID":   "1",
		"APP_NAME": "rest-tester",
	})
	js, err := NewJSScriptler(envs, NewTestSuiteCaseContext())
	if err != nil {
		t.Fatal(err)
	}
	script := `
		var name = env.APP_NAME
		var id = env.APP_ID
		name + "----" + id
`
	actual, err := js.ExpandScript(script)
	if err != nil {
		t.Fatal(err)
	}
	if actual != "rest-tester----1" {
		t.Fatalf("expect 'hello 1', but got %s", actual)
	}
}

func TestJSScriptler_RunScriptWithBody(t *testing.T) {
	envs := env.NewReadWriteEnv(env.NewOSEnv(), map[string]string{
		"APP_ID":   "1",
		"APP_NAME": "rest-tester",
	})
	js, err := NewJSScriptler(envs, NewTestSuiteCaseContext())
	if err != nil {
		t.Fatal(err)
	}
	script := `
		var name = env.APP_NAME
		var id = env.APP_ID
		name + "----" + id + "----" + body
`
	actual, err := js.ExpandScriptWithBody(script, "body")
	if err != nil {
		t.Fatal(err)
	}
	if actual != "rest-tester----1----body" {
		t.Fatalf("expect 'hello 1', but got %s", actual)
	}
}

func TestAA(t *testing.T) {
	var a any = "1"
	var b any = float64(1)
	ta := reflect.ValueOf(a)
	tb := reflect.ValueOf(b)

	fmt.Println(ta.CanConvert(tb.Type()))
	fmt.Println(tb.CanConvert(ta.Type()))

}
