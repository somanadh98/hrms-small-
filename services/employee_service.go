package services

import (
    "gorm.io/gorm"
    "github.com/example/hrms-backend/models"
)

type EmployeeService struct { db *gorm.DB }

func NewEmployeeService(db *gorm.DB) *EmployeeService { return &EmployeeService{db: db} }

func (s *EmployeeService) Create(e *models.Employee) error { return s.db.Create(e).Error }
func (s *EmployeeService) Update(id uint, e *models.Employee) error {
    e.ID = id
    return s.db.Model(&models.Employee{ID: id}).Updates(e).Error
}
func (s *EmployeeService) Delete(id uint) error { return s.db.Delete(&models.Employee{}, id).Error }
func (s *EmployeeService) Get(id uint) (*models.Employee, error) {
    var m models.Employee
    if err := s.db.First(&m, id).Error; err != nil { return nil, err }
    return &m, nil
}
func (s *EmployeeService) List() ([]models.Employee, error) {
    var list []models.Employee
    if err := s.db.Find(&list).Error; err != nil { return nil, err }
    return list, nil
}


