package scheduler

type Scheduler interface {
	SelectCAndidateNodes()
	Score()
	Pick()
}
