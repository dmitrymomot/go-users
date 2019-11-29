package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	// q, p := users.ConditionsToQuery(
	// 	users.OrderBy(users.CreatedAtDesc, users.UpdatedAtDesc),
	// 	users.Offset(2),
	// 	users.Limit(23),
	// 	users.Confirmed(false),
	// 	users.Disabled(true),
	// )
	// fmt.Printf("query: %s; params: %+v", q, p)

	ul, err := ui.GetList(
		users.OrderBy(users.CreatedAtDesc, users.UpdatedAtDesc),
		// users.Offset(20),
		users.Limit(3),
		// users.Confirmed(true),
		// users.Disabled(false),
	)
	if err != nil {
		panic(err)
	}
	log.Printf("users count %d", len(ul))

	for _, i := range ul {
		fmt.Println(fmt.Printf("user: %s created at: %s", i.Email, time.Unix(i.CreatedAt, 0)))
	}
}
