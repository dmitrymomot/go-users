package users

import (
	"time"

	"github.com/dmitrymomot/go-utilities/hash"
)

type (
	// User model structure
	User struct {
		ID        string     `json:"id"`
		Email     string     `json:"email"`
		Password  string     `json:"-"`
		Confirmed bool       `json:"confirmed"`
		Disabled  bool       `json:"disabled"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
	}
)

// SetPassword helper
func (u *User) SetPassword(password string) error {
	h, err := hash.New(password)
	if err != nil {
		return ErrCouldNotSetPassword
	}
	u.Password = h
	return nil
}
