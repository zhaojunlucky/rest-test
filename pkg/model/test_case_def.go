package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
)

var counter int = 0

func incrementCounter() int {
	counter++
	return counter
}

type TestCaseDef struct {
	ID          int
	SuiteDef    *TestSuiteDef
	Name        string
	Description string
	Enabled     bool
	Environment map[string]string
	Request     *RestTestRequestDef
	RequestRef  string
	Response    *RestTestResponseDef
}

func (t *TestCaseDef) GetID() string {
	caseId := t.ID
	suiteId := t.SuiteDef.ID

	if t.SuiteDef.PlanDef != nil {
		return fmt.Sprintf("%d_%d_%d", suiteId, caseId, t.SuiteDef.PlanDef.ID)
	} else {
		return fmt.Sprintf("0_%d_%d", suiteId, caseId)
	}
}

func (t *TestCaseDef) CloneRequestRef(src *RestTestRequestDef) error {
	log.Infof("clone request ref for test case %s from test case %s", t.Description, src.CaseDef.Description)
	if t.Request != nil || len(t.RequestRef) <= 0 {
		return fmt.Errorf("cannot clone request ref for test case %s", t.Name)
	}
	t.Request = src.Clone(t)
	return t.Response.UpdateRequest(t.Request)
}

func (t *TestCaseDef) Parse(caseDef map[string]any) error {
	t.ID = incrementCounter()

	mapWrapper := collection.NewMapWrapper(caseDef)

	if mapWrapper.Has("name") {
		err := mapWrapper.Get("name", &t.Name)
		if err != nil {
			log.Errorf("key name parse error in test case %s: %s", t.Name, err.Error())
			return err
		}
	}

	err := mapWrapper.Get("desc", &t.Description)
	if err != nil {
		log.Errorf("key desc parsing error in test case %s: %s", t.Name, err.Error())
		return err
	}

	if !mapWrapper.Has("enabled") {
		t.Enabled = true
	} else {
		err = mapWrapper.Get("enabled", &t.Enabled)
		if err != nil {
			log.Errorf("key enabled parsing error in test case %s: %s", t.Name, err.Error())
			return err
		}
	}

	if !mapWrapper.Has("environment") {
		t.Environment = map[string]string{}
	} else {
		err = mapWrapper.Get("environment", &t.Environment)
		if err != nil {
			log.Errorf("key environment parsing error in test case %s: %s", t.Name, err.Error())
			return err
		}
	}

	if mapWrapper.Has("request") {
		t.Request = &RestTestRequestDef{
			CaseDef: t,
		}
		err = t.Request.Parse(mapWrapper)
		if err != nil {
			log.Errorf("key request parsing error in test case %s: %s", t.Name, err.Error())
			return err
		}

	} else if mapWrapper.Has("requestRef") {
		err = mapWrapper.Get("requestRef", &t.RequestRef)
		if err != nil {
			log.Errorf("key requestRef parsing error in test case %s: %s", t.Name, err.Error())
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
