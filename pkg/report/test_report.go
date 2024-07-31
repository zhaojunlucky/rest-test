package report

type TestReport[T any] interface {
	GetExecutionTime() float64
	GetTotalTime() float64

	GetTestDef() *T
}
