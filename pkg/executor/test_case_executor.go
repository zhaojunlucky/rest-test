package executor

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/core/execution"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"time"
)

type TestCaseExecutor struct {
}

func (t *TestCaseExecutor) Execute(ctx *core.RestTestContext, env env.Env, global *model.GlobalSetting, testDef *model.TestCaseDef, testCaseExecResult *execution.TestCaseExecutionResult) (*report.TestCaseReport, error) {

	testCaseReport := report.TestCaseReport{
		TestCase: testDef,
	}

	start := time.Now()
	testCaseReport.ExecutionTime = time.Since(start).Seconds()
	return &testCaseReport, nil
}

func (t *TestCaseExecutor) Prepare(ctx *execution.TestSuiteExecutionResult, def model.TestCaseDef) error {
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test case %s", def.Name)
	}

	testCaseExecResult := &execution.TestCaseExecutionResult{
		TestCaseDef: &def,
	}
	ctx.AddTestCaseExecResults(testCaseExecResult)
	return nil
}