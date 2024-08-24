package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.NotNil(t, jsonOP.GetErrors())

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
