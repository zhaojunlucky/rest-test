package validator

import (
	"errors"
	"fmt"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"strings"
)

const (
	And = "and"
	Or  = "or"
)

type JSONPathNode interface {
	Validate(js core.JSEnvExpander, v any) error
}

func IsLogicOperator(op string) bool {
	return strings.EqualFold(op, And) || strings.EqualFold(op, Or)
}

type JSONOperator struct {
	expectCount  int
	OP           string
	errors       []error
	successCount int
}

func (j *JSONOperator) Add(err error) error {
	if err != nil {
		j.errors = append(j.errors, err)
	} else {
		j.successCount++
	}

	if len(j.errors)+j.successCount > j.expectCount {
		return fmt.Errorf("failed to verify json path, got %d results, want %d", j.successCount+len(j.errors), j.expectCount)
	}
	return nil
}

func (j *JSONOperator) GetErrors() error {
	return errors.Join(j.errors...)
}

func (j *JSONOperator) Passed() bool {
	if strings.EqualFold(j.OP, Or) {
		return j.successCount > 0
	} else if strings.EqualFold(j.OP, And) {
		return j.successCount == j.expectCount
	}
	return false
}
