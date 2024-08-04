package report

const (
	ConfigError     = "ConfigError"
	DependencyError = "DependencyError"
	Skipped         = "Skipped"
	Completed       = "Completed"
	InitError       = "InitError"
	ExecutionError  = "ExecutionError"
)

type TestReport[T any] interface {
	GetExecutionTime() float64
	GetTotalTime() float64

	GetTestDef() *T
	GetError() []error
	GetStatus() string
}
