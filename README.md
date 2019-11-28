# go-users

Simple implementation of users interaction

## Installation

```
go get -u github.com/dmitrymomot/go-users
```

## Usage

```
package main

import (
    "github.com/dmitrymomot/go-users"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3" // or another driver
)

func main() {
    // DB connection
    db, err := sqlx.Connect("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Init user repository
    // Table `users` must be exists in database
    ur := users.NewRepository(db, "users")
    // Init user interactor
    ui := users.NewInteractor(ur, "your secret signing key")

    // your code ...
}

```