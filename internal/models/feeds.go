package models

type Feeds struct {
	Id      int    `json:"id" db:"id" bson:"id, omitempty"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	PubDate string `json:"pub_date" db:"pub_date"`
	Link    string `json:"link" db:"link"`
}
