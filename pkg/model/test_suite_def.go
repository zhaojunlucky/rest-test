package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type TestSuiteDef struct {
	Name        string
	Depends     []string
	Enabled     bool
	Environment map[string]string
	Global      GlobalSetting
	Cases       []*TestCaseDef
	path        string
}

func (t *TestSuiteDef) Parse(file string) error {
	fi, err := os.Open(file)
	if err != nil {
		log.Errorf("open file %s error: %s", file, err.Error())
		return err
	}
	t.path = filepath.Dir(file)
	bytes, err := io.ReadAll(fi)
	if err != nil {
		log.Errorf("read file %s error: %s", file, err.Error())
		return err
	}
	def := make(map[string]any)
	err = yaml.Unmarshal(bytes, &def)
	if err != nil {
		log.Errorf("unmarshal file %s error: %s", file, err.Error())
		return err
	}

	mapWrapper := collection.NewMapWrapper(def)

	err = mapWrapper.Get("name", &t.Name)
	if err != nil {
		log.Errorf("key name not found in test suite %s", t.Name)
		return err
	}

	if mapWrapper.Has("depends") {
		depends, err := mapWrapper.GetAny("depends")
		if err != nil {
			log.Warningf("key depends not found in test plan %s", t.Name)
		} else {
			t.Depends, err = collection.GetObjAsSlice[string](depends)
			if err != nil {
				err = fmt.Errorf("key depends in test plan %s is not a string or a string list. %w", t.Name, err)
				log.Errorf(err.Error())
				return err
			}
		}
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

	t.Global = GlobalSetting{}
	err = t.Global.Parse(mapWrapper)
	if err != nil {
		return err
	}
	if !filepath.IsAbs(t.Global.DataDir) {
		t.Global.DataDir = filepath.Join(t.path, t.Global.DataDir)
	}
	log.Infof("test suite %s data dir: %s", t.Name, t.Global.DataDir)
	t.Cases, err = t.parseCases(mapWrapper)
	return err
}

func (t *TestSuiteDef) parseCases(mapWrapper *collection.MapWrapper) ([]*TestCaseDef, error) {
	var caseList []any
	err := mapWrapper.Get("cases", &caseList)
	if err != nil {
		return nil, err
	}
	if len(caseList) <= 0 {
		log.Errorf("test suite %s has no cases", t.Name)
		return nil, fmt.Errorf("test suite %s has no cases", t.Name)
	}

	var caseListDef []*TestCaseDef

	for i, caseName := range caseList {
		caseObj, ok := caseName.(map[string]any)

		if !ok {
			log.Errorf("the %d case in test suite %s is not a map", i, t.Name)
			return nil, fmt.Errorf("the %d case in test suite %s is not a map", i, t.Name)
		}
		caseDef := &TestCaseDef{}
		err = caseDef.Parse(caseObj)
		if err != nil {
			log.Errorf("parse %d case %s error: %s", i, caseName, err.Error())
			return nil, err
		}
		caseListDef = append(caseListDef, caseDef)
	}

	return caseListDef, nil
}
