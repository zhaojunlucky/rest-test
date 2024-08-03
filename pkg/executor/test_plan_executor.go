package executor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (t *TestPlanExecutor) Execute(ctx *core.RestTestContext, environ env.Env, testDef *model.TestPlanDef) (*report.TestPlanReport, error) {
	planReport := report.TestPlanReport{
		TestPlan: testDef,
	}

	testPlanExecResults, err := t.Prepare(testDef)
	testPlanExecResults.TestPlanReport = &planReport

	if err != nil {
		planReport.Error = err
		planReport.Status = report.ConfigError
		return &planReport, nil
	}

	if testDef.Enabled == false {
		planReport.Error = fmt.Errorf("test plan %s is disabled", testDef.Name)
		planReport.Status = report.Skipped
		return &planReport, nil
	}
	t.testSuiteExecutor = &TestSuiteExecutor{}
	start := time.Now()
	for i, testSuiteDef := range testDef.Suites {
		newEnv := env.NewReadWriteEnv(environ, testDef.Environment)
		testSuiteReport, err := t.testSuiteExecutor.Execute(ctx, newEnv, &testDef.Global, &testSuiteDef, testPlanExecResults.TestSuiteExecResults[i])
		if err != nil {
			log.Errorf("test suite %s failed, error: %v", testSuiteReport.TestSuite.Name, err)
		} else {
			log.Infof("test suite %s passed", testSuiteReport.TestSuite.Name)
		}
		planReport.ExecutionTime += testSuiteReport.ExecutionTime
		planReport.AddTestSuiteReport(testSuiteReport)
	}

	planReport.TotalTime = time.Since(start).Seconds()
	planReport.Status = report.Completed
	return &planReport, nil
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
