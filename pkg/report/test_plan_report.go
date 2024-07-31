package report

import (
	"github.com/zhaojunlucky/rest-test/pkg/model"
)

type TestPlanReport struct {
	executionTime float64
	totalTime     float64
	testPlan      *model.TestPlanDef
}

func (t *TestPlanReport) GetTestDef() *model.TestPlanDef {
	return t.testPlan
}

func (t *TestPlanReport) GetExecutionTime() float64 {
	return 0
}

func (t *TestPlanReport) GetTotalTime() float64 {
	return 0
}

func (t *TestPlanReport) GetChildren() []*TestSuiteReport {
	return nil
}
