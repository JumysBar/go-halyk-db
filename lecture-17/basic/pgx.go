package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	username = "postgres"
	password = "postgres"
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

type User struct {
	ID      int64
	Name    string
	Created time.Time
}

func dbInsertUser(db *pgxpool.Pool, user *User) error {
	query := "INSERT INTO users(Name) VALUES ($1)"
	// Очень очень плохо
	result, err := db.Exec(context.Background(), query, user.Name)
	if err != nil {
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	affected := result.RowsAffected()
	fmt.Printf("Rows affected: %d\n", affected)
	return nil
}

func dbGetUsers(db *pgxpool.Pool) ([]*User, error) {
	var (
		id      int64
		name    string
		created time.Time
	)

	result := make([]*User, 0, 10)

	rows, err := db.Query(context.Background(), "SELECT ID, Name, _created FROM users")
	if err != nil {
		return nil, fmt.Errorf("dbGetUsers: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &created); err != nil {
			return nil, fmt.Errorf("dbGetUsers: %w", err)
		}
		result = append(result, &User{
			ID:      id,
			Name:    name,
			Created: created,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dbGetUsers: %w", err)
	}
	return result, nil
}

func dbGetUsersCount(db *pgxpool.Pool) (int, error) {
	var result int
	if err := db.QueryRow(context.Background(), "SELECT count(*) FROM users").Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func dbDeleteUsers(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), "DELETE FROM users")
	return err
}

func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)

	config, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		fmt.Println(err)
		return
	}

	config.MaxConns = 25
	config.MaxConnLifetime = 5 * time.Minute

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Ping(context.Background()); err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	fmt.Println("Successfully connected to postgres")

	err = dbInsertUser(db, &User{
		Name: "Nomi",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	count, err := dbGetUsersCount(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Users count: %d\n", count)

	users, err := dbGetUsers(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d. Name: %s. Created at: %v\n", user.ID, user.Name, user.Created)
	}

	if err := dbDeleteUsers(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All commands were successfully completed")
}
