package middlewares

import (
    "context"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
    CtxUserID ctxKey = "user_id"
    CtxUserRole ctxKey = "user_role"
)

func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            http.Error(w, "missing bearer token", http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(auth, "Bearer ")
        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "invalid claims", http.StatusUnauthorized)
            return
        }
        uid, _ := claims["sub"].(float64)
        role, _ := claims["role"].(string)
        ctx := context.WithValue(r.Context(), CtxUserID, uint(uid))
        ctx = context.WithValue(ctx, CtxUserRole, role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


