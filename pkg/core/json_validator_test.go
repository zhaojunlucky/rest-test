package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONValidator_SimplePass(t *testing.T) {
	jsonValidator := NewJSONValidator()
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	validators := map[string]any{
		"or": []map[string]any{
			{"$.a": 1},
		},
	}
	err := jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c": 3},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 3},
			{"$.a": 1},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONValidator_SimpleFail(t *testing.T) {
	jsonValidator := NewJSONValidator()
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	validators := map[string]any{
		"or": []map[string]any{
			{"$.a": 2},
		},
	}
	err := jsonValidator.validate(jsonData, validators)
	if err == nil {
		t.Fatal("expect failed")
	}

	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c": 2},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err == nil {
		t.Fatal("expect failed")
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 3},
			{"$.a": 2},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err == nil {
		t.Fatal("expect failed")
	}

}

func TestJSONValidator_ComplexPass(t *testing.T) {
	jsonValidator := NewJSONValidator()
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": map[string]any{
			"d": 3,
			"e": "4",
			"f": []int{5, 6},
		},
	}

	validators := map[string]any{
		"or": []map[string]any{
			{"$.c.f": []int{5, 6}},
		},
	}
	err := jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 2},
			{"$.c.f": []int{5, 6}},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c.f": []int{5, 6}},
			{"$.c.e": "4"},
		},
	}
	err = jsonValidator.validate(jsonData, validators)
	if err != nil {
		t.Fatal(err)
	}

}

func TestJSONOperatorOr_Pass(t *testing.T) {
	jsonOP := JSONOperator{
		expectCount: 3,
		OP:          "or",
	}

	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))

	assert.Equal(t, true, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error")))
	assert.Equal(t, true, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))
	assert.Equal(t, true, jsonOP.Passed())
	assert.Nil(t, jsonOP.GetErrors())

	assert.NotNil(t, t, jsonOP.Add(nil))
}

func TestJSONOperatorOr_Fail(t *testing.T) {
	jsonOP := JSONOperator{
		expectCount: 3,
		OP:          "or",
	}

	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error")))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error")))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error")))
	assert.Equal(t, false, jsonOP.Passed())

	assert.NotNil(t, jsonOP.GetErrors())

	assert.NotNil(t, t, jsonOP.Add(nil))
}

func TestJSONOperator_And_pass(t *testing.T) {
	jsonOP := JSONOperator{
		expectCount: 3,
		OP:          "and",
	}

	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))
	assert.Equal(t, true, jsonOP.Passed())
	assert.Nil(t, jsonOP.GetErrors())

	assert.NotNil(t, t, jsonOP.Add(nil))
}

func TestJSONOperator_And_fail(t *testing.T) {
	jsonOP := JSONOperator{
		expectCount: 3,
		OP:          "and",
	}

	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(nil))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error")))
	assert.Equal(t, false, jsonOP.Passed())

	assert.Equal(t, nil, jsonOP.Add(fmt.Errorf("test error2")))
	assert.Equal(t, false, jsonOP.Passed())

	assert.NotNil(t, jsonOP.GetErrors())
	assert.NotNil(t, t, jsonOP.Add(nil))

}
