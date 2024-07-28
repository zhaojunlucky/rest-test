package model

type RestTestRequestBodyDef struct {
	File        string
	Environment map[string]string
	Body        string
	Script      string
}

func (d RestTestRequestBodyDef) Parse(bodyObj any) error {

}
