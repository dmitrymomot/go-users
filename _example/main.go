package main

import (
	"encoding/json"
	"fmt"
	"log"

	users "github.com/dmitrymomot/go-users"
	"github.com/dmitrymomot/random"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // or another driver
)

func main() {
	db, err := sqlx.Connect("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		disabled INTEGER,
		confirmed INTEGER,
		created_at INTEGER NOT NULL,
		updated_at INTEGER
	);
	`)

	ui := users.NewInteractor(users.NewRepository(db, "users"), "secret%Key")

	d := &users.User{
		Email: random.String(10) + "@dmomot.com",
	}
	d.SetPassword("passwd")
	if err = ui.Create(d); err != nil {
		panic(err)
	}

	u, err := ui.GetByEmail(d.Email)
	if err != nil {
		panic(err)
	}
	log.Printf("user %+v", u)
	log.Printf("check password (wrong) %v", ui.VerifyPassword(u, "secret"))
	log.Printf("check password (valid) %v", ui.VerifyPassword(u, "passwd"))

	token, err := ui.ConfirmationToken(u, 60)
	if err != nil {
		panic(err)
	}
	if err := ui.Confirm(token); err != nil {
		panic(err)
	}

	u2, err := ui.GetByID(u.ID)
	if err != nil {
		panic(err)
	}
	log.Printf("user %+v", u2)

	b, err := json.Marshal(u2)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
