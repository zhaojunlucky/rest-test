package model

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"io"
	"math"
	"net/http"
	"regexp"
)

type RestTestResponsePlainBody struct {
	RestTestRequest *RestTestRequestDef
	Length          int64
	Regex           *regexp.Regexp
}

func (d RestTestResponsePlainBody) Validate(ctx *core.RestTestContext, resp *http.Response) error {
	if d.Length != math.MinInt && d.Length != resp.ContentLength {
		return fmt.Errorf("invalid content length: %d, expect %d", resp.ContentLength, d.Length)
	}

	if d.Regex != nil {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if !d.Regex.MatchString(string(data)) {
			return fmt.Errorf("invalid content, expect match %s", d.Regex)
		}
	}

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
