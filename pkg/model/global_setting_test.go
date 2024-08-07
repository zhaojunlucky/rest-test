package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobalSetting_With(t *testing.T) {

	gs1 := GlobalSetting{
		Headers: map[string]string{
			"key": "value",
		},
		DataDir: "dataDir",
	}

	gs2 := GlobalSetting{
		Headers: map[string]string{
			"key1": "value1",
		},
		DataDir: "dataDir2",
	}

	gs3 := gs1.With(&gs2)

	assert.Equal(t, gs3.DataDir, gs2.DataDir)
	assert.Equal(t, gs3.Headers["key1"], gs2.Headers["key1"])
	assert.Equal(t, gs3.Headers["key"], gs1.Headers["key"])
}
