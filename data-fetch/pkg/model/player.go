package model

type PlayerDetails struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Birth     struct {
		Date    string `json:"date"`
		Place   string `json:"place"`
		Country string `json:"country"`
	} `json:"birth"`
	Nationality string `json:"nationality"`
	Height      string `json:"height"`
	Weight      string `json:"weight"`
	Injured     bool   `json:"injured"`
	Photo       string `json:"photo"`
}

type BirthInfo struct {
	Date    string `json:"date"`
	Place   string `json:"place"`
	Country string `json:"country"`
}

type Games struct {
	Appearances int    `json:"appearances"`
	Lineups     int    `json:"lineups"`
	Minutes     int    `json:"minutes"`
	Number      int    `json:"number"`
	Position    string `json:"position"`
	Rating      string `json:"rating"`
	Captain     bool   `json:"captain"`
}

type Substitutes struct {
	In    int `json:"in"`
	Out   int `json:"out"`
	Bench int `json:"bench"`
}

type Shots struct {
	Total int `json:"total"`
	On    int `json:"on"`
}

type Goals struct {
	Total    int `json:"total"`
	Conceded int `json:"conceded"`
	Assists  int `json:"assists"`
	Saves    int `json:"saves"`
}

type Passes struct {
	Total    int `json:"total"`
	Key      int `json:"key"`
	Accuracy int `json:"accuracy"`
}

type Tackles struct {
	Total         int `json:"total"`
	Blocks        int `json:"blocks"`
	Interceptions int `json:"interceptions"`
}

type Duels struct {
	Total int `json:"total"`
	Won   int `json:"won"`
}

type Dribbles struct {
	Attempts int `json:"attempts"`
	Success  int `json:"success"`
	Past     int `json:"past"`
}

type Fouls struct {
	Drawn     int `json:"drawn"`
	Committed int `json:"committed"`
}

type Cards struct {
	Yellow    int `json:"yellow"`
	YellowRed int `json:"yellowred"`
	Red       int `json:"red"`
}

type Penalty struct {
	Won       int `json:"won"`
	Committed int `json:"committed"`
	Scored    int `json:"scored"`
	Missed    int `json:"missed"`
	Saved     int `json:"saved"`
}

type PlayerStatistics struct {
	Team        Team        `json:"team"`
	League      League      `json:"league"`
	Games       Games       `json:"games"`
	Substitutes Substitutes `json:"substitutes"`
	Shots       Shots       `json:"shots"`
	Goals       Goals       `json:"goals"`
	Passes      Passes      `json:"passes"`
	Tackles     Tackles     `json:"tackles"`
	Duels       Duels       `json:"duels"`
	Dribbles    Dribbles    `json:"dribbles"`
	Fouls       Fouls       `json:"fouls"`
	Cards       Cards       `json:"cards"`
	Penalty     Penalty     `json:"penalty"`
}

type Player struct {
	PlayerDetails PlayerDetails      `json:"player"`
	Statistics    []PlayerStatistics `json:"statistics"`
}

type PlayerParameters struct {
	League string `json:"league"`
	Page   string `json:"page"`
	Season string `json:"season"`
}

type PlayerAPIResponse struct {
	Get        string           `json:"get"`
	Parameters PlayerParameters `json:"parameters"`
	Errors     []interface{}    `json:"errors"`
	Results    float64          `json:"results"`
	Paging     Paging           `json:"paging"`
	Response   []Player         `json:"response"`
}
