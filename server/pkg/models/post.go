package models

type Posts struct {
	Id          uint    `json:"id" gorm:"primaryKey; autoIncrement:true"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PostBody    string `json: "postbody"`
}

