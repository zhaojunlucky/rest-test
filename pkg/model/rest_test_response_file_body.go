package model

import (
	"github.com/zhaojunlucky/golib/pkg/collection"
	"math"
	"net/http"
)

type RestTestResponseFileBody struct {
	RestTestRequest *RestTestRequestDef
	Max             int
	Length          int
	Min             int
	Sha256          string
}

func (d RestTestResponseFileBody) Validate(ctx *RestTestContext, resp *http.Response) error {
	return nil
}

func (d RestTestResponseFileBody) Parse(mapWrapper *collection.MapWrapper) error {
	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if mapWrapper.Has("sha256") {
		err := mapWrapper.Get("sha256", &d.Sha256)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("min") {
		err := mapWrapper.Get("min", &d.Min)
		if err != nil {
			return err
		}
	} else {
		d.Min = math.MinInt
	}

	if mapWrapper.Has("max") {
		err := mapWrapper.Get("max", &d.Max)
		if err != nil {
			return err
		}
	} else {
		d.Max = math.MinInt
	}
	return nil
}
