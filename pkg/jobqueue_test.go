package meetup

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobQueue", func() {
	var jq *JobQueue

	var n int
	var nMu sync.Mutex

	BeforeEach(func() {
		jq = NewJobQueue(2)

		n = 0
		nMu = sync.Mutex{}
	})

	It("can add jobs", func() {
		iterCount := 10
		wg := sync.WaitGroup{}
		wg.Add(iterCount)

		for i := 0; i < iterCount; i++ {
			jq.Run(func() {
				nMu.Lock()
				defer nMu.Unlock()

				n += 1
				wg.Done()
			})
		}

		wg.Wait()

		Expect(n).To(Equal(iterCount))
	})

	When("queue is full", Ordered, func() {
		var wg sync.WaitGroup

		BeforeAll(func() {
			wg = sync.WaitGroup{}
			wg.Add(2)

			for i := 0; i < 2; i++ {
				jq.Run(func() {
					nMu.Lock()
					n += 1
					nMu.Unlock()

					wg.Done()
					wg.Wait()
				})
			}
		})

		It("should block until space is freed", func() {
			blockWg := sync.WaitGroup{}
			blockWg.Add(1)

			go jq.Run(func() {
				time.Sleep(time.Millisecond * 100)
				wg.Wait()

				nMu.Lock()
				defer nMu.Unlock()

				n += 1

				blockWg.Done()
			})

			wg.Wait()
			Expect(n).To(Equal(2))

			blockWg.Wait()
			nMu.Lock()
			defer nMu.Unlock()
			Eventually(n).Should(Equal(3))
		})
	})
})
