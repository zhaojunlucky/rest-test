package executor

import (
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"time"
)

type TestCaseExecutor struct {
}

func (t *TestCaseExecutor) Execute(ctx *core.RestTestContext, env env.Env, global *model.GlobalSetting, testDef *model.TestCaseDef) (*report.TestCaseReport, error) {

	testCaseReport := report.TestCaseReport{
		TestCase: testDef,
	}

	start := time.Now()
	testCaseReport.ExecutionTime = time.Since(start).Seconds()
	return &testCaseReport, nil
}
