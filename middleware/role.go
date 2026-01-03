package middleware

import (
	"alfdwirhmn/inventory/repository"
	"net/http"
)

const (
	RoleStaff      = "staff"
	RoleAdmin      = "admin"
	RoleSuperAdmin = "super_admin"
)

type RoleMiddleware struct {
	SessionRepo repository.SessionRepository
	UserRepo    repository.UserRepository
}

func NewRoleMiddleware(
	sessionRepo repository.SessionRepository,
	userRepo repository.UserRepository,
) *RoleMiddleware {
	return &RoleMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	}
}

// READ: staff, admin, super_admin
func (r *RoleMiddleware) AllowRead() func(http.Handler) http.Handler {
	return AuthMiddleware(
		r.SessionRepo,
		r.UserRepo,
		RoleStaff,
		RoleAdmin,
		RoleSuperAdmin,
	)
}

// ADMIN: admin, super_admin
func (r *RoleMiddleware) AllowAdmin() func(http.Handler) http.Handler {
	return AuthMiddleware(
		r.SessionRepo,
		r.UserRepo,
		RoleAdmin,
		RoleSuperAdmin,
	)
}

// SUPER ADMIN ONLY
func (r *RoleMiddleware) AllowSuperAdmin() func(http.Handler) http.Handler {
	return AuthMiddleware(
		r.SessionRepo,
		r.UserRepo,
		RoleSuperAdmin,
	)
}
