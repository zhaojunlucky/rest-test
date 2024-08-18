package executor

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/zhaojunlucky/golib/pkg/env"
	"strings"
)

type JSScriptler struct {
	vm *goja.Runtime
}

func (js *JSScriptler) Expand(val string) (string, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return "", nil
	}
	if val[0] == '`' || val[len(val)-1] == '`' {
		return "", fmt.Errorf("cannot expand %s string with start or end backticks", val)
	}

	if strings.Contains(val, "\n") {
		return "", fmt.Errorf("cannot expand %s string with new line character", val)
	}

	o, err := js.vm.RunString(fmt.Sprintf("`%s`", val))
	if err != nil {
		return "", err
	}
	str, ok := o.Export().(string)
	if !ok {
		return "", fmt.Errorf("expect string but got %T", o.Export())
	}
	return str, nil
}

func (js *JSScriptler) ExpandMap(val map[string]string) (map[string]string, error) {
	expandMap := make(map[string]string)
	var err error
	for k, v := range val {
		expandMap[k], err = js.Expand(v)
		if err != nil {
			return nil, err
		}
	}
	return expandMap, nil
}

func (js *JSScriptler) Set(key string, val any) error {
	return js.vm.GlobalObject().Set(key, val)
}

func (js *JSScriptler) ExpandScript(script string) (string, error) {
	o, err := js.vm.RunString(script)
	if err != nil {
		return "", err
	}
	str, ok := o.Export().(string)
	if !ok {
		return "", fmt.Errorf("expect string but got %T", o.Export())
	}
	return str, nil
}

func (js *JSScriptler) ExpandScriptWithBody(script string, body string) (string, error) {
	err := js.vm.Set("body", body)
	if err != nil {
		return "", err
	}
	o, err := js.vm.RunString(script)
	if err != nil {
		return "", err
	}
	str, ok := o.Export().(string)
	if !ok {
		return "", fmt.Errorf("expect string but got %T", o.Export())
	}
	return str, nil
}

func (js *JSScriptler) RunScriptWithBody(script string, body any) (goja.Value, error) {
	err := js.vm.Set("body", body)
	if err != nil {
		return nil, err
	}
	o, err := js.vm.RunString(script)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func NewJSScriptler(env env.Env, testSuiteCases *TestSuiteCaseContext) (*JSScriptler, error) {

	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	registry := new(require.Registry)
	registry.Enable(vm)
	console.Enable(vm)

	err := vm.GlobalObject().Set("ctx", testSuiteCases.CaseResult)
	if err != nil {
		return nil, err
	}

	err = vm.GlobalObject().Set("env", env.GetAll())
	if err != nil {
		return nil, err
	}
	js := &JSScriptler{vm: vm}

	envs, err := js.ExpandMap(env.GetAll())
	if err != nil {
		return nil, err
	}

	err = vm.GlobalObject().Set("env", envs)
	if err != nil {
		return nil, err
	}

	return js, nil
}
