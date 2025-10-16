package controllers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "gorm.io/gorm"

    "github.com/example/hrms-backend/middlewares"
    "github.com/example/hrms-backend/models"
    "github.com/example/hrms-backend/services"
    "github.com/example/hrms-backend/utils"
)

type EmployeeController struct { db *gorm.DB; svc *services.EmployeeService }

func NewEmployeeController(db *gorm.DB) *EmployeeController { return &EmployeeController{db: db, svc: services.NewEmployeeService(db)} }

// @Summary List employees (HR)
// @Tags Employees
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Router /employees [get]
func (c *EmployeeController) List(w http.ResponseWriter, r *http.Request) {
    list, err := c.svc.List(); if err != nil { utils.Error(w, "error", http.StatusInternalServerError); return }
    utils.Success(w, "ok", list, http.StatusOK)
}

// @Summary Create employee (HR)
// @Tags Employees
// @Security BearerAuth
// @Param input body models.Employee true "Employee"
// @Success 201 {object} utils.APIResponse
// @Router /employees [post]
func (c *EmployeeController) Create(w http.ResponseWriter, r *http.Request) {
    var req models.Employee
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { utils.Error(w, "invalid body", http.StatusBadRequest); return }
    if err := c.svc.Create(&req); err != nil { utils.Error(w, "create error", http.StatusBadRequest); return }
    utils.Success(w, "created", req, http.StatusCreated)
}

// @Summary Update employee (HR)
// @Tags Employees
// @Security BearerAuth
// @Param id path int true "ID"
// @Param input body models.Employee true "Employee"
// @Success 200 {object} utils.APIResponse
// @Router /employees/{id} [put]
func (c *EmployeeController) Update(w http.ResponseWriter, r *http.Request) {
    id64, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
    var req models.Employee
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { utils.Error(w, "invalid body", http.StatusBadRequest); return }
    if err := c.svc.Update(uint(id64), &req); err != nil { utils.Error(w, "update error", http.StatusBadRequest); return }
    utils.Success(w, "updated", req, http.StatusOK)
}

// @Summary Delete employee (HR)
// @Tags Employees
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 204 {object} nil
// @Router /employees/{id} [delete]
func (c *EmployeeController) Delete(w http.ResponseWriter, r *http.Request) {
    id64, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
    if err := c.svc.Delete(uint(id64)); err != nil { utils.Error(w, "delete error", http.StatusBadRequest); return }
    w.WriteHeader(http.StatusNoContent)
}

// @Summary Get my employee profile (Employee)
// @Tags Employees
// @Security BearerAuth
// @Success 200 {object} utils.APIResponse
// @Router /employees/me [get]
func (c *EmployeeController) GetMe(w http.ResponseWriter, r *http.Request) {
    uid := r.Context().Value(middlewares.CtxUserID).(uint)
    var emp models.Employee
    if err := c.db.Where("user_id = ?", uid).First(&emp).Error; err != nil { utils.Error(w, "not found", http.StatusNotFound); return }
    utils.Success(w, "ok", emp, http.StatusOK)
}


