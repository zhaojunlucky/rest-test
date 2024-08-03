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

func (t *TestCaseExecutor) Execute(ctx *core.RestTestContext, env env.Env, global *model.GlobalSetting, testCaseExecResult *execution.TestCaseExecutionResult) *report.TestCaseReport {

	//if err != nil {
	//	log.Errorf("test case %s failed, error: %v", testCaseDef.Name, err)
	//} else {
	//	log.Infof("test case %s passed", testCaseDef.Name)
	//}
	testCaseReport := report.TestCaseReport{
		TestCase: testCaseExecResult.TestCaseDef,
	}

	testCaseExecResult.TestCaseReport = &testCaseReport

	start := time.Now()
	testCaseReport.ExecutionTime = time.Since(start).Seconds()
	return &testCaseReport
}

func (t *TestCaseExecutor) Prepare(ctx *execution.TestSuiteExecutionResult, def model.TestCaseDef) error {
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test case %s", def.Name)
	}

	testCaseExecResult := &execution.TestCaseExecutionResult{
		TestCaseDef:              &def,
		TestSuiteExecutionResult: ctx,
	}
	ctx.AddTestCaseExecResults(testCaseExecResult)
	return nil
}

func (t *TestCaseExecutor) Validate(result *execution.TestCaseExecutionResult) error {

	return nil
}
