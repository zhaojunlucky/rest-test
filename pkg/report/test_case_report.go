package report

import "github.com/zhaojunlucky/rest-test/pkg/model"

type TestCaseReport struct {
	executionTime float64
	totalTime     float64
	testCase      *model.TestCaseDef
}

func (t *TestCaseReport) GetTestDef() *model.TestCaseDef {
	return t.testCase
}

func (t *TestCaseReport) GetExecutionTime() float64 {
	return t.executionTime
}

func (t *TestCaseReport) GetTotalTime() float64 {
	return t.totalTime
}
