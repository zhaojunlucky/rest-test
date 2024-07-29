package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"net/http"
)

type RestTestResponseBodyValidator interface {
	Parse(mapWrapper *collection.MapWrapper) error
	Validate(ctx *RestTestContext, resp *http.Response) error
}
