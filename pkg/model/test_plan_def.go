package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

type TestPlanDef struct {
	Name        string
	depends     []string // not supported now
	Enabled     bool
	Environment map[string]string
	Global      GlobalSetting
	Suites      []TestSuiteDef
	path        string
}

func (t *TestPlanDef) Parse(file string) error {
	fi, err := os.Open(file)
	if err != nil {
		return err
	}
	t.path = path.Base(file)
	bytes, err := io.ReadAll(fi)
	if err != nil {
		return err
	}
	def := make(map[string]any)
	err = yaml.Unmarshal(bytes, &def)
	if err != nil {
		return err
	}

	mapWrapper := collection.NewMapWrapper(def)
	err = mapWrapper.Get("name", &t.Name)
	if err != nil {
		return err
	}

	depends, err := mapWrapper.GetAny("depends")
	if err != nil {
		log.Warningf("key depends not found in test plan %s", t.Name)
	} else {
		t.depends, err = collection.GetObjAsSlice[string](depends)
		if err != nil {
			err = fmt.Errorf("key depends in test plan %s is not a string or a string list. %w", t.Name, err)
			return err
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

	t.Suites, err = t.parseSuites(mapWrapper)

	return nil
}

func (t *TestPlanDef) parseSuites(mapWrapper *collection.MapWrapper) ([]TestSuiteDef, error) {
	var suiteNames []string
	err := mapWrapper.Get("suites", &suiteNames)
	if err != nil {
		return nil, err
	}

	if len(suiteNames) <= 0 {
		return nil, fmt.Errorf("test plan %s has no suite", t.Name)
	}

	baseDir := t.path
	if len(t.Global.DataDir) > 0 {
		baseDir = t.Global.DataDir
	}
	var suites []TestSuiteDef
	for _, name := range suiteNames {
		suiteDef := TestSuiteDef{}
		err = suiteDef.Parse(path.Join(baseDir, fmt.Sprintf("%s.yml", name)))
		if err != nil {
			return nil, err
		}
		suites = append(suites, suiteDef)
	}

	return suites, nil
}
