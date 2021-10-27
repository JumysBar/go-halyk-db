package main

import (
	"database/sql"
	"fmt"
)

func dbTransaction(db *sql.DB, users []*User) error {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("dbInsertUser: %w", err)
	}

	query := "INSERT INTO users(Name) VALUES ($1)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	defer stmt.Close()

	for _, user := range users {
		result, err := stmt.Exec(user.Name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("dbInsertUser: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("dbInsertUser: %w", err)
		}
		fmt.Printf("Rows affected: %d\n", affected)

	}
	// Закрепляем изменения
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	return nil
}
