package model

import (
	"crypto/sha256"
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"io"
	"math"
	"net/http"
)

type RestTestResponseFileBody struct {
	RestTestRequest *RestTestRequestDef
	Max             int64
	Length          int64
	Min             int64
	Sha256          string
}

func (d RestTestResponseFileBody) Validate(ctx *core.RestTestContext, resp *http.Response) (any, error) {
	if d.Length != math.MinInt && d.Length != resp.ContentLength {
		return "", fmt.Errorf("invalid content length: %d, expect %d", resp.ContentLength, d.Length)
	}

	if d.Min != math.MinInt && d.Min > resp.ContentLength {
		return "", fmt.Errorf("invalid content length: %d, expect >= %d", resp.ContentLength, d.Min)
	}

	if d.Max != math.MinInt && d.Max < resp.ContentLength {
		return "", fmt.Errorf("invalid content length: %d, expect <= %d", resp.ContentLength, d.Max)
	}

	if d.Sha256 != "" {
		realSha256, err := d.CalcSha256(resp.Body)

		if err != nil {
			return "", err
		}

		if realSha256 != d.Sha256 {
			return "", fmt.Errorf("invalid content sha256: %s, expect %s", realSha256, d.Sha256)
		}
	}
	return "", nil
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

func (d RestTestResponseFileBody) CalcSha256(body io.ReadCloser) (string, error) {
	h := sha256.New()

	var buf [4096]byte
	for {
		n, err := body.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		h.Write(buf[:n])
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}
