package report

import "github.com/zhaojunlucky/rest-test/pkg/model"

type TestSuiteReport struct {
	ExecutionTime   float64
	TotalTime       float64
	TestSuite       *model.TestSuiteDef
	Error           error
	testCaseReports []*TestCaseReport
	Status          string
}

func (t *TestSuiteReport) GetTestDef() *model.TestSuiteDef {
	return t.TestSuite
}

func (t *TestSuiteReport) GetExecutionTime() float64 {
	return t.ExecutionTime
}

func (t *TestSuiteReport) GetTotalTime() float64 {
	return t.TotalTime
}

func (t *TestSuiteReport) GetChildren() []*TestCaseReport {
	return t.testCaseReports
}

func (t *TestSuiteReport) AddTestCaseReport(testCaseReport *TestCaseReport) {
	t.testCaseReports = append(t.testCaseReports, testCaseReport)
}

func (t *TestSuiteReport) GetError() error {
	return t.Error
}

func (t *TestSuiteReport) GetStatus() string {
	return t.Status
}
