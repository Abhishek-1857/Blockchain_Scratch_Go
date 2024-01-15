package models

type Response struct {
	Block Block
	Vote  bool
	Id    string
}

var VoteMap = make(map[string][]Response)

type HttpResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
