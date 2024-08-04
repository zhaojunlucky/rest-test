package model

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
)

type RestTestRequestBodyDef struct {
	File        string
	UploadFile  bool
	Environment map[string]string
	Body        string
	Script      string
	parsed      bool
}

func (d RestTestRequestBodyDef) Parse(bodyObj any) error {
	d.parsed = true
	if bodyObj == nil {
		return nil
	}
	bodyType := reflect.TypeOf(bodyObj)

	if bodyType.Kind() == reflect.String {
		d.Body = bodyObj.(string)
		return nil
	} else if bodyType.Kind() == reflect.Map {

		mapWrapper, err := collection.NewMapWrapperAny(bodyObj)
		if err != nil {
			return err
		}
		return d.parse(mapWrapper)
	} else {
		bodyStrType := reflect.ValueOf(d.Body)

		if bodyType.ConvertibleTo(bodyStrType.Type()) {
			d.Body = bodyStrType.Convert(bodyType).Interface().(string)
			return nil
		}

	}
	return fmt.Errorf("unsupported body type: %v", bodyType)
}

func (d RestTestRequestBodyDef) HasBody() bool {
	return d.parsed
}

func (d RestTestRequestBodyDef) parse(mapWrapper *collection.MapWrapper) error {
	if mapWrapper.Has("file") {
		err := mapWrapper.Get("file", &d.File)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("uploadFile") {
		err := mapWrapper.Get("uploadFile", &d.UploadFile)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("environment") {
		err := mapWrapper.Get("environment", &d.Environment)
		if err != nil {
			return err
		}
	}
	if mapWrapper.Has("body") {
		if len(d.File) > 0 {
			return fmt.Errorf("file and body cannot be set at the same time")
		}
		err := mapWrapper.Get("body", &d.Body)
		if err != nil {
			return err
		}
	}
	if mapWrapper.Has("script") {
		err := mapWrapper.Get("script", &d.Script)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d RestTestRequestBodyDef) GetBody(dataDir string, js core.JSEnvExpander) (io.Reader, error) {
	if !d.parsed {
		return nil, nil
	}
	var file io.Reader
	var err error
	var body = d.Body
	if len(d.File) > 0 {
		filePath := path.Join(dataDir, d.File)
		file, err = os.Open(filePath)
		if err != nil {
			return nil, err
		}
		if d.UploadFile {
			return file, nil
		} else {
			data, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			body = string(data)
		}
	}

	if len(d.Script) > 0 {
		return js.RunScript(d.Script)
	}

	return strings.NewReader(body), nil
}
