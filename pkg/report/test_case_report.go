package report

import "github.com/zhaojunlucky/rest-test/pkg/model"

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
