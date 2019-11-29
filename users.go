package users

import (
	"github.com/jmoiron/sqlx"
)

// New function is a factory which returns users Interactor instance with injected users repository
// Can be used as a helper to make the code shorter
func New(db *sqlx.DB, tableName, signingKey string) *Interactor {
	return NewInteractor(NewRepository(db, tableName), signingKey)
}
