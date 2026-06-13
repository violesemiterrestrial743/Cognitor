package workers

type Job[T any] struct {
	Value T
}

type Result[T any] struct {
	Value T
	Err   error
}
