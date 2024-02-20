package model

type PlayerSquad struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Number   int    `json:"number"`
	Position string `json:"position"`
	Photo    string `json:"photo"`
}

type Squad struct {
	Team    Team          `json:"team"`
	Players []PlayerSquad `json:"players"`
}

type SquadResponse struct {
	Get        string        `json:"get"`
	Parameters interface{}   `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int           `json:"results"`

	Response []Squad `json:"response"`
}
