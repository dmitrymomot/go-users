package users

import (
	"github.com/dmitrymomot/go-utilities/hash"
)

type (
	// User model structure
	User struct {
		ID        string `db:"id" json:"id"`
		Email     string `db:"email" json:"email"`
		Password  string `db:"password" json:"-"`
		Confirmed bool   `db:"confirmed" json:"confirmed"`
		Disabled  bool   `db:"disabled" json:"disabled"`
		CreatedAt int64  `db:"created_at" json:"created_at"`
		UpdatedAt *int64 `db:"updated_at" json:"updated_at,omitempty,"`
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
