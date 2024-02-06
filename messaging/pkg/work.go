package messaging

type Instructions string

type Work struct {
	File        string       `json:"file"`
	Filter      Instructions `json:"filter"`
	RequesterId string       `json:"requester_id"`
	URL         string       `json:"url"`
}
