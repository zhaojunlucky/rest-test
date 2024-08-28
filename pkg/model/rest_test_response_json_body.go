package model

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/core/validator"
	"io"
	"math"
	"net/http"
	"strings"
	"unicode"
)

type RestTestResponseJSONBody struct {
	RestTestRequest *RestTestRequestDef
	Array           bool
	Length          int
	JSONPathOPNode  *validator.JSONPathOperatorNode
	Script          string
}

func (d *RestTestResponseJSONBody) UpdateRequest(req *RestTestRequestDef) error {
	d.RestTestRequest = req
	return nil
}

func (d *RestTestResponseJSONBody) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(data)

	if err = writeBody(ctx, d.RestTestRequest.CaseDef, resp, bodyStr); err != nil {
		log.Warnf("write test case response body error: %s", err.Error())
	}
	for _, r := range bodyStr {
		if unicode.IsSpace(r) {
			continue
		}

		if r == '{' {
			if d.Array {
				return nil, fmt.Errorf("invalid json body, must be array, but got object")
			}
			break
		} else if r == '[' {
			if !d.Array {
				return nil, fmt.Errorf("invalid json body, must be object, but got array")
			}
		} else {
			return nil, fmt.Errorf("invalid json body, must start with [ or {, but got %s", string(r))
		}
	}
	jsonDecoder := json.NewDecoder(strings.NewReader(bodyStr))
	jsonDecoder.UseNumber()

	if d.Array {
		var arr []any
		err = jsonDecoder.Decode(&arr)
		if err != nil {
			return nil, err
		}

		if d.Length != math.MinInt && d.Length != len(arr) {
			return nil, fmt.Errorf("invalid JSON Array length: %d, expect %d", len(arr), d.Length)
		}
		return d.validate(core.ConvertArr(arr), js)
	} else {
		var obj map[string]any
		err = jsonDecoder.Decode(&obj)
		if err != nil {
			return nil, err
		}
		return d.validate(core.ConvertObj(obj), js)
	}
}

func (d *RestTestResponseJSONBody) Parse(mapWrapper *collection.MapWrapper) error {

	if mapWrapper.Has("array") {
		err := mapWrapper.Get("array", &d.Array)
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

	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if !mapWrapper.Has("validators") {
		return nil
	}

	validators, err := mapWrapper.GetAny("validators")
	if err != nil {
		return err
	}

	d.JSONPathOPNode, err = validator.NewJSONPathRootValidator(validators)

	if err != nil {
		return err
	}

	return nil
}

func (d *RestTestResponseJSONBody) validate(obj any, js core.JSEnvExpander) (any, error) {
	if len(d.Script) > 0 {
		log.Info("validate body with script")
		_, err := js.RunScriptWithBody(d.Script, obj)
		if err != nil {
			log.Errorf("script error: %s", err.Error())
			return nil, err
		}
	}
	if d.JSONPathOPNode != nil {
		log.Info("validate body with jsonpath")
		err := d.JSONPathOPNode.Validate(js, obj)
		if err != nil {
			log.Errorf("jsonpath validate error: %s", err.Error())
			return nil, err
		}
	}
	return obj, nil
}
