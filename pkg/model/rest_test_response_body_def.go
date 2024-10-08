package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (d *RestTestResponseBodyDef) UpdateRequest(req *RestTestRequestDef) error {
	log.Infof("update request for body")
	d.RestTestRequest = req
	return d.BodyValidator.UpdateRequest(req)
}

func (d *RestTestResponseBodyDef) Parse(bodyObj any) error {
	mapWrapper, err := collection.NewMapWrapperAny(bodyObj)
	if err != nil {
		log.Errorf("parse body error: %s", err.Error())
		return err
	}

	if !mapWrapper.Has("type") {
		log.Debugf("body type not found, use default: %s", JSON)
		d.Type = JSON
	} else {
		err = mapWrapper.Get("type", &d.Type)
		if err != nil {
			log.Errorf("parse body type error: %s", err.Error())
			return err
		}
	}
	log.Debugf("body type: %s", d.Type)
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

func (d *RestTestResponseBodyDef) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	return d.BodyValidator.Validate(ctx, resp, js)
}
