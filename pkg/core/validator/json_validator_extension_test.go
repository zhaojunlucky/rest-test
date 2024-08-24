package validator

import (
	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"testing"
)

func TestInFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`2 in [1, 2, 3]`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}
}

func TestContainsFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`contain("hello", "helloworld")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`contain(2, [1, 2, 3])`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`contain("a", $.obj)`, map[string]any{"obj": map[any]any{"a": 1, "b": 2, "c": 3}})
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}
}

func TestIntFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`int(123)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123 {
		t.Fatal("expect 123")
	}

	val, err = lang.Evaluate(`int("123")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123 {
		t.Fatal("expect 123")
	}

	val, err = lang.Evaluate(`int(123.45)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123 {
		t.Fatal("expect 123")
	}
}

func TestFloatFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`float(123)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123.0 {
		t.Fatal("expect 123")
	}

	val, err = lang.Evaluate(`float("123.1")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123.1 {
		t.Fatal("expect 123.1")
	}

	val, err = lang.Evaluate(`float(123.45)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 123.45 {
		t.Fatal("expect 123.45")
	}
}

func TestStringFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`string(123)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != "123.000000" {
		t.Fatal("expect 123.000000")
	}

	val, err = lang.Evaluate(`string("123")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != "123" {
		t.Fatal("expect 123")
	}

	val, err = lang.Evaluate(`string(123.45)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != "123.450000" {
		t.Fatal("expect 123.450000")
	}
}

func TestLenFunc(t *testing.T) {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`len("123")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != 3 {
		t.Fatal("expect 3")
	}

	val, err = lang.Evaluate(`len($.arr)`, map[string]any{"arr": []int{1, 2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	if val != 3 {
		t.Fatal("expect 3")
	}

	val, err = lang.Evaluate(`len($.obj)`, map[string]any{"obj": map[string]any{"a": 1, "b": 2, "c": 3}})
	if err != nil {
		t.Fatal(err)
	}
	if val != 3 {
		t.Fatal("expect 3")
	}
}

func TestBoolFunc(t *testing.T) {

	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lang := gval.Full(allLangs...)

	val, err := lang.Evaluate(`bool("true")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`bool("false")`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != false {
		t.Fatal("expect false")
	}

	val, err = lang.Evaluate(`bool(24.1)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`bool(0)`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if val != false {
		t.Fatal("expect false")
	}

	val, err = lang.Evaluate(`bool($.arr)`, map[string]any{"arr": []int{1, 2, 3}})
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`bool($.arr)`, map[string]any{"arr": []int{}})
	if err != nil {
		t.Fatal(err)
	}
	if val != false {
		t.Fatal("expect false")
	}

	val, err = lang.Evaluate(`bool($.obj)`, map[string]any{"obj": map[string]any{"a": 1, "b": 2, "c": 3}})
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Fatal("expect true")
	}

	val, err = lang.Evaluate(`bool($.obj)`, map[string]any{"obj": map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if val != false {
		t.Fatal("expect false")
	}

}
