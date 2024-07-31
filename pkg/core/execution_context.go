package core

type ExecutionContext struct {
	execCtx map[string]any
}

func (e *ExecutionContext) Get(key string) any {
	return e.execCtx[key]
}

func (e *ExecutionContext) Set(key string, value any) {
	e.execCtx[key] = value
}

func (e *ExecutionContext) Has(key string) bool {
	_, ok := e.execCtx[key]
	return ok
}

func (e *ExecutionContext) GetExecCtx() map[string]any {
	return e.execCtx
}
