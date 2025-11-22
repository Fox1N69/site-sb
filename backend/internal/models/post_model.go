package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray is a custom type for storing array of strings in PostgreSQL
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray: invalid type")
	}

	return json.Unmarshal(bytes, s)
}

type Post struct {
	GormCustom
	PreviewImage string      `json:"preview_image" gorm:"type:text"` // First image (cover/thumbnail)
	Images       StringArray `json:"images" gorm:"type:jsonb"`       // Array of additional images
	Title        string      `json:"title" gorm:"type:varchar(255)"`
	Description  string      `json:"description" gorm:"type:text"`
	Text         string      `json:"text" gorm:"type:text"`
}
