package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/example/hrms-backend/middlewares"
	"github.com/example/hrms-backend/models"
	"github.com/example/hrms-backend/services"
	"github.com/example/hrms-backend/utils"
)

type AttendanceController struct {
	db  *gorm.DB
	svc *services.AttendanceService
}

func NewAttendanceController(db *gorm.DB) *AttendanceController {
	return &AttendanceController{db: db, svc: services.NewAttendanceService(db)}
}

type attendanceReq struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

// Employee
func (c *AttendanceController) AddMine(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	var emp models.Employee
	if err := c.db.Where("user_id = ?", uid).First(&emp).Error; err != nil {
		utils.Error(w, "employee not found", http.StatusBadRequest)
		return
	}
	var req attendanceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	d, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.Error(w, "invalid date", http.StatusBadRequest)
		return
	}
	if err := c.svc.Add(emp.ID, d, models.AttendanceStatus(req.Status)); err != nil {
		utils.Error(w, "save error", http.StatusBadRequest)
		return
	}
	utils.Success(w, "saved", nil, http.StatusCreated)
}

func (c *AttendanceController) DeleteMine(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	var emp models.Employee
	if err := c.db.Where("user_id = ?", uid).First(&emp).Error; err != nil {
		utils.Error(w, "employee not found", http.StatusBadRequest)
		return
	}
	id64, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.Error(w, "invalid attendance ID", http.StatusBadRequest)
		return
	}
	if err := c.svc.Delete(emp.ID, uint(id64)); err != nil {
		utils.Error(w, "delete error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *AttendanceController) ListMine(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	var emp models.Employee
	if err := c.db.Where("user_id = ?", uid).First(&emp).Error; err != nil {
		utils.Error(w, "employee not found", http.StatusBadRequest)
		return
	}
	list, err := c.svc.ListByEmployee(emp.ID)
	if err != nil {
		utils.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.Success(w, "ok", list, http.StatusOK)
}

// HR
func (c *AttendanceController) ListAll(w http.ResponseWriter, r *http.Request) {
	list, err := c.svc.ListAll()
	if err != nil {
		utils.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.Success(w, "ok", list, http.StatusOK)
}

func (c *AttendanceController) UpdateAny(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.Error(w, "invalid attendance ID", http.StatusBadRequest)
		return
	}
	var req attendanceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := c.svc.UpdateAny(uint(id64), models.AttendanceStatus(req.Status)); err != nil {
		utils.Error(w, "update error", http.StatusBadRequest)
		return
	}
	utils.Success(w, "updated", nil, http.StatusOK)
}
