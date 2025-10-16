package controllers

import (
    "encoding/json"
    "net/http"

    "gorm.io/gorm"

    "github.com/example/hrms-backend/models"
    "github.com/example/hrms-backend/services"
    "github.com/example/hrms-backend/utils"
)

type AuthController struct {
    db  *gorm.DB
    svc *services.AuthService
}

func NewAuthController(db *gorm.DB) *AuthController {
    return &AuthController{db: db, svc: services.NewAuthService(db)}
}

type registerReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

// @Summary Register new user
// @Tags Auth
// @Param input body registerReq true "Register"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Router /auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
    var req registerReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { utils.Error(w, "invalid body", http.StatusBadRequest); return }
    if utils.IsBlank(req.Username) || utils.IsBlank(req.Password) || utils.IsBlank(req.Role) {
        utils.Error(w, "missing fields", http.StatusBadRequest); return
    }
    hash, err := c.svc.HashPassword(req.Password)
    if err != nil { utils.Error(w, "hash error", http.StatusInternalServerError); return }
    user := models.User{Username: req.Username, PasswordHash: hash, Role: models.UserRole(req.Role)}
    if err := c.db.Create(&user).Error; err != nil { utils.Error(w, "create error", http.StatusBadRequest); return }
    utils.Success(w, "registered", map[string]any{"id": user.ID}, http.StatusCreated)
}

type loginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// @Summary Login
// @Tags Auth
// @Param input body loginReq true "Login"
// @Success 200 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    var req loginReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { utils.Error(w, "invalid body", http.StatusBadRequest); return }
    var user models.User
    if err := c.db.Where("username = ?", req.Username).First(&user).Error; err != nil { utils.Error(w, "invalid credentials", http.StatusUnauthorized); return }
    if !c.svc.CheckPassword(user.PasswordHash, req.Password) { utils.Error(w, "invalid credentials", http.StatusUnauthorized); return }
    at, err := c.svc.GenerateAccessToken(&user)
    if err != nil { utils.Error(w, "token error", http.StatusInternalServerError); return }
    rt, err := c.svc.GenerateRefreshToken(&user)
    if err != nil { utils.Error(w, "token error", http.StatusInternalServerError); return }
    utils.Success(w, "ok", map[string]string{"access_token": at, "refresh_token": rt}, http.StatusOK)
}

type refreshReq struct { RefreshToken string `json:"refresh_token"` }

// @Summary Refresh token
// @Tags Auth
// @Param input body refreshReq true "Refresh"
// @Success 200 {object} utils.APIResponse
// @Failure 401 {object} utils.APIResponse
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
    var req refreshReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { utils.Error(w, "invalid body", http.StatusBadRequest); return }
    uid, err := c.svc.ParseRefresh(req.RefreshToken)
    if err != nil { utils.Error(w, "invalid refresh", http.StatusUnauthorized); return }
    var user models.User
    if err := c.db.First(&user, uid).Error; err != nil { utils.Error(w, "user not found", http.StatusUnauthorized); return }
    at, err := c.svc.GenerateAccessToken(&user)
    if err != nil { utils.Error(w, "token error", http.StatusInternalServerError); return }
    utils.Success(w, "ok", map[string]string{"access_token": at}, http.StatusOK)
}


