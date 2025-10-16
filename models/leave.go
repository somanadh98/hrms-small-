package models

import "time"

type LeaveStatus string

const (
    LeavePending  LeaveStatus = "PENDING"
    LeaveApproved LeaveStatus = "APPROVED"
    LeaveRejected LeaveStatus = "REJECTED"
)

type Leave struct {
    ID         uint        `gorm:"primaryKey" json:"id"`
    CreatedAt  time.Time   `json:"created_at"`
    UpdatedAt  time.Time   `json:"updated_at"`
    EmployeeID uint        `gorm:"index;not null" json:"employee_id"`
    StartDate  time.Time   `gorm:"type:date;not null" json:"start_date"`
    EndDate    time.Time   `gorm:"type:date;not null" json:"end_date"`
    Reason     string      `gorm:"size:255" json:"reason"`
    Status     LeaveStatus `gorm:"type:varchar(16);not null;default:PENDING" json:"status"`
    Version    uint        `gorm:"default:1" json:"version"`
}


