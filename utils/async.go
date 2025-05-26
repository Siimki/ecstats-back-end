package utils

// RunAsync runs any function that returns (T, error) in a goroutine
// and returns a result channel and error channel.
func RunAsync[T any](fn func() (T, error)) (<-chan T, <-chan error) {
	resultCh := make(chan T, 1)
	errorCh := make(chan error, 1)

	go func() {
		result, err := fn()
		if err != nil {
			errorCh <- err
			return
		}
		resultCh <- result
	}()

	return resultCh, errorCh
}
