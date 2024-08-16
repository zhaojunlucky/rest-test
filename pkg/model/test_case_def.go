package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
)

type TestCaseDef struct {
	Name        string
	Description string
	Enabled     bool
	Environment map[string]string
	Request     *RestTestRequestDef
	RequestRef  string
	Response    *RestTestResponseDef
}

func (t *TestCaseDef) Parse(caseDef map[string]any) error {
	mapWrapper := collection.NewMapWrapper(caseDef)

	if mapWrapper.Has("name") {
		err := mapWrapper.Get("name", &t.Name)
		if err != nil {
			return err
		}
	}

	err := mapWrapper.Get("desc", &t.Description)
	if err != nil {
		return err
	}

	if !mapWrapper.Has("enabled") {
		t.Enabled = true
	} else {
		err = mapWrapper.Get("enabled", &t.Enabled)
		if err != nil {
			return err
		}
	}

	if !mapWrapper.Has("environment") {
		t.Environment = map[string]string{}
	} else {
		err = mapWrapper.Get("environment", &t.Environment)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("request") {
		t.Request = &RestTestRequestDef{}

		err = t.Request.Parse(mapWrapper)
		if err != nil {
			return err
		}

	} else if mapWrapper.Has("requestRef") {
		err = mapWrapper.Get("requestRef", &t.RequestRef)
		if err != nil {
			return err
		}
	} else {
		log.Warnf("key request or requestRef not found in test case %s", t.Name)
		return fmt.Errorf("key request or requestRef not found in test case %s", t.Name)
	}

	t.Response = &RestTestResponseDef{
		RestTestRequest: t.Request,
	}
	return t.Response.Parse(mapWrapper)
}
