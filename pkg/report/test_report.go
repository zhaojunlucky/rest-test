package report

const (
	ConfigError = "ConfigError"
	Skipped     = "Skipped"
	Completed   = "Completed"
)

type TestReport[T any] interface {
	GetExecutionTime() float64
	GetTotalTime() float64

	GetTestDef() *T
	GetError() []error
	GetStatus() string
}
