package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"net/http"
)

type RestTestResponseDef struct {
	RestTestRequest *RestTestRequestDef
	Code            int
	ContentType     string
	Body            *RestTestResponseBodyDef
}

func (t *RestTestResponseDef) Parse(mapWrapper *collection.MapWrapper) error {
	respWrapper, err := mapWrapper.GetChild("response")
	if err != nil {
		return err
	}

	err = respWrapper.Get("code", &t.Code)
	if err != nil {
		return err
	}

	err = respWrapper.Get("contentType", &t.ContentType)
	if err != nil {
		return err
	}

	bodyObj, err := respWrapper.GetAny("body")
	if err != nil {
		return err
	}

	t.Body = &RestTestResponseBodyDef{
		RestTestRequest: t.RestTestRequest,
	}
	err = t.Body.Parse(bodyObj)
	if err != nil {
		return err
	}

	return nil

}

func (t *RestTestResponseDef) Validate(ctx *core.RestTestContext, resp *http.Response) (any, error) {
	return t.Body.Validate(ctx, resp)
}
