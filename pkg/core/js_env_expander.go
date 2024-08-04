package core

type JSEnvExpander interface {
	Expand(string) (string, error)
	ExpandMap(map[string]string) (map[string]string, error)
	RunScript(string) (string, error)
	RunScriptWithBody(script string, body string) (string, error)
}
