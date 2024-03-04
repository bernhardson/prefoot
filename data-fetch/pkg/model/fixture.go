package model

type Fixture struct {
	Fixture FixtureMeta `json:"fixture"`
	League  League      `json:"league"`
	Teams   Teams       `json:"teams"`
	Goals   GoalsFD     `json:"goals"`
	Score   Score       `json:"score"`
}

type FixtureMeta struct {
	ID        int     `json:"id"`
	Referee   string  `json:"referee"`
	Timezone  string  `json:"timezone"`
	Date      string  `json:"date"`
	Timestamp int64   `json:"timestamp"`
	Periods   Periods `json:"periods"`
	Venue     Venue   `json:"venue"`
	Status    Status  `json:"status"`
}
