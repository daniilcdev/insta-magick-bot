package ports

type WorkDoneReporter interface {
	Done(work any)
}

type WorkFailedReporter interface {
	Failed(work any)
}

type WorkReporter interface {
	WorkDoneReporter
	WorkFailedReporter
}
