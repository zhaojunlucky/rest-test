package execution

import (
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
)

type TestPlanExecutionResult struct {
	TestPlanDef               *model.TestPlanDef
	TestPlanReport            *report.TestPlanReport
	TestSuiteExecResults      []*TestSuiteExecutionResult
	NamedTestSuiteExecResults map[string]*TestSuiteExecutionResult
	Executed                  bool
}

func (r TestPlanExecutionResult) HasNamed(name string) bool {
	_, ok := r.NamedTestSuiteExecResults[name]
	return ok
}

func (r TestPlanExecutionResult) GetNamed(name string) *TestSuiteExecutionResult {
	return r.NamedTestSuiteExecResults[name]
}

func (r TestPlanExecutionResult) AddTestSuiteExecResults(suite *TestSuiteExecutionResult) {
	r.TestSuiteExecResults = append(r.TestSuiteExecResults, suite)
	if suite.TestSuiteDef.Name != "" {
		r.NamedTestSuiteExecResults[suite.TestSuiteDef.Name] = suite
	}

}
