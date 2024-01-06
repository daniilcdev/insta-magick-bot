package telegram

type Storage interface {
	NewRequest(file, requesterId string)
	GetRequester(file string) (string, error)
	RemoveRequest(file string)
}
