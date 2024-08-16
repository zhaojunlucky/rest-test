package model

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"net/http"
	"strings"
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

func (t *RestTestResponseDef) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	if t.Code != 0 && resp.StatusCode != t.Code {
		return nil, fmt.Errorf("invalid response code: %d, expect %d", resp.StatusCode, t.Code)
	}
	if len(t.ContentType) != 0 && !strings.HasPrefix(resp.Header.Get("Content-Type"), t.ContentType) {
		return nil, fmt.Errorf("invalid response content type: %s, expect %s", resp.Header.Get("Content-Type"), t.ContentType)
	}
	return t.Body.Validate(ctx, resp, js)
}
