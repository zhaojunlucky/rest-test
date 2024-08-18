package core

import "github.com/dop251/goja"

type JSEnvExpander interface {
	Expand(string) (string, error)
	ExpandMap(map[string]string) (map[string]string, error)
	ExpandScript(string) (string, error)
	ExpandScriptWithBody(script string, body string) (string, error)
	RunScriptWithBody(script string, body any) (goja.Value, error)
}
