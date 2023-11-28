package meetup

type JobFunc func()

// JobQueue is a mindnummingly simple job queue, to support our minimal needs.
type JobQueue struct {
	jobChan chan *JobFunc
}

func NewJobQueue(n int) *JobQueue {
	return &JobQueue{
		jobChan: make(chan *JobFunc, n),
	}
}

func (jq *JobQueue) Run(fn JobFunc) {
	jq.jobChan <- &fn

	go func() {
		fn()
		<-jq.jobChan
	}()
}
