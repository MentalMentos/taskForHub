package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:string" json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `gorm:"type:string" json:"-"`
	Role      string    `gorm:"type:string" json:"role"`
	IP        string    `gorm:"type:string" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
