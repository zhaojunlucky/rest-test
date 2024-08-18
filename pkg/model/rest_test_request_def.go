package model

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
)

type RestTestRequestDef struct {
	URL        string
	Method     string
	Body       *RestTestRequestBodyDef
	Headers    map[string]string
	Parameters map[string]string
}

func (t *RestTestRequestDef) Parse(mapWrapper *collection.MapWrapper) error {
	reqWrapper, err := mapWrapper.GetChild("request")
	if err != nil {
		log.Errorf("parse request error: %s", err.Error())
		return err
	}

	err = reqWrapper.Get("url", &t.URL)
	if err != nil {
		log.Errorf("parse url error: %s", err.Error())
		return err
	}

	err = reqWrapper.Get("method", &t.Method)
	if err != nil {
		log.Errorf("parse method error: %s", err.Error())
		return err
	}

	if reqWrapper.Has("headers") {
		err = reqWrapper.Get("headers", &t.Headers)
		if err != nil {
			log.Errorf("parse headers error: %s", err.Error())
			return err
		}
	}

	if reqWrapper.Has("parameters") {
		err = reqWrapper.Get("parameters", &t.Parameters)
		if err != nil {
			log.Errorf("parse parameters error: %s", err.Error())
			return err
		}
	}

	t.Body = &RestTestRequestBodyDef{}
	if reqWrapper.Has("body") {
		var bodyObj any
		bodyObj, err = reqWrapper.GetAny("body")
		if err != nil {
			log.Errorf("parse body error: %s", err.Error())
			return err
		}
		return t.Body.Parse(bodyObj)

	}
	return nil
}
