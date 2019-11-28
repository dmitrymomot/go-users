package users

import (
	"time"

	"github.com/dmitrymomot/go-signature"
	"github.com/dmitrymomot/go-utilities/hash"
	"github.com/google/uuid"
)

type (
	// Interactor structure
	Interactor struct {
		repository Repository
	}

	// UserRepository interface
	UserRepository interface {
		GetByID(id string) (*User, error)
		GetByEmail(email string) (*User, error)
		GetList(limit, offset int, sort string, where ...interface{}) ([]*User, error)
		Insert(*User) error
		Update(*User) error
		Delete(id string) error
	}
)

// NewInteractor factory
func NewInteractor(r Repository, signingKey string) *Interactor {
	if signingKey == "" {
		signingKey = "secret%key"
	}
	signature.SetSigningKey(signingKey)
	return &Interactor{repository: r}
}

// GetByID fetch user by id
func (i *Interactor) GetByID(id string) (*User, error) {
	return i.repository.GetByID(id)
}

// GetByEmail fetch user by email
func (i *Interactor) GetByEmail(email string) (*User, error) {
	return i.repository.GetByEmail(email)
}

// GetList od users with sorting and optional conditional
func (i *Interactor) GetList(condition ...interface{}) ([]*User, error) {
	return i.repository.GetList(condition...)
}

// Create new user
func (i *Interactor) Create(u *User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.Email == "" {
		return ErrEmailMissed
	}
	if _, err := i.repository.GetByEmail(u.Email); err == nil {
		return ErrTakenEmail
	}
	return i.repository.Insert(u)
}

// Update existed user
func (i *Interactor) Update(u *User) error {
	if u.ID == "" {
		return ErrNotExistedUser
	}
	if u.UpdatedAt == nil {
		t := time.Now()
		u.UpdatedAt = &t
	}
	if u.Email == "" {
		return ErrEmailMissed
	}
	return i.repository.Insert(u)
}

// Delete user by id
func (i *Interactor) Delete(id string) error {
	return i.repository.Delete(id)
}

// VerifyPassword compare password string with user password hash,
// returns nil if password is valid
func (i *Interactor) VerifyPassword(u *User, passwordStr string) error {
	if err := hash.Compare(u.Password, passwordStr); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// ConfirmationToken returns confirmation token string
func (i *Interactor) ConfirmationToken(u *User, ttl int64) (string, error) {
	token, err := signature.NewTemporary(u.Email, ttl)
	if err != nil {
		return "", ErrCouldNotGenerateToken
	}
	return token, nil
}

// Confirm function checks token and set confirmed=true if it's valid
func (i *Interactor) Confirm(token string) error {
	payload, err := signature.Parse(token)
	if err != nil {
		return ErrInvalidToken
	}
	email, ok := payload.(string)
	if !ok {
		return ErrInvalidToken
	}
	u, err := i.repository.GetByEmail(email)
	if err != nil {
		return err
	}
	u.Confirmed = true
	t := time.Now()
	u.UpdatedAt = &t
	if err := i.repository.Update(u); err != nil {
		return err
	}
	return nil
}

// ResetPasswordToken returns reset password token string
func (i *Interactor) ResetPasswordToken(u *User, ttl int64) (string, error) {
	token, err := signature.NewTemporary(u.ID, ttl)
	if err != nil {
		return "", ErrCouldNotGenerateToken
	}
	return token, nil
}

// ResetPassword fucntion validates token and updates password
func (i *Interactor) ResetPassword(token string, newPassword string) error {
	payload, err := signature.Parse(token)
	if err != nil {
		return ErrInvalidToken
	}
	uid, ok := payload.(string)
	if !ok {
		return ErrInvalidToken
	}
	u, err := i.repository.GetByID(uid)
	if err != nil {
		return err
	}
	if err := u.SetPassword(newPassword); err != nil {
		return err
	}
	t := time.Now()
	u.UpdatedAt = &t
	if err := i.repository.Update(u); err != nil {
		return err
	}
	return nil
}
