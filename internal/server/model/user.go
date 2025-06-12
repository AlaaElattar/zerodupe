package model

// User represents a user in the system
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `json:"-"`
	Salt     string `json:"-"`
}
