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

type TestPlanExecutor struct {
	testSuiteExecutor *TestSuiteExecutor
}

func (t *TestPlanExecutor) ExecutePlan(ctx *core.RestTestContext, environ env.Env, testPlanDef *model.TestPlanDef) (*report.TestPlanReport, error) {
	log.Infof("[Plan] start run test plan: %s", testPlanDef.Name)

	defer func() {
		log.Infof("[Plan] end run test plan: %s", testPlanDef.Name)
	}()

	testPlanExecCtx, err := t.Prepare(testPlanDef)
	if testPlanExecCtx == nil {
		return nil, fmt.Errorf("failed to prepare test plan: %v", err)
	}

	testPlanExecCtx.TestPlanReport = &report.TestPlanReport{
		TestPlan: testPlanDef,
	}

	planReport := testPlanExecCtx.TestPlanReport
	if err != nil {
		log.Errorf("failed to prepare test plan: %v", err)
		planReport.Error = err
		planReport.Status = report.ConfigError
		return planReport, err
	}

	if err = t.Validate(testPlanExecCtx); err != nil {
		log.Errorf("failed to validate test plan: %v", err)
		planReport.Error = err
		planReport.Status = report.ConfigError
		return planReport, err
	}
	planEnv := env.NewReadWriteEnv(environ, testPlanDef.Environment)

	t.Execute(ctx, planEnv, &testPlanDef.Global, testPlanExecCtx)

	start := time.Now()

	planReport.ExecutionTime = time.Since(start).Seconds()
	planReport.TotalTime = time.Since(start).Seconds()
	planReport.Status = report.Completed

	return planReport, nil
}

func (t *TestPlanExecutor) Execute(ctx *core.RestTestContext, environ env.Env, global *model.GlobalSetting, testPlanExecCtx *execution.TestPlanExecutionResult) {
	defer func() {
		testPlanExecCtx.Executed = true
		log.Infof("[Plan] end execute test plan: %s", testPlanExecCtx.TestPlanDef.Name)
	}()

	log.Infof("[Plan] start execute test plan: %s", testPlanExecCtx.TestPlanDef.Name)
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
	defer func() {
		log.Infof("[Plan] end prepare test plan: %s", def.Name)
	}()

	log.Infof("[Plan] prepare test plan: %s", def.Name)
	testPlanExecCtx = &execution.TestPlanExecutionResult{
		TestPlanDef:               def,
		NamedTestSuiteExecResults: make(map[string]*execution.TestSuiteExecutionResult),
	}
	for _, testSuiteDef := range def.Suites {
		err = t.testSuiteExecutor.Prepare(testPlanExecCtx, testSuiteDef)
		if err != nil {
			log.Errorf("failed to prepare test suite %s: %v", testSuiteDef.Name, err)
			return
		}
	}
	return
}

func (t *TestPlanExecutor) Validate(ctx *execution.TestPlanExecutionResult) error {
	defer func() {
		log.Infof("[Plan] end validate test plan: %s", ctx.TestPlanDef.Name)
	}()
	log.Infof("[Plan] validate test plan: %s", ctx.TestPlanDef.Name)
	for _, testSuiteExecResult := range ctx.TestSuiteExecResults {
		for _, dep := range testSuiteExecResult.TestSuiteDef.Depends {
			if !ctx.HasNamed(dep) {
				return fmt.Errorf("depends %s of test suite %s not found",
					dep, testSuiteExecResult.TestSuiteDef.Name)
			}
		}

		if err := t.testSuiteExecutor.Validate(testSuiteExecResult); err != nil {
			log.Errorf("failed to validate test suite %s: %v", testSuiteExecResult.TestSuiteDef.Name, err)
			return err
		}
	}
	return nil
}

func NewTestPlanExecutor() *TestPlanExecutor {
	return &TestPlanExecutor{
		testSuiteExecutor: &TestSuiteExecutor{},
	}
}
