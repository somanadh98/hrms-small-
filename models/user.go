package models

import "time"

type UserRole string

const (
    RoleHR       UserRole = "HR"
    RoleEmployee UserRole = "EMPLOYEE"
)

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Username     string    `gorm:"uniqueIndex;size:80;not null" json:"username"`
    PasswordHash string    `gorm:"not null" json:"-"`
    Role         UserRole  `gorm:"type:varchar(16);not null" json:"role"`
    // Optimistic locking version
    Version uint `gorm:"default:1" json:"version"`
}


