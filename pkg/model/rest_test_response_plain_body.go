package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (d *RestTestResponsePlainBody) UpdateRequest(req *RestTestRequestDef) error {
	d.RestTestRequest = req
	return nil
}

func (d *RestTestResponsePlainBody) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyStr := string(data)

	if d.Length != math.MinInt && d.Length != resp.ContentLength {
		log.Errorf("invalid content length: %d, expect %d", resp.ContentLength, d.Length)
		return bodyStr, fmt.Errorf("invalid content length: %d, expect %d", resp.ContentLength, d.Length)
	}

	if d.Regex != nil {

		if !d.Regex.MatchString(bodyStr) {
			log.Errorf("invalid content, expect match %s", d.Regex)
			return bodyStr, fmt.Errorf("invalid content, expect match %s", d.Regex)
		}
	}

	return bodyStr, nil
}

func (d *RestTestResponsePlainBody) Parse(mapWrapper *collection.MapWrapper) error {
	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			log.Errorf("parse length error: %s", err.Error())
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if mapWrapper.Has("regex") {
		var regStr string
		err := mapWrapper.Get("regex", &regStr)
		if err != nil {
			log.Errorf("parse regex error: %s", err.Error())
			return err
		}

		d.Regex, err = regexp.Compile(regStr)
		if err != nil {
			log.Errorf("parse regex error: %s", err.Error())
			return err
		}
	}
	return nil
}
