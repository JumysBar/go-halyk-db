package main

import (
	"database/sql"
	"fmt"
)

func dbInsertUsers(db *sql.DB, users []*User) error {
	query := "INSERT INTO users(Name, favorite_book) VALUES ($1, $2)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	defer stmt.Close()

	var result sql.Result
	for _, user := range users {
		if user.FavoriteBook == "" {
			result, err = stmt.Exec(user.Name, nil)
		} else {
			result, err = stmt.Exec(user.Name, user.FavoriteBook)
		}
		if err != nil {
			return fmt.Errorf("dbInsertUser: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("dbInsertUser: %w", err)
		}
		fmt.Printf("Rows affected: %d\n", affected)

	}
	return nil
}
