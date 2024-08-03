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

type TestSuiteExecutor struct {
	testCaseExecutor *TestCaseExecutor
}

func (t *TestSuiteExecutor) Execute(ctx *core.RestTestContext, environ env.Env, global *model.GlobalSetting, testSuiteExecResult *execution.TestSuiteExecutionResult) *report.TestSuiteReport {

	testSuiteExecResult.TestSuiteReport = &report.TestSuiteReport{
		TestSuite: testSuiteExecResult.TestSuiteDef,
	}
	suiteReport := testSuiteExecResult.TestSuiteReport
	testSuiteDef := testSuiteExecResult.TestSuiteDef

	//if err != nil {
	//	log.Errorf("test suite %s failed, error: %v", testSuiteReport.TestSuite.Name, err)
	//} else {
	//	log.Infof("test suite %s passed", testSuiteReport.TestSuite.Name)
	//}
	if testSuiteDef.Enabled == false {
		suiteReport.Error = fmt.Errorf("test suite %s is disabled", testSuiteDef.Name)
		suiteReport.Status = report.Skipped
		return suiteReport
	}

	if len(testSuiteDef.Depends) > 0 {
		for _, dep := range testSuiteDef.Depends {
			if !testSuiteExecResult.TestPlanExecutionResult.HasNamed(dep) {
				suiteReport.Error = fmt.Errorf("depends test suite %s not found", dep)
				suiteReport.Status = report.DependencyError
				return suiteReport
			}
			depTSExecResult := testSuiteExecResult.TestPlanExecutionResult.NamedTestSuiteExecResults[dep]

			if !depTSExecResult.Executed {
				suiteReport.Error = fmt.Errorf("depends test suite %s not executed", dep)
				suiteReport.Status = report.DependencyError
				return suiteReport
			}

			if !depTSExecResult.TestSuiteReport.HasPassed() {
				suiteReport.Error = fmt.Errorf("depends test suite %s failed", dep)
				suiteReport.Status = report.DependencyError
				return suiteReport
			}
		}
	}

	start := time.Now()
	for _, testCaseExecResult := range testSuiteExecResult.TestCasesExecResults {
		newEnv := env.NewReadWriteEnv(environ, testSuiteDef.Environment)
		testCaseReport := t.testCaseExecutor.Execute(ctx, newEnv, global, testCaseExecResult)

		suiteReport.ExecutionTime += testCaseReport.ExecutionTime
		suiteReport.AddTestCaseReport(testCaseReport)
	}
	suiteReport.TotalTime = time.Since(start).Seconds()
	return suiteReport
}

func (t *TestSuiteExecutor) Prepare(ctx *execution.TestPlanExecutionResult, def model.TestSuiteDef) error {
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test suite %s", def.Name)
	}

	testSuiteExecResult := &execution.TestSuiteExecutionResult{
		TestSuiteDef:              &def,
		NamedTestCasesExecResults: make(map[string]*execution.TestCaseExecutionResult),
		TestPlanExecutionResult:   ctx,
	}
	ctx.AddTestSuiteExecResults(testSuiteExecResult)
	t.testCaseExecutor = &TestCaseExecutor{}

	for _, testCaseDef := range def.Cases {
		if err := t.testCaseExecutor.Prepare(testSuiteExecResult, testCaseDef); err != nil {
			return err
		}
	}

	return nil
}

func (t *TestSuiteExecutor) Validate(result *execution.TestSuiteExecutionResult) error {

	for _, testCaseResult := range result.TestCasesExecResults {

		if len(testCaseResult.TestCaseDef.RequestRef) > 0 {
			if !result.HasNamed(testCaseResult.TestCaseDef.RequestRef) {
				return fmt.Errorf("depends %s test case %s not found", testCaseResult.TestCaseDef.RequestRef,
					testCaseResult.TestCaseDef.Name)
			} else {
				testCaseResult.TestCaseDef.Request = result.NamedTestCasesExecResults[testCaseResult.TestCaseDef.RequestRef].TestCaseDef.Request
			}
		}

		if err := t.testCaseExecutor.Validate(testCaseResult); err != nil {
			return err
		}
	}

	return nil
}
