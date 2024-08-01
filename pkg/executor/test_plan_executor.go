package executor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"time"
)

type TestPlanExecutor struct {
}

func (t *TestPlanExecutor) Execute(ctx *core.RestTestContext, environ env.Env, testDef *model.TestPlanDef) (*report.TestPlanReport, error) {
	planReport := report.TestPlanReport{
		TestPlan: testDef,
	}
	if testDef.Enabled == false {
		planReport.Error = fmt.Errorf("test plan %s is disabled", testDef.Name)
		planReport.Status = report.Skipped
		return &planReport, nil
	}
	testSuiteExecutor := &TestSuiteExecutor{}
	start := time.Now()
	for _, testSuiteDef := range testDef.Suites {
		newEnv := env.NewReadWriteEnv(environ, testDef.Environment)
		testSuiteReport, err := testSuiteExecutor.Execute(ctx, newEnv, &testDef.Global, &testSuiteDef)
		if err != nil {
			log.Errorf("test suite %s failed, error: %v", testSuiteReport.TestSuite.Name, err)
		} else {
			log.Infof("test suite %s passed", testSuiteReport.TestSuite.Name)
		}
		planReport.ExecutionTime += testSuiteReport.ExecutionTime
		planReport.AddTestSuiteReport(testSuiteReport)
	}
	planReport.TotalTime = time.Since(start).Seconds()
	return &planReport, nil
}
