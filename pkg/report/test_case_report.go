package report

import (
	"fmt"
	"github.com/zhaojunlucky/rest-test/pkg/model"
)

type TestCaseReport struct {
	ExecutionTime float64
	TotalTime     float64
	TestCase      *model.TestCaseDef
	Error         error
	Status        string
}

func (t *TestCaseReport) GetTestDef() *model.TestCaseDef {
	return t.TestCase
}

func (t *TestCaseReport) GetExecutionTime() float64 {
	return t.ExecutionTime
}

func (t *TestCaseReport) GetTotalTime() float64 {
	return t.TotalTime
}

func (t *TestCaseReport) GetError() error {
	return t.Error
}

func (t *TestCaseReport) GetStatus() string {
	return t.Status
}

func (t *TestCaseReport) GetReportData() (map[string]any, error) {
	if t.Status == "" {
		return nil, fmt.Errorf("test case desc: %s name: %s is not executed", t.TestCase.Description, t.TestCase.Name)
	}
	return map[string]any{
		"type":          "case",
		"desc":          t.TestCase.Description,
		"name":          t.TestCase.Name,
		"executionTime": t.ExecutionTime,
		"totalTime":     t.TotalTime,
		"error":         getErrorStr(t.Error),
		"status":        t.Status,
	}, nil
}
