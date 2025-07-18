package middlewares

import (
    "context"
    "net/http"
    "strings"
    "fmt"
    
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("TU_SECRETO_SUPER_SEGURO")

type key string

const userCtxKey key = "user"

func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }


        if !hasPermission(claims, r.URL.Path) {
            http.Error(w, "Permission denied", http.StatusForbidden)
            return
        }

        ctx := context.WithValue(r.Context(), userCtxKey, claims["email"])
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func hasPermission(claims jwt.MapClaims, path string) bool {

    return true
}

func GetUserEmail(ctx context.Context) string {
    userEmail, _ := ctx.Value(userCtxKey).(string)
    return userEmail
}
