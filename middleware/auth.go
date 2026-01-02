package middleware

import (
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/utils"
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey  contextKey = "user"
	TokenContextKey contextKey = "token"
)

func AuthMiddleware(
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
	allowedRoles ...string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set token
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				utils.JSONError(w, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			token, err := uuid.Parse(tokenStr)
			if err != nil {
				utils.JSONError(w, http.StatusUnauthorized, "Invalid token", nil)
				return
			}

			// find token dari session
			session, err := sessionRepo.FindByToken(r.Context(), token)
			if err != nil || !session.IsValid() {
				utils.JSONError(w, http.StatusUnauthorized, "Session expired", nil)
				return
			}

			// find user id dari repo
			user, err := userRepo.FindByID(r.Context(), session.UserID)
			if err != nil {
				utils.JSONError(w, http.StatusUnauthorized, "User not found", nil)
				return
			}

			if len(allowedRoles) > 0 {
				allowed := false
				for _, role := range allowedRoles {
					if user.Role == role {
						allowed = true
						break
					}
				}

				if !allowed {
					utils.JSONError(w, http.StatusForbidden, "Forbidden", nil)
					return
				}
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, TokenContextKey, tokenStr)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
