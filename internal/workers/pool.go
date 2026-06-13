package workers

import (
	"context"
	"sync"
)

func Map[T any, R any](ctx context.Context, workers int, values []T, fn func(context.Context, T) (R, error)) []Result[R] {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan T)
	results := make(chan Result[R], len(values))
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for value := range jobs {
				if ctx.Err() != nil {
					var zero R
					results <- Result[R]{Value: zero, Err: ctx.Err()}
					continue
				}
				result, err := fn(ctx, value)
				results <- Result[R]{Value: result, Err: err}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, value := range values {
			select {
			case <-ctx.Done():
				return
			case jobs <- value:
			}
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	var collected []Result[R]
	for result := range results {
		collected = append(collected, result)
	}
	return collected
}
