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

type TestPlanExecutor struct {
	testSuiteExecutor *TestSuiteExecutor
}

func (t *TestPlanExecutor) ExecutePlan(environ env.Env, testPlanDef *model.TestPlanDef) (*report.TestPlanReport, error) {
	t.testSuiteExecutor = &TestSuiteExecutor{}

	testPlanExecCtx, err := t.Prepare(testPlanDef)
	if testPlanExecCtx == nil {
		return nil, fmt.Errorf("failed to prepare test plan: %v", err)
	}

	testPlanExecCtx.TestPlanReport = &report.TestPlanReport{
		TestPlan: testPlanDef,
	}

	planReport := testPlanExecCtx.TestPlanReport
	if err != nil {
		planReport.Error = err
		planReport.Status = report.ConfigError
		return planReport, err
	}

	if err = t.Validate(testPlanExecCtx); err != nil {
		planReport.Error = err
		planReport.Status = report.ConfigError
		return planReport, err
	}
	ctx := &core.RestTestContext{}
	planEnv := env.NewReadWriteEnv(environ, testPlanDef.Environment)

	t.Execute(ctx, planEnv, &testPlanDef.Global, testPlanExecCtx)

	start := time.Now()

	planReport.ExecutionTime = time.Since(start).Seconds()
	planReport.TotalTime = time.Since(start).Seconds()
	planReport.Status = report.Completed
	return planReport, nil
}

func (t *TestPlanExecutor) Execute(ctx *core.RestTestContext, environ env.Env, global *model.GlobalSetting, testPlanExecCtx *execution.TestPlanExecutionResult) {

	testPlanDef := testPlanExecCtx.TestPlanDef
	planReport := testPlanExecCtx.TestPlanReport
	if testPlanDef.Enabled == false {
		planReport.Error = fmt.Errorf("test plan %s is disabled", testPlanDef.Name)
		planReport.Status = report.Skipped
		return
	}

	start := time.Now()
	for _, testSuiteExecResult := range testPlanExecCtx.TestSuiteExecResults {
		newEnv := env.NewReadWriteEnv(environ, testSuiteExecResult.TestSuiteDef.Environment)
		suiteGlobal := global.With(&testSuiteExecResult.TestSuiteDef.Global)

		testSuiteReport := t.testSuiteExecutor.Execute(ctx, newEnv, suiteGlobal, testSuiteExecResult)

		planReport.ExecutionTime += testSuiteReport.ExecutionTime
		planReport.AddTestSuiteReport(testSuiteReport)
	}

	planReport.TotalTime = time.Since(start).Seconds()
	planReport.Status = report.Completed
}

func (t *TestPlanExecutor) Prepare(def *model.TestPlanDef) (testPlanExecCtx *execution.TestPlanExecutionResult, err error) {

	testPlanExecCtx = &execution.TestPlanExecutionResult{
		TestPlanDef:               def,
		NamedTestSuiteExecResults: make(map[string]*execution.TestSuiteExecutionResult),
	}
	for _, testSuiteDef := range def.Suites {
		err = t.testSuiteExecutor.Prepare(testPlanExecCtx, testSuiteDef)
		if err != nil {
			return
		}
	}
	return
}

func (t *TestPlanExecutor) Validate(ctx *execution.TestPlanExecutionResult) error {
	for _, testSuiteExecResult := range ctx.TestSuiteExecResults {
		for _, dep := range testSuiteExecResult.TestSuiteDef.Depends {
			if !ctx.HasNamed(dep) {
				return fmt.Errorf("depends %s of test suite %s not found",
					dep, testSuiteExecResult.TestSuiteDef.Name)
			}
		}

		if err := t.testSuiteExecutor.Validate(testSuiteExecResult); err != nil {
			return err
		}
	}
	return nil
}

func NewTestPlanExecutor() *TestPlanExecutor {
	return &TestPlanExecutor{}
}
