package execution

import (
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
)

type TestCaseExecutionResult struct {
	TestCaseDef              *model.TestCaseDef
	Executed                 bool
	TestSuiteExecutionResult *TestSuiteExecutionResult
	TestCaseReport           *report.TestCaseReport
}
