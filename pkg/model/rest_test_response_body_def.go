package model

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"strings"
)

const JSON = "JSON"
const Plain = "Plain"
const File = "file"

type RestTestResponseBodyDef struct {
	Type          string
	BodyValidator RestTestResponseBodyValidator
}

func (d RestTestResponseBodyDef) Parse(bodyObj any) error {
	mapWrapper, err := collection.NewMapWrapperAny(bodyObj)
	if err != nil {
		return err
	}

	err = mapWrapper.Get("type", &d.Type)
	if err != nil {
		return err
	}

	if strings.EqualFold(JSON, d.Type) {
		d.BodyValidator = &RestTestResponseJSONBody{}
	} else if strings.HasSuffix(d.Type, File) {
		d.BodyValidator = &RestTestResponseFileBody{}
	} else if strings.EqualFold(Plain, d.Type) {
		d.BodyValidator = &RestTestResponsePlainBody{}
	} else {
		return fmt.Errorf("unsupported body type: %s", d.Type)
	}

	return d.BodyValidator.Parse(mapWrapper)
}
