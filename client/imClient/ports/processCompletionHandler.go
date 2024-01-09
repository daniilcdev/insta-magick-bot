package ports

type CompletionHandler interface {
	OnProcessCompleted(dir string)
}
