package report

import "github.com/zhaojunlucky/rest-test/pkg/model"

type TestSuiteReport struct {
	executionTime float64
	totalTime     float64
	testSuite     *model.TestSuiteDef
}

func (t *TestSuiteReport) GetTestDef() *model.TestSuiteDef {
	return t.testSuite
}

func (t *TestSuiteReport) GetExecutionTime() float64 {
	return t.executionTime
}

func (t *TestSuiteReport) GetTotalTime() float64 {
	return t.totalTime
}

func (t *TestSuiteReport) GetChildren() []*TestCaseReport {
	return nil
}
