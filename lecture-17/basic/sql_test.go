package main

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	err = dbInsertUser(db, &User{
		Name: "Nomi",
	})

	if err != nil {
		t.Fatal(err)
	}

	users, err := dbGetUsers(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatalf("Error! Expected len: %d. Got: %d", 1, len(users))
	}

	if users[0].Name != "Nomi" {
		t.Fatalf("Error! Expected name: %s. Got: %s", "Nomi", users[0].Name)
	}
}
