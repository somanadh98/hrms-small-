package services

import (
    "errors"
    "sync"
    "time"

    "gorm.io/gorm"
    "github.com/example/hrms-backend/models"
)

type LeaveService struct {
    db *gorm.DB
    mu sync.Mutex // protect approval/rejection state transitions
}

func NewLeaveService(db *gorm.DB) *LeaveService { return &LeaveService{db: db} }

func (s *LeaveService) Apply(employeeID uint, start, end time.Time, reason string) error {
    lv := models.Leave{EmployeeID: employeeID, StartDate: start, EndDate: end, Reason: reason, Status: models.LeavePending}
    return s.db.Create(&lv).Error
}

func (s *LeaveService) DeleteMine(employeeID, id uint) error {
    return s.db.Where("id = ? AND employee_id = ? AND status = ?", id, employeeID, models.LeavePending).
        Delete(&models.Leave{}).Error
}

func (s *LeaveService) ListMine(employeeID uint) ([]models.Leave, error) {
    var list []models.Leave
    if err := s.db.Where("employee_id = ?", employeeID).Order("created_at desc").Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func (s *LeaveService) ListAll() ([]models.Leave, error) {
    var list []models.Leave
    if err := s.db.Order("created_at desc").Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func (s *LeaveService) setStatus(id uint, from models.LeaveStatus, to models.LeaveStatus) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.db.Transaction(func(tx *gorm.DB) error {
        var m models.Leave
        if err := tx.First(&m, id).Error; err != nil { return err }
        if m.Status != from {
            return errors.New("invalid status transition")
        }
        return tx.Model(&models.Leave{}).
            Where("id = ? AND version = ?", m.ID, m.Version).
            Updates(map[string]interface{}{"status": to, "version": m.Version + 1}).Error
    })
}

func (s *LeaveService) Approve(id uint) error { return s.setStatus(id, models.LeavePending, models.LeaveApproved) }
func (s *LeaveService) Reject(id uint) error  { return s.setStatus(id, models.LeavePending, models.LeaveRejected) }


