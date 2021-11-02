package concurrency

import (
	"sync"
)

// MergeErrorChannels joins multiple channels of error type into one. It observes all the input
// channels waiting for error value or closing. If all input channels are closed or returned a
// value, it closes its own result channel. Should be used with context.WithCancel in order to
// end awaiting goroutines in case of single error with use of context.CancelFunc.
// The application of this function is in the scope of waiting for the very first error or no errors.
func MergeErrorChannels(channels ...<-chan error) <-chan error {
	// TODO maybe without len
	resultChannel := make(chan error, len(channels))

	go func() {
		var wg sync.WaitGroup

		for _, ch := range channels {
			wg.Add(1)
			go func(ch <-chan error) {
				if err := <-ch; err != nil {
					resultChannel <- err
				}
				wg.Done()
			}(ch)
		}

		wg.Wait()
		close(resultChannel)
	}()

	return resultChannel
}
