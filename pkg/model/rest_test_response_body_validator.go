package model

import "net/http"

type RestTestResponseBodyValidator interface {
	Validate(resp http.Response) error
}
