package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users\\(Name, favorite_book\\)").WithArgs("Nomi", nil).WillReturnResult(sqlmock.NewResult(1, 1))
	mockRows := sqlmock.NewRows([]string{"ID", "Name", "created_at", "favorite_book"}).AddRow(1, "Nomi", time.Now(), sql.NullString{})
	mock.ExpectQuery("SELECT").WillReturnRows(mockRows)

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
