package model

type TeamVenue struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

type TeamResponse struct {
	Get        string        `json:"get"`
	Parameters Parameters    `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`
	Paging     Paging        `json:"paging"`
	TeamVenues []TeamVenue   `json:"response"`
}

type Parameters struct {
	League string `json:"league"`
	Season string `json:"season"`
}
