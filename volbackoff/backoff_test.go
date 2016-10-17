package volbackoff_test

import (
	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/voldriver/volbackoff"
	"context"
	"github.com/docker/distribution/Godeps/_workspace/src/github.com/bugsnag/bugsnag-go/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Volbackoff", func() {

	var (
		err error

		now time.Time

		backoff volbackoff.ExponentialBackoff

		ctx       context.Context
		fakeClock *fakeclock.FakeClock

		op func(context.Context) error
	)

	JustBeforeEach(func() {
		backoff = volbackoff.NewExponentialBackOff(ctx, fakeClock)

		go func() {
			err = backoff.Retry(op)
		}()
	})

	BeforeEach(func() {
		now = time.Now()
		fakeClock = fakeclock.NewFakeClock(now)

		ctx = context.TODO()

		op = func(context.Context) error {
			return nil
		}
	})

	It("does something", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when operation fails and deadline is exceeded", func() {
		BeforeEach(func() {
			op = func(context.Context) error {
				return errors.Errorf("badness")
			}
		})

		It("Retry should return an error", func() {
			Eventually(func() bool {
				return err != nil
			}).Should(BeTrue())
		})
	})

	Context("when operation fails and deadline has not exceeded", func() {
		BeforeEach(func() {
			ctx, _ = context.WithDeadline(ctx, now.Add(30*time.Second))

			op = func(context.Context) error {
				return errors.Errorf("badness")
			}

		})

		It("Retry should return an error", func() {
			fakeClock.WaitForWatcherAndIncrement(time.Second * 31)

			Eventually(func() bool {
				return err != nil
			}).Should(BeTrue())
		})
	})

	//BeforeEach(func() {
	//	fmt.Println("BeforeEach ran")
	//	panic("wtf!")
	//})
	//
	//AfterEach(func() {
	//	fmt.Println("After Each ran")
	//})
	//
	//It("fails this example", func() {
	//	fmt.Println("ran 1")
	//})
	//
	//It("fails this example, too", func() {
	//	fmt.Println("ran 2")
	//})
	//
	//Context("nested group", func() {
	//	It("fails this third example", func() {
	//		fmt.Println("ran 1.1")
	//	})
	//
	//	It("fails this fourth example", func() {
	//		fmt.Println("ran 1.2")
	//	})
	//
	//	Describe("yet another level deep", func() {
	//		It("fails this last example", func() {
	//			fmt.Println("ran 1.3")
	//		})
	//	})
	//})
})
