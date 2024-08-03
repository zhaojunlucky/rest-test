package execution

import "github.com/zhaojunlucky/rest-test/pkg/model"

type TestCaseExecutionResult struct {
	TestCaseDef *model.TestCaseDef
	Executed    bool
}
