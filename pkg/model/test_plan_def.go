package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"path/filepath"
)

var planCounter = 0

func increasePlanCounter() int {
	planCounter++
	return planCounter
}

type TestPlanDef struct {
	ID          int
	Name        string
	depends     []string // not supported now
	Enabled     bool
	Environment map[string]string
	Global      GlobalSetting
	Suites      []TestSuiteDef
	path        string // static
}

func (t *TestPlanDef) Parse(file string) error {
	t.ID = increasePlanCounter()
	fi, err := os.Open(file)
	if err != nil {
		log.Errorf("open file error: %s", err.Error())
		return err
	}
	t.path = filepath.Dir(file)
	bytes, err := io.ReadAll(fi)
	if err != nil {
		log.Errorf("read file error: %s", err.Error())
		return err
	}
	def := make(map[string]any)
	err = yaml.Unmarshal(bytes, &def)
	if err != nil {
		log.Errorf("unmarshal file error: %s", err.Error())
		return err
	}

	mapWrapper := collection.NewMapWrapper(def)
	err = mapWrapper.Get("name", &t.Name)
	if err != nil {
		log.Errorf("key name not found in test plan %s", t.Name)
		return err
	}

	if mapWrapper.Has("depends") {
		depends, err := mapWrapper.GetAny("depends")
		if err != nil {
			log.Warningf("key depends not found in test plan %s", t.Name)
		} else {
			t.depends, err = collection.GetObjAsSlice[string](depends)
			if err != nil {
				err = fmt.Errorf("key depends in test plan %s is not a string or a string list. %w", t.Name, err)
				log.Error(err)
				return err
			}
		}
	}

	if !mapWrapper.Has("enabled") {
		t.Enabled = true
	} else {
		err = mapWrapper.Get("enabled", &t.Enabled)
		if err != nil {
			log.Errorf("key enabled not found in test plan %s", t.Name)
			return err
		}
	}

	if !mapWrapper.Has("environment") {
		t.Environment = map[string]string{}
	} else {
		err = mapWrapper.Get("environment", &t.Environment)
		if err != nil {
			log.Errorf("key environment not found in test plan %s", t.Name)
			return err
		}
	}

	t.Global = GlobalSetting{}
	err = t.Global.Parse(mapWrapper)
	if err != nil {
		log.Errorf("parse global error: %s", err.Error())
		return err
	}
	if !filepath.IsAbs(t.Global.DataDir) {
		t.Global.DataDir = filepath.Join(t.path, t.Global.DataDir)
	}
	log.Infof("test plan %s data dir is %s", t.Name, t.Global.DataDir)

	t.Suites, err = t.parseSuites(mapWrapper)
	return err
}

func (t *TestPlanDef) parseSuites(mapWrapper *collection.MapWrapper) ([]TestSuiteDef, error) {
	var suiteNames []string
	err := mapWrapper.Get("suites", &suiteNames)
	if err != nil {
		log.Errorf("key suites not found in test plan %s", t.Name)
		return nil, err
	}

	if len(suiteNames) <= 0 {
		log.Errorf("test plan %s has no suite", t.Name)
		return nil, fmt.Errorf("test plan %s has no suite", t.Name)
	}

	baseDir := t.path
	if len(t.Global.DataDir) > 0 {
		baseDir = t.Global.DataDir
	}
	var suites []TestSuiteDef
	for i, name := range suiteNames {
		suiteDef := TestSuiteDef{
			PlanDef: t,
		}
		err = suiteDef.Parse(path.Join(baseDir, fmt.Sprintf("%s.yml", name)))
		if err != nil {
			log.Errorf("parse %d suite %s error: %s", i, name, err.Error())
			return nil, err
		}
		suites = append(suites, suiteDef)
	}

	return suites, nil
}
