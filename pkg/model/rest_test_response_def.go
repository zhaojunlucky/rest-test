package model

import "github.com/zhaojunlucky/golib/pkg/collection"

type RestTestResponseDef struct {
	Code        int
	ContentType string
	Body        RestTestResponseBodyDef
}

func (t RestTestResponseDef) Parse(mapWrapper *collection.MapWrapper) error {
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

	t.Body = RestTestResponseBodyDef{}
	err = t.Body.Parse(bodyObj)
	if err != nil {
		return err
	}

	return nil

}
