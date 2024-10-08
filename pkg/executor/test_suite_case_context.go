package executor

import (
	"fmt"
	"github.com/zhaojunlucky/rest-test/pkg/execution"
)

type TestSuiteCaseContext struct {
	CaseResult map[string]map[string]any
}

func NewTestSuiteCaseContext() *TestSuiteCaseContext {
	return &TestSuiteCaseContext{
		CaseResult: make(map[string]map[string]any),
	}
}

func (t *TestSuiteCaseContext) Add(caseResult *execution.TestCaseExecutionResult, body any) error {
	if !caseResult.Executed {
		return fmt.Errorf("case %s not executed", caseResult.TestCaseDef.Name)
	}
	if len(caseResult.TestCaseDef.Name) <= 0 {
		return nil
	}
	t.CaseResult[caseResult.TestCaseDef.Name] = map[string]any{
		"requestDef":  caseResult.TestCaseDef.Request,
		"responseDef": caseResult.TestCaseDef.Response,
		"resp":        body,
	}
	return nil
}
