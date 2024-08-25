package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (t *RestTestResponseDef) UpdateRequest(req *RestTestRequestDef) error {
	t.RestTestRequest = req
	return t.Body.UpdateRequest(req)
}

func (t *RestTestResponseDef) Parse(mapWrapper *collection.MapWrapper) error {
	respWrapper, err := mapWrapper.GetChild("response")
	if err != nil {
		log.Errorf("parse response error: %s", err.Error())
		return err
	}

	fieldCnt := 0
	if respWrapper.Has("code") {
		fieldCnt++
		err = respWrapper.Get("code", &t.Code)
		if err != nil {
			log.Errorf("parse response code error: %s", err.Error())
			return err
		}
	}

	if respWrapper.Has("contentType") {
		fieldCnt++
		err = respWrapper.Get("contentType", &t.ContentType)
		if err != nil {
			log.Errorf("parse response contentType error: %s", err.Error())
			return err
		}
	}

	if respWrapper.Has("body") {
		fieldCnt++
		bodyObj, err := respWrapper.GetAny("body")
		if err != nil {
			log.Errorf("parse response body error: %s", err.Error())
			return err
		}

		t.Body = &RestTestResponseBodyDef{
			RestTestRequest: t.RestTestRequest,
		}
		err = t.Body.Parse(bodyObj)
		if err != nil {
			log.Errorf("parse response body error: %s", err.Error())
			return err
		}
	}

	if fieldCnt <= 0 {
		log.Warnf("invalid response definition, at least provide one field to validate response")
		return fmt.Errorf("invalid response definition, at least provide one field to validate response")
	}
	return nil

}

func (t *RestTestResponseDef) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	if t.Code != 0 && resp.StatusCode != t.Code {
		log.Errorf("invalid response code: %d, expect %d", resp.StatusCode, t.Code)
		return nil, fmt.Errorf("invalid response code: %d, expect %d", resp.StatusCode, t.Code)
	}
	if len(t.ContentType) != 0 && !strings.HasPrefix(resp.Header.Get("Content-Type"), t.ContentType) {
		log.Errorf("invalid response content type: %s, expect %s", resp.Header.Get("Content-Type"), t.ContentType)
		return nil, fmt.Errorf("invalid response content type: %s, expect %s", resp.Header.Get("Content-Type"), t.ContentType)
	}

	if t.Body != nil {
		return t.Body.Validate(ctx, resp, js)

	} else {
		return nil, nil
	}
}
