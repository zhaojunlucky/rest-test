package executor

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"time"
)

type TestSuiteExecutor struct {
}

func (t *TestSuiteExecutor) Execute(ctx *core.RestTestContext, env env.Env, global *model.GlobalSetting, testDef *model.TestSuiteDef) (*report.TestSuiteReport, error) {

	suiteReport := report.TestSuiteReport{
		TestSuite: testDef,
	}
	testCaseExecutor := &TestCaseExecutor{}
	start := time.Now()
	for _, testCaseDef := range testDef.Cases {
		testCaseReport, err := testCaseExecutor.Execute(ctx, env, &testCaseDef)
		if err != nil {
			log.Errorf("test case %s failed, error: %v", testCaseReport.TestCase.Name, err)

		} else {
			log.Infof("test case %s passed", testCaseReport.TestCase.Name)
		}
		suiteReport.ExecutionTime += testCaseReport.ExecutionTime
		suiteReport.AddTestCaseReport(testCaseReport)
	}
	suiteReport.TotalTime = time.Since(start).Seconds()
	return &suiteReport, nil
}
