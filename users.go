package users

import "database/sql"

// New function is a factory which returns users Interactor instance with injected users repository
// Can be used as a helper to make the code shorter
func New(db *sql.DB, driverName, tableName, signingKey string) *Interactor {
	return NewInteractor(NewRepository(db, driverName, tableName), signingKey)
}
