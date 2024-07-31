package executor

import (
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
)

type TestSuiteExecutor []interface {
	Execute(ctx *core.RestTestContext, env *env.Env) (report.TestReport[model.TestSuiteDef], error)
}
