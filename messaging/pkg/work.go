package messaging

type Instructions string

type Work struct {
	RequestId   int64        `json:"request_id"`
	File        string       `json:"file"`
	Filter      string       `json:"filter"`
	Instruction Instructions `json:"filter_args"`
	URL         string       `json:"url"`
}
