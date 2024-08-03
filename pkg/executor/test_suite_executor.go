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

type TestSuiteExecutor struct {
	testCaseExecutor *TestCaseExecutor
}

func (t *TestSuiteExecutor) Execute(ctx *core.RestTestContext, environ env.Env, global *model.GlobalSetting, testDef *model.TestSuiteDef, testSuiteExecResult *execution.TestSuiteExecutionResult) (*report.TestSuiteReport, error) {

	suiteReport := report.TestSuiteReport{
		TestSuite: testDef,
	}

	if testDef.Enabled == false {
		suiteReport.Error = fmt.Errorf("test suite %s is disabled", testDef.Name)
		suiteReport.Status = report.Skipped
		return &suiteReport, nil
	}

	if len(testDef.Depends) > 0 {
		for _, dep := range testDef.Depends {
			if !testSuiteExecResult.HasNamed(dep) {
				suiteReport.Error = fmt.Errorf("depends test suite %s not found", dep)
				suiteReport.Status = report.ConfigError
				return &suiteReport, nil
			}

			if !Executed {
				suiteReport.Error = fmt.Errorf("depends test suite %s not executed", dep)
				suiteReport.Status = report.ConfigError
				return &suiteReport, nil
			}
		}
	}

	start := time.Now()
	for i, testCaseDef := range testDef.Cases {
		newEnv := env.NewReadWriteEnv(environ, testDef.Environment)
		testCaseReport, err := t.testCaseExecutor.Execute(ctx, newEnv, global, &testCaseDef, testSuiteExecResult.TestCasesExecResults[i])
		if err != nil {
			log.Errorf("test case %s failed, error: %v", testCaseDef.Name, err)
		} else {
			log.Infof("test case %s passed", testCaseDef.Name)
		}
		suiteReport.ExecutionTime += testCaseReport.ExecutionTime
		suiteReport.AddTestCaseReport(testCaseReport)
	}
	suiteReport.TotalTime = time.Since(start).Seconds()
	return &suiteReport, nil
}

func (t *TestSuiteExecutor) Prepare(ctx *execution.TestPlanExecutionResult, def model.TestSuiteDef) error {
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test suite %s", def.Name)
	}

	testSuiteExecResult := &execution.TestSuiteExecutionResult{
		TestSuiteDef:              &def,
		NamedTestCasesExecResults: make(map[string]*execution.TestCaseExecutionResult),
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
