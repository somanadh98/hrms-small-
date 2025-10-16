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

type LeaveController struct {
	db  *gorm.DB
	svc *services.LeaveService
}

func NewLeaveController(db *gorm.DB) *LeaveController {
	return &LeaveController{db: db, svc: services.NewLeaveService(db)}
}

type leaveReq struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`
}

// Employee routes
func (c *LeaveController) Apply(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	var emp models.Employee
	if err := c.db.Where("user_id = ?", uid).First(&emp).Error; err != nil {
		utils.Error(w, "employee not found", http.StatusBadRequest)
		return
	}
	var req leaveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	s, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.Error(w, "invalid start", http.StatusBadRequest)
		return
	}
	e, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.Error(w, "invalid end", http.StatusBadRequest)
		return
	}
	if err := c.svc.Apply(emp.ID, s, e, req.Reason); err != nil {
		utils.Error(w, "apply error", http.StatusBadRequest)
		return
	}
	utils.Success(w, "applied", nil, http.StatusCreated)
}

func (c *LeaveController) DeleteMine(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	id64, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.Error(w, "invalid leave ID", http.StatusBadRequest)
		return
	}
	if err := c.svc.DeleteMine(uid, uint(id64)); err != nil {
		utils.Error(w, "delete error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *LeaveController) ListMine(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middlewares.CtxUserID).(uint)
	list, err := c.svc.ListMine(uid)
	if err != nil {
		utils.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.Success(w, "ok", list, http.StatusOK)
}

// HR routes
func (c *LeaveController) ListAll(w http.ResponseWriter, r *http.Request) {
	list, err := c.svc.ListAll()
	if err != nil {
		utils.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.Success(w, "ok", list, http.StatusOK)
}

func (c *LeaveController) Approve(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.Error(w, "invalid leave ID", http.StatusBadRequest)
		return
	}
	if err := c.svc.Approve(uint(id64)); err != nil {
		utils.Error(w, "approve error", http.StatusBadRequest)
		return
	}
	utils.Success(w, "approved", nil, http.StatusOK)
}

func (c *LeaveController) Reject(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.Error(w, "invalid leave ID", http.StatusBadRequest)
		return
	}
	if err := c.svc.Reject(uint(id64)); err != nil {
		utils.Error(w, "reject error", http.StatusBadRequest)
		return
	}
	utils.Success(w, "rejected", nil, http.StatusOK)
}
