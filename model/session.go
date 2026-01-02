package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           int        `json:"id" db:"id"`
	UserID       int        `json:"user_id" db:"user_id"`
	Token        uuid.UUID  `json:"token" db:"token"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	ExpiredAt    time.Time  `json:"expired_at" db:"expired_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	LastActivity time.Time  `json:"last_activity" db:"last_activity"`
	IPAddress    string     `json:"ip_address" db:"ip_address"`
	UserAgent    string     `json:"user_agent" db:"user_agent"`
}

// IsValid checks if session is still valid
func (s *Session) IsValid() bool {
	now := time.Now()

	// Check if revoked
	if s.RevokedAt != nil {
		return false
	}

	// Check if expired
	if now.After(s.ExpiredAt) {
		return false
	}

	return true
}

// IsExpired checks if session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiredAt)
}

// IsRevoked checks if session has been revoked
func (s *Session) IsRevoked() bool {
	return s.RevokedAt != nil
}
