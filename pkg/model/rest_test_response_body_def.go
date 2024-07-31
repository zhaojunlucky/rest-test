package model

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"net/http"
	"strings"
)

const JSON = "JSON"
const Plain = "Plain"
const File = "file"

type RestTestResponseBodyDef struct {
	RestTestRequest *RestTestRequestDef
	Type            string
	BodyValidator   RestTestResponseBodyValidator
}

func (d RestTestResponseBodyDef) Parse(bodyObj any) error {
	mapWrapper, err := collection.NewMapWrapperAny(bodyObj)
	if err != nil {
		return err
	}

	err = mapWrapper.Get("type", &d.Type)
	if err != nil {
		return err
	}

	if strings.EqualFold(JSON, d.Type) {
		d.BodyValidator = &RestTestResponseJSONBody{
			RestTestRequest: d.RestTestRequest,
		}
	} else if strings.EqualFold(d.Type, File) {
		d.BodyValidator = &RestTestResponseFileBody{
			RestTestRequest: d.RestTestRequest,
		}
	} else if strings.EqualFold(Plain, d.Type) {
		d.BodyValidator = &RestTestResponsePlainBody{
			RestTestRequest: d.RestTestRequest,
		}
	} else {
		return fmt.Errorf("unsupported body type: %s", d.Type)
	}

	return d.BodyValidator.Parse(mapWrapper)
}

func (d RestTestResponseBodyDef) Validate(ctx *core.RestTestContext, resp *http.Response) error {
	return d.BodyValidator.Validate(ctx, resp)
}
