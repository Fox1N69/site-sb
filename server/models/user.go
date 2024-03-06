package models

type Users struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password []byte `json:"-"`
}
