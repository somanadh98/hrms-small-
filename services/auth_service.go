package services

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "github.com/example/hrms-backend/models"
)

type AuthService struct {
    db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService { return &AuthService{db: db} }

func (s *AuthService) HashPassword(password string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(b), err
}

func (s *AuthService) CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "sub":  user.ID,
        "role": string(user.Role),
        "exp":  time.Now().Add(15 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
}

func (s *AuthService) GenerateRefreshToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func (s *AuthService) ParseRefresh(tokenStr string) (uint, error) {
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
    })
    if err != nil || !token.Valid {
        return 0, errors.New("invalid refresh token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid claims")
    }
    if sub, ok := claims["sub"].(float64); ok {
        return uint(sub), nil
    }
    return 0, errors.New("invalid subject")
}


