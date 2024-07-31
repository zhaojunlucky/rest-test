package executor

import (
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/report"
)

type TestExecutor[T any] interface {
	Execute(ctx *core.RestTestContext, env *env.Env) (report.TestReport[T], error)
}
