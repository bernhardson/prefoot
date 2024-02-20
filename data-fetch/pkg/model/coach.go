package model

type Coach struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         int    `json:"age"`
	Birth       Birth  `json:"birth"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Photo       string `json:"photo"`
	Team        Team   `json:"team"`
	Career      Career `json:"career"`
}

type Career []struct {
	Team  Team   `json:"team"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type CoachResponse struct {
	Get        string        `json:"get"`
	Parameters interface{}   `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Paging     Paging        `json:"paging"`
	Response   []Coach       `json:"response"`
}
