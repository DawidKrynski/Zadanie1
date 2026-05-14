package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name"`
	PasswordHash  string         `json:"-"`
	OAuthAccounts []OAuthAccount `json:"oauth_accounts,omitempty" gorm:"foreignKey:UserID"`
	Sessions      []AuthSession  `json:"-" gorm:"foreignKey:UserID"`
}

type OAuthAccount struct {
	gorm.Model
	UserID         uint      `json:"user_id" gorm:"index;not null"`
	Provider       string    `json:"provider" gorm:"uniqueIndex:idx_oauth_provider_user;not null"`
	ProviderUserID string    `json:"provider_user_id" gorm:"uniqueIndex:idx_oauth_provider_user;not null"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	TokenType      string    `json:"-"`
	Expiry         time.Time `json:"expiry"`
	User           User      `json:"-" gorm:"foreignKey:UserID"`
}

type AuthSession struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	TokenHash string    `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}

type OAuthState struct {
	gorm.Model
	Provider  string    `json:"provider" gorm:"not null"`
	StateHash string    `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at"`
}
