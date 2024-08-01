package report

import (
	"github.com/zhaojunlucky/rest-test/pkg/model"
)

type TestPlanReport struct {
	ExecutionTime     float64
	TotalTime         float64
	TestPlan          *model.TestPlanDef
	testSuitesReports []*TestSuiteReport
	Error             error
	Status            string
}

func (t *TestPlanReport) GetTestDef() *model.TestPlanDef {
	return t.TestPlan
}

func (t *TestPlanReport) GetExecutionTime() float64 {
	return t.ExecutionTime
}

func (t *TestPlanReport) GetTotalTime() float64 {
	return t.TotalTime
}

func (t *TestPlanReport) GetChildren() []*TestSuiteReport {
	return t.testSuitesReports
}

func (t *TestPlanReport) AddTestSuiteReport(testSuiteReport *TestSuiteReport) {
	t.testSuitesReports = append(t.testSuitesReports, testSuiteReport)
}

func (t *TestPlanReport) GetError() error {
	return t.Error
}

func (t *TestPlanReport) GetStatus() string {
	return t.Status
}
