package services

import (
    "errors"
    "sync"
    "time"

    "gorm.io/gorm"
    "github.com/example/hrms-backend/models"
)

type AttendanceService struct {
    db *gorm.DB
    mu sync.Mutex // protect critical sections like duplicate insert/update checks
}

func NewAttendanceService(db *gorm.DB) *AttendanceService { return &AttendanceService{db: db} }

// Add or update attendance for a date atomically using DB transaction
func (s *AttendanceService) Add(employeeID uint, date time.Time, status models.AttendanceStatus) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Attendance
        err := tx.Where("employee_id = ? AND date = ?", employeeID, date.Format("2006-01-02")).First(&existing).Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            rec := models.Attendance{EmployeeID: employeeID, Date: date, Status: status}
            return tx.Create(&rec).Error
        } else if err != nil {
            return err
        }
        // optimistic increment of Version to detect concurrent updates
        existing.Status = status
        return tx.Model(&models.Attendance{}).
            Where("id = ? AND version = ?", existing.ID, existing.Version).
            Updates(map[string]interface{}{"status": status, "version": existing.Version + 1}).Error
    })
}

func (s *AttendanceService) Delete(employeeID, id uint) error {
    return s.db.Where("id = ? AND employee_id = ?", id, employeeID).Delete(&models.Attendance{}).Error
}

func (s *AttendanceService) ListByEmployee(employeeID uint) ([]models.Attendance, error) {
    var list []models.Attendance
    if err := s.db.Where("employee_id = ?", employeeID).Order("date desc").Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func (s *AttendanceService) ListAll() ([]models.Attendance, error) {
    var list []models.Attendance
    if err := s.db.Order("date desc").Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func (s *AttendanceService) UpdateAny(id uint, status models.AttendanceStatus) error {
    return s.db.Model(&models.Attendance{ID: id}).Update("status", status).Error
}


