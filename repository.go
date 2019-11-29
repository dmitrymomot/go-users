package users

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined users list ordering
const (
	// ORDER BY created_at ASC
	CreatedAtAsc Order = iota + 1
	// ORDER BY created_at DESC
	CreatedAtDesc
	// ORDER BY updated_at ASC
	UpdatedAtAsc
	// ORDER BY updated_at DESC
	UpdatedAtDesc
)

var orderQueryMap = map[Order]string{
	CreatedAtAsc:  "created_at ASC",
	CreatedAtDesc: "created_at DESC",
	UpdatedAtAsc:  "updated_at ASC",
	UpdatedAtDesc: "updated_at DESC",
}

type (
	// Order type
	Order int

	// Repository structure is implementation of UserRepository interface
	Repository struct {
		db        *sqlx.DB
		tableName string
	}
)

func (o Order) String() string {
	if v, ok := orderQueryMap[o]; ok {
		return v
	}
	return ""
}

// NewRepository factory
func NewRepository(db *sqlx.DB, tableName string) *Repository {
	return &Repository{db: db, tableName: tableName}
}

// GetByID fetch user record by id
func (r *Repository) GetByID(id string) (*User, error) {
	q := "SELECT * FROM %s WHERE id = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	u := &User{}
	if err := r.db.Get(u, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "get user by id")
	}
	return u, nil
}

// GetByEmail fetch user record by email
func (r *Repository) GetByEmail(email string) (*User, error) {
	q := "SELECT * FROM %s WHERE email = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	u := &User{}
	if err := r.db.Get(u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "get user by email")
	}
	return u, nil
}

// GetList fetch users list
func (r *Repository) GetList(c ...Condition) ([]*User, error) {
	q := "SELECT * FROM %s "
	q = fmt.Sprintf(q, r.tableName)
	sq, params := ConditionsToQuery(c...)
	q = q + sq
	q = r.db.Rebind(q)
	ul := make([]*User, 0)
	fmt.Println(fmt.Printf("query: %s; params: %+v", q, params))
	if err := r.db.Select(&ul, q, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "get users list")
	}
	return ul, nil
}

// Insert a new user record
func (r *Repository) Insert(u *User) error {
	q := "INSERT INTO %s (`id`, `email`, `password`, `confirmed`, `disabled`, `created_at`) VALUES (?, ?, ?, ?, ?, ?);"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(
		q,
		u.ID, u.Email, u.Password,
		u.Confirmed, u.Disabled,
		u.CreatedAt,
	); err != nil {
		return errors.Wrap(err, "store user")
	}
	return nil
}

// Update existed user record
func (r *Repository) Update(u *User) error {
	q := "UPDATE %s SET `email`=?, `password`=?, `confirmed`=?, `disabled`=?, `updated_at`=? WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(
		q,
		u.Email, u.Password,
		u.Confirmed, u.Disabled,
		u.UpdatedAt,
		u.ID,
	); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "update user")
	}
	return nil
}

// Delete a user record by id
func (r *Repository) Delete(id string) error {
	q := "DELETE FROM %s WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "delete user")
	}
	return nil
}
