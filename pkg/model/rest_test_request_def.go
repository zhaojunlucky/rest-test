package model

import "github.com/zhaojunlucky/golib/pkg/collection"

type RestTestRequestDef struct {
	Url        string
	Method     string
	Body       RestTestRequestBodyDef
	Headers    map[string]string
	Parameters map[string]string
}

func (t *RestTestRequestDef) Parse(mapWrapper *collection.MapWrapper) error {
	reqWrapper, err := mapWrapper.GetChild("request")
	if err != nil {
		return err
	}

	err = reqWrapper.Get("url", &t.Url)
	if err != nil {
		return err
	}

	err = reqWrapper.Get("method", &t.Method)
	if err != nil {
		return err
	}

	if reqWrapper.Has("headers") {
		err = reqWrapper.Get("headers", &t.Headers)
		if err != nil {
			return err
		}
	}

	if reqWrapper.Has("parameters") {
		err = reqWrapper.Get("parameters", &t.Parameters)
		if err != nil {
			return err
		}
	}

	t.Body = RestTestRequestBodyDef{}
	if reqWrapper.Has("body") {
		var bodyObj any
		bodyObj, err = reqWrapper.GetAny("body")
		if err != nil {
			return err
		}
		return t.Body.Parse(bodyObj)

	}
	return nil
}
