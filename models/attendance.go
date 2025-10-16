package models

import "time"

type AttendanceStatus string

const (
    StatusPresent AttendanceStatus = "PRESENT"
    StatusAbsent  AttendanceStatus = "ABSENT"
)

type Attendance struct {
    ID         uint             `gorm:"primaryKey" json:"id"`
    CreatedAt  time.Time        `json:"created_at"`
    UpdatedAt  time.Time        `json:"updated_at"`
    EmployeeID uint             `gorm:"index;not null" json:"employee_id"`
    Date       time.Time        `gorm:"type:date;index:idx_emp_date,unique" json:"date"`
    Status     AttendanceStatus `gorm:"type:varchar(16);not null" json:"status"`
    Version    uint             `gorm:"default:1" json:"version"`
}


