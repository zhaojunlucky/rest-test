package executor

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("cannot expand %s string with start or end backticks", val)
		return "", fmt.Errorf("cannot expand %s string with start or end backticks", val)
	}

	if strings.Contains(val, "\n") {
		log.Errorf("cannot expand %s string with new line character", val)
		return "", fmt.Errorf("cannot expand %s string with new line character", val)
	}

	o, err := js.vm.RunString(fmt.Sprintf("`%s`", val))
	if err != nil {
		log.Errorf("expand %s error: %s", val, err.Error())
		return "", fmt.Errorf("expand %s error: %s", val, err.Error())
	}
	str, ok := o.Export().(string)
	if !ok {
		log.Errorf("expect string but got %T", o.Export())
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
			log.Errorf("expand %s error: %s", v, err.Error())
			return nil, fmt.Errorf("expand %s error: %s", v, err.Error())
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
		log.Errorf("expand %s error: %s", script, err.Error())
		return "", fmt.Errorf("expand %s error: %s", script, err.Error())
	}
	str, ok := o.Export().(string)
	if !ok {
		log.Errorf("expect string but got %T", o.Export())
		return "", fmt.Errorf("expect string but got %T", o.Export())
	}
	return str, nil
}

func (js *JSScriptler) ExpandScriptWithBody(script string, body string) (string, error) {
	err := js.vm.Set("body", body)
	if err != nil {
		log.Errorf("set body error: %s", err.Error())
		return "", fmt.Errorf("set body error: %s", err.Error())
	}
	o, err := js.vm.RunString(script)
	if err != nil {
		log.Errorf("expand %s error: %s", script, err.Error())
		return "", fmt.Errorf("expand %s error: %s", script, err.Error())
	}
	str, ok := o.Export().(string)
	if !ok {
		log.Errorf("expect string but got %T", o.Export())
		return "", fmt.Errorf("expect string but got %T", o.Export())
	}
	return str, nil
}

func (js *JSScriptler) RunScriptWithBody(script string, body any) (goja.Value, error) {
	err := js.vm.Set("body", body)
	if err != nil {
		log.Errorf("set body error: %s", err.Error())
		return nil, fmt.Errorf("set body error: %s", err.Error())
	}
	o, err := js.vm.RunString(script)
	if err != nil {
		log.Errorf("expand %s error: %s", script, err.Error())
		return nil, fmt.Errorf("expand %s error: %s", script, err.Error())
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
		log.Errorf("set ctx error: %s", err.Error())
		return nil, err
	}

	err = vm.GlobalObject().Set("env", env.GetAll())
	if err != nil {
		log.Errorf("set env error: %s", err.Error())
		return nil, err
	}
	js := &JSScriptler{vm: vm}

	envs, err := js.ExpandMap(env.GetAll())
	if err != nil {
		log.Errorf("expand env error: %s", err.Error())
		return nil, err
	}

	err = vm.GlobalObject().Set("env", envs)
	if err != nil {
		log.Errorf("set env error: %s", err.Error())
		return nil, err
	}

	return js, nil
}
