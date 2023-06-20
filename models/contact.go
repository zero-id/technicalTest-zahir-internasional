package models

import "time"

type Contact struct {
	ID        string    `json:"id" gorm:"type:char(36);primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	Gender    string    `json:"gender" gorm:"type:varchar(255)"`
	Phone     string    `json:"phone" gorm:"type:varchar(255)"`
	Email     string    `json:"email" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
