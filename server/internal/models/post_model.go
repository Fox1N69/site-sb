package models

type Post struct {
	GormCustom
	Title       string `json:"title"`
	Description string `json:"description"`
	Text        string `json:"text"`
}
