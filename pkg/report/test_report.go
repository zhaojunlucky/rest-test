package report

const (
	Passed  = "Passed"
	Skipped = "Skipped"
	Failed  = "Failed"
)

type TestReport[T any] interface {
	GetExecutionTime() float64
	GetTotalTime() float64

	GetTestDef() *T
	GetError() []error
	GetStatus() string
}
