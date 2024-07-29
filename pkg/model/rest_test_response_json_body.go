package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"math"
	"net/http"
	"reflect"
)

const OR = "or"
const AND = "and"

type RestTestResponseJSONBody struct {
	RestTestRequest     *RestTestRequestDef
	Array               bool
	Length              int
	ContainsRequestJSON bool
	Validators          map[string]any
}

func (d RestTestResponseJSONBody) Validate(ctx *RestTestContext, resp *http.Response) error {
	return nil
}
func (d RestTestResponseJSONBody) Parse(mapWrapper *collection.MapWrapper) error {

	if mapWrapper.Has("array") {
		err := mapWrapper.Get("array", &d.Array)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if mapWrapper.Has("containsRequestJSON") {
		err := mapWrapper.Get("containsRequestJSON", &d.ContainsRequestJSON)
		if err != nil {
			return err
		}
	}

	if !mapWrapper.Has("validators") {
		return nil
	}

	valType, err := mapWrapper.GetType("validators")
	if err != nil {
		return err
	}

	if valType.Type().Kind() == reflect.Map {
		err = mapWrapper.Get("validators", &d.Validators)
		if err != nil {
			return err
		}
	} else if valType.Type().Kind() == reflect.Slice {
		var validators []string
		err = mapWrapper.Get("validators", &validators)
		if err != nil {
			return err
		}

		d.Validators = map[string]any{
			AND: validators,
		}
	}

	return nil
}
