package models

import "time"

type Employee struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    UserID     uint      `gorm:"uniqueIndex;not null" json:"user_id"`
    Name       string    `gorm:"size:120;not null" json:"name"`
    Position   string    `gorm:"size:120;not null" json:"position"`
    Department string    `gorm:"size:120;not null" json:"department"`
    Salary     float64   `gorm:"not null" json:"salary"`
    Version    uint      `gorm:"default:1" json:"version"`
}


