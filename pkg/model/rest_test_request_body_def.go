package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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
	bodyValid   bool
}

func (d *RestTestRequestBodyDef) Parse(bodyObj any) error {
	if bodyObj == nil {
		return nil
	}
	bodyType := reflect.TypeOf(bodyObj)

	if bodyType.Kind() == reflect.String {
		d.Body = bodyObj.(string)
		d.bodyValid = true
		return nil
	} else if bodyType.Kind() == reflect.Map {

		mapWrapper, err := collection.NewMapWrapperAny(bodyObj)
		if err != nil {
			log.Errorf("parse body error: %s", err.Error())
			return err
		}
		return d.parse(mapWrapper)
	}
	return fmt.Errorf("unsupported body type: %v", bodyType)
}

func (d *RestTestRequestBodyDef) HasBody() bool {
	return d.bodyValid
}

func (d *RestTestRequestBodyDef) parse(mapWrapper *collection.MapWrapper) error {
	if mapWrapper.Has("file") {
		err := mapWrapper.Get("file", &d.File)
		if err != nil {
			log.Errorf("parse file error %s", err.Error())
			return err
		}
	}

	if mapWrapper.Has("uploadFile") {
		err := mapWrapper.Get("uploadFile", &d.UploadFile)
		if err != nil {
			log.Errorf("parse uploadFile error %s", err.Error())
			return err
		}
	}

	if mapWrapper.Has("environment") {
		err := mapWrapper.Get("environment", &d.Environment)
		if err != nil {
			log.Errorf("parse environment error %s", err.Error())
			return err
		}
	}
	if mapWrapper.Has("body") {
		if len(d.File) > 0 {
			log.Errorf("file and body cannot be set at the same time")
			return fmt.Errorf("file and body cannot be set at the same time")
		}
		err := mapWrapper.Get("body", &d.Body)
		if err != nil {
			log.Errorf("parse body error %s", err.Error())
			return err
		}
	}
	if mapWrapper.Has("script") {
		err := mapWrapper.Get("script", &d.Script)
		if err != nil {
			log.Errorf("parse script error %s", err.Error())
			return err
		}
	}
	d.bodyValid = true
	return nil
}

func (d *RestTestRequestBodyDef) GetBody(dataDir string, js core.JSEnvExpander) (io.Reader, *string, error) {
	if !d.bodyValid {
		return nil, nil, nil
	}
	var file io.Reader
	var err error
	var body = d.Body
	if len(d.File) > 0 {
		filePath := path.Join(dataDir, d.File)
		file, err = os.Open(filePath)
		if err != nil {
			log.Errorf("failed to open file %s: %s", filePath, err.Error())
			return nil, nil, err
		}
		if d.UploadFile {
			fileBody := fmt.Sprintf("@%s", filePath)
			return file, &fileBody, nil
		} else {
			data, err := io.ReadAll(file)
			if err != nil {
				log.Errorf("failed to open file %s: %s", filePath, err.Error())

				return nil, nil, err
			}
			body = string(data)
		}
	}

	if len(d.Script) > 0 {
		body, err = js.ExpandScriptWithBody(d.Script, body)

		if err != nil {
			log.Errorf("failed to expand script %s: %s", d.Script, err.Error())
			return nil, nil, err
		}
	}

	return strings.NewReader(body), &body, nil
}
