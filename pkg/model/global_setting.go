package model

import (
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
)

type GlobalSetting struct {
	Headers   map[string]string
	DataDir   string
	APIPrefix string
}

func (g *GlobalSetting) Parse(mapWrapper *collection.MapWrapper) error {

	var globalObj map[string]any
	err := mapWrapper.Get("global", &globalObj)
	if err != nil {
		log.Errorf("parse global error: %s", err.Error())
		return err
	}

	globalWrapper := collection.NewMapWrapper(globalObj)

	if globalWrapper.Has("dataDir") {
		err = globalWrapper.Get("dataDir", &g.DataDir)
		if err != nil {
			log.Errorf("parse dataDir error: %s", err.Error())
			return err
		}
	}

	if globalWrapper.Has("headers") {
		err = globalWrapper.Get("headers", &g.Headers)
		if err != nil {
			log.Errorf("parse headers error: %s", err.Error())
			return err
		}
	}

	if globalWrapper.Has("apiPrefix") {
		err = globalWrapper.Get("apiPrefix", &g.APIPrefix)
		if err != nil {
			log.Errorf("parse apiPrefix error: %s", err.Error())
			return err
		}
	}

	return nil
}

func (g *GlobalSetting) With(global *GlobalSetting) *GlobalSetting {
	if global == nil {
		return g
	}

	return &GlobalSetting{
		Headers:   g.mergeHeaders(global.Headers),
		DataDir:   g.mergeDataDir(global.DataDir),
		APIPrefix: g.mergeAPIPrefix(global.APIPrefix),
	}
}

func (g *GlobalSetting) Expand(js core.JSEnvExpander) (*GlobalSetting, error) {

	global := &GlobalSetting{}

	var err error
	global.DataDir, err = js.Expand(g.DataDir)
	if err != nil {
		log.Errorf("expand dataDir error: %s", err.Error())
		return nil, err
	}

	global.Headers, err = js.ExpandMap(g.Headers)
	if err != nil {
		log.Errorf("expand headers error: %s", err.Error())
		return nil, err
	}

	global.APIPrefix, err = js.Expand(g.APIPrefix)
	if err != nil {
		log.Errorf("expand apiPrefix error: %s", err.Error())
		return nil, err
	}
	return global, nil
}

func (g *GlobalSetting) mergeHeaders(headers map[string]string) map[string]string {
	copyMap := make(map[string]string)

	for k, v := range g.Headers {
		copyMap[k] = v
	}

	for k, v := range headers {
		copyMap[k] = v
	}
	return copyMap
}

func (g *GlobalSetting) mergeDataDir(dataDir string) string {
	if len(dataDir) > 0 {
		return dataDir
	}
	return g.DataDir
}

func (g *GlobalSetting) mergeAPIPrefix(prefix string) string {
	if len(prefix) > 0 {
		return prefix
	}
	return g.APIPrefix
}
