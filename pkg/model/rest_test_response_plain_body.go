package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"math"
	"net/http"
	"regexp"
)

type RestTestResponsePlainBody struct {
	RestTestRequest *RestTestRequestDef
	Length          int
	Regex           *regexp.Regexp
}

func (d RestTestResponsePlainBody) Validate(ctx *RestTestContext, resp *http.Response) error {
	return nil
}

func (d RestTestResponsePlainBody) Parse(mapWrapper *collection.MapWrapper) error {
	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if mapWrapper.Has("regex") {
		var regStr string
		err := mapWrapper.Get("regex", &regStr)
		if err != nil {
			return err
		}

		d.Regex = regexp.MustCompile(regStr)
	}
	return nil
}
