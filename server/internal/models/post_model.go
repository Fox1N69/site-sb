package models

type Post struct {
	GormCustom
	PreviewImage string `json:"preview_image"`
	FullImage    string `json:"full_image"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Text         string `json:"text"`
}
