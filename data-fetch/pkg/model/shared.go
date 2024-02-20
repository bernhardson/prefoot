package model


type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type Birth struct {
	Date    string `json:"date"`
	Place   string `json:"place"`
	Country string `json:"country"`
}
