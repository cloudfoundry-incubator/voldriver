package volbackoff

import (
	"time"

	"code.cloudfoundry.org/clock"
	"context"
	"fmt"
)

const (
	backoffInitialInterval = 500 * time.Millisecond
	backoffIncrement       = 1.5
)

type ExponentialBackoff interface {
	Retry(operation func(context.Context) error) error
}

type exponentialBackOff struct {
	deadline time.Time
	ctx      context.Context
	clock    clock.Clock
}

// newExponentialBackOff takes a maximum elapsed time, after which the
// exponentialBackOff stops retrying the operation.
func NewExponentialBackOff(ctx context.Context, clock clock.Clock) ExponentialBackoff {
	deadline, _ := ctx.Deadline()
	fmt.Printf("deadline is %#v\n", deadline)
	//if ok == false {
	//	panic("illegal context, must set a deadline")
	//}

	return &exponentialBackOff{
		deadline: deadline,
		ctx:      ctx,
		clock:    clock,
	}
}

// Retry takes a retriable operation, and calls it until either the operation
// succeeds, or the retry timeout occurs.
func (b *exponentialBackOff) Retry(operation func(ctx context.Context) error) error {
	var (
		startTime       time.Time = b.clock.Now()
		backoffInterval time.Duration
		backoffExpired  bool
	)

	for {
		err := operation(b.ctx)
		if err == nil {
			return nil
		}

		if err == context.Canceled || err == context.DeadlineExceeded {
			return err
		}

		backoffInterval, backoffExpired = b.incrementInterval(startTime, backoffInterval)
		if backoffExpired {
			return err
		}

		fmt.Printf("here")
		b.clock.Sleep(backoffInterval)
	}
}

func (b *exponentialBackOff) incrementInterval(startTime time.Time, currentInterval time.Duration) (nextInterval time.Duration, expired bool) {
	if b.clock.Now().After(b.deadline) {
		return 0, true
	}

	switch {
	case currentInterval == 0:
		nextInterval = backoffInitialInterval
	//case elapsedTime+backoff(currentInterval) > b.deadline:
	//	nextInterval = time.Millisecond + b.maxElapsedTime - elapsedTime
	default:
		nextInterval = backoff(currentInterval)
	}

	return nextInterval, false
}

func backoff(interval time.Duration) time.Duration {
	return time.Duration(float64(interval) * backoffIncrement)
}
