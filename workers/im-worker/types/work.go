package types

type Instructions string

type Work struct {
	File   string
	Filter Instructions
}

type Worker interface {
	OnWorkReceived(work Work)
}
