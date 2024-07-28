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
