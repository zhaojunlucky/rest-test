package model

type RestTestResponseBodyDef struct {
	Type          string
	BodyValidator RestTestResponseBodyValidator
}

func (d RestTestResponseBodyDef) Parse(bodyObj any) error {

}
