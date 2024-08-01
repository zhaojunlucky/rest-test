package model

import "github.com/zhaojunlucky/golib/pkg/collection"

type GlobalSetting struct {
	Headers map[string]string
	DataDir string
}

func (g *GlobalSetting) Parse(mapWrapper *collection.MapWrapper) error {

	var globalObj map[string]any
	err := mapWrapper.Get("global", &globalObj)
	if err != nil {
		return err
	}

	globalWrapper := collection.NewMapWrapper(globalObj)

	if globalWrapper.Has("dataDir") {
		err = globalWrapper.Get("dataDir", &g.DataDir)
		if err != nil {
			return err
		}
	}

	if globalWrapper.Has("headers") {
		err = mapWrapper.Get("headers", &g.Headers)
		if err != nil {
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
		Headers: g.mergeHeaders(global.Headers),
		DataDir: g.mergeDataDir(global.DataDir),
	}
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
