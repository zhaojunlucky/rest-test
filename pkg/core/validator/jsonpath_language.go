package validator

import (
	"context"
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

var lang = newJSONPathLanguage()

type JSONPathLanguage struct {
	language   *gval.Language
	expression map[string]gval.Evaluable
}

func (l *JSONPathLanguage) Compile(expr string) error {
	if _, ok := l.expression[expr]; ok {
		return nil
	}
	eval, err := l.language.NewEvaluable(expr)
	if err != nil {
		return err
	}

	l.expression[expr] = eval
	return nil
}

func (l *JSONPathLanguage) Evaluate(expr string, v any) (any, error) {
	eval, ok := l.expression[expr]
	if !ok {
		return nil, fmt.Errorf("uncompiled expression: %s", expr)
	}
	return eval(context.Background(), v)
}

func newJSONPathLanguage() *JSONPathLanguage {
	extLangs := GetExtensionLanguage()
	allLangs := append(extLangs, jsonpath.Language())
	lan := gval.Full(allLangs...)

	return &JSONPathLanguage{
		language:   &lan,
		expression: make(map[string]gval.Evaluable),
	}
}

func NewJSONPathLanguage() *JSONPathLanguage {
	return lang
}
