package models

type User struct {
	ChatID int64 `gorm:"primaryKey;uniqueIndex"`
	State  string
	Data   string `gorm:"type:text"`
}
