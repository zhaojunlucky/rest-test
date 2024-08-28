package executor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/execution"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"time"
)

type TestSuiteExecutor struct {
	testCaseExecutor *TestCaseExecutor
}

func (t *TestSuiteExecutor) Execute(ctx *core.RestTestContext, environ env.Env, global *model.GlobalSetting, testSuiteExecResult *execution.TestSuiteExecutionResult) *report.TestSuiteReport {
	defer func() {
		testSuiteExecResult.Executed = true
		log.Infof("[Suite] end executetest suite: %s-%s", testSuiteExecResult.TestSuiteDef.GetID(), testSuiteExecResult.TestSuiteDef.Name)
	}()
	log.Infof("[Suite] start execute test suite: %s-%s", testSuiteExecResult.TestSuiteDef.GetID(), testSuiteExecResult.TestSuiteDef.Name)
	if testSuiteExecResult.TestSuiteReport == nil {
		testSuiteExecResult.TestSuiteReport = &report.TestSuiteReport{
			TestSuite: testSuiteExecResult.TestSuiteDef,
		}
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
	testSuiteCases := NewTestSuiteCaseContext()
	start := time.Now()

	for _, testCaseExecResult := range testSuiteExecResult.TestCasesExecResults {
		newEnv := env.NewReadWriteEnv(environ, testSuiteDef.Environment)
		testCaseReport := t.testCaseExecutor.Execute(ctx, newEnv, global, testCaseExecResult, testSuiteCases)

		suiteReport.ExecutionTime += testCaseReport.ExecutionTime
		suiteReport.AddTestCaseReport(testCaseReport)
	}
	suiteReport.TotalTime = time.Since(start).Seconds()
	suiteReport.Status = report.Completed
	return suiteReport
}

func (t *TestSuiteExecutor) Prepare(ctx *execution.TestPlanExecutionResult, def model.TestSuiteDef) error {
	defer func() {
		log.Infof("[Suite] end prepare test suite: %s-%s", def.GetID(), def.Name)
	}()
	log.Infof("[Suite] start prepare test suite: %s-%s", def.GetID(), def.Name)
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test suite %s", def.Name)
	}

	testSuiteExecResult := &execution.TestSuiteExecutionResult{
		TestSuiteDef:              &def,
		NamedTestCasesExecResults: make(map[string]*execution.TestCaseExecutionResult),
		TestPlanExecutionResult:   ctx,
	}
	ctx.AddTestSuiteExecResults(testSuiteExecResult)
	t.testCaseExecutor = NewTestCaseExecutor()

	for _, testCaseDef := range def.Cases {
		if err := t.testCaseExecutor.Prepare(testSuiteExecResult, testCaseDef); err != nil {
			return err
		}
	}

	return nil
}

func (t *TestSuiteExecutor) Validate(result *execution.TestSuiteExecutionResult) error {
	defer func() {
		log.Infof("[Suite] end validate test suite: %s-%s", result.TestSuiteDef.GetID(), result.TestSuiteDef.Name)
	}()
	log.Infof("[Suite] start validate test suite: %s-%s", result.TestSuiteDef.GetID(), result.TestSuiteDef.Name)
	for _, testCaseResult := range result.TestCasesExecResults {

		if len(testCaseResult.TestCaseDef.RequestRef) > 0 {
			if !result.HasNamed(testCaseResult.TestCaseDef.RequestRef) {
				return fmt.Errorf("depends %s test case %s not found", testCaseResult.TestCaseDef.RequestRef,
					testCaseResult.TestCaseDef.Name)
			} else {
				err := testCaseResult.TestCaseDef.CloneRequestRef(result.NamedTestCasesExecResults[testCaseResult.TestCaseDef.RequestRef].TestCaseDef.Request)
				if err != nil {
					log.Infof("failed to clone requestRef %s, error %v", testCaseResult.TestCaseDef.RequestRef, err)
					return err
				}
			}
		}

		if err := t.testCaseExecutor.Validate(testCaseResult); err != nil {
			return err
		}
	}

	return nil
}

func (t *TestSuiteExecutor) ExecuteSuite(ctx *core.RestTestContext, osEnv env.Env, testSuiteDef *model.TestSuiteDef) (*report.TestSuiteReport, error) {
	defer func() {
		log.Infof("[Suite] end run test suite: %s-%s", testSuiteDef.GetID(), testSuiteDef.Name)
	}()

	log.Infof("[Suite] start run test suite: %s-%s", testSuiteDef.GetID(), testSuiteDef.Name)
	testSuiteExecCtx, err := t.prepare(testSuiteDef)
	if testSuiteExecCtx == nil {
		return nil, fmt.Errorf("failed to prepare test plan: %v", err)
	}

	testSuiteExecCtx.TestSuiteReport = &report.TestSuiteReport{
		TestSuite: testSuiteDef,
	}

	suiteReport := testSuiteExecCtx.TestSuiteReport
	if err != nil {
		suiteReport.Error = err
		suiteReport.Status = report.ConfigError
		return suiteReport, err
	}
	log.Infof("validate test suite: %s-%s", testSuiteDef.GetID(), testSuiteDef.Name)
	if err = t.Validate(testSuiteExecCtx); err != nil {
		suiteReport.Error = err
		suiteReport.Status = report.ConfigError
		return suiteReport, err
	}
	planEnv := env.NewReadWriteEnv(osEnv, testSuiteDef.Environment)
	log.Infof("execute test suite: %s-%s", testSuiteDef.GetID(), testSuiteDef.Name)
	t.Execute(ctx, planEnv, &testSuiteDef.Global, testSuiteExecCtx)

	start := time.Now()

	suiteReport.ExecutionTime = time.Since(start).Seconds()
	suiteReport.TotalTime = time.Since(start).Seconds()
	return suiteReport, nil
}

func (t *TestSuiteExecutor) prepare(def *model.TestSuiteDef) (*execution.TestSuiteExecutionResult, error) {
	defer func() {
		log.Infof("[Suite] end prepare test suite: %s", def.Name)
	}()
	log.Infof("[Suite] prepare test suite: %s", def.Name)
	testSuiteExecResult := &execution.TestSuiteExecutionResult{
		TestSuiteDef:              def,
		NamedTestCasesExecResults: make(map[string]*execution.TestCaseExecutionResult),
		TestPlanExecutionResult:   nil,
	}
	t.testCaseExecutor = NewTestCaseExecutor()
	for _, testCaseDef := range def.Cases {
		if err := t.testCaseExecutor.Prepare(testSuiteExecResult, testCaseDef); err != nil {
			return nil, err
		}
	}

	return testSuiteExecResult, nil
}

func NewTestSuiteExecutor() *TestSuiteExecutor {
	return &TestSuiteExecutor{
		testCaseExecutor: NewTestCaseExecutor(),
	}
}
