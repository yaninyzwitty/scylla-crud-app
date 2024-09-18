package middleware

// dont use this please

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

const AuthUserID = "middleware.auth.userID"

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		// Check that the header begins with a prefix of Bearer
		if !strings.HasPrefix(authorization, "Bearer ") {
			writeUnauthed(w)
			return
		}

		// Pull out the token
		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		// Decode the token from base 64
		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			writeUnauthed(w)
			return
		}

		// We're just assuming a valid base64 token is a valid user id.
		userID := string(token)
		fmt.Println("userID:", userID)

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
