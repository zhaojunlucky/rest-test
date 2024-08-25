package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"net/http"
)

type RestTestResponseBodyValidator interface {
	Parse(mapWrapper *collection.MapWrapper) error
	Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error)
	UpdateRequest(req *RestTestRequestDef) error
}
