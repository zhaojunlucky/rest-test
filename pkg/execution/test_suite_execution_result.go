package execution

import (
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
)

type TestSuiteExecutionResult struct {
	TestSuiteDef              *model.TestSuiteDef
	TestCasesExecResults      []*TestCaseExecutionResult
	NamedTestCasesExecResults map[string]*TestCaseExecutionResult
	Executed                  bool
	TestSuiteReport           *report.TestSuiteReport
	TestPlanExecutionResult   *TestPlanExecutionResult
}

func (r *TestSuiteExecutionResult) HasNamed(name string) bool {
	_, ok := r.NamedTestCasesExecResults[name]
	return ok
}

func (r *TestSuiteExecutionResult) AddTestCaseExecResults(result *TestCaseExecutionResult) {

	r.TestCasesExecResults = append(r.TestCasesExecResults, result)
	if result.TestCaseDef.Name != "" {
		r.NamedTestCasesExecResults[result.TestCaseDef.Name] = result
	}
}
