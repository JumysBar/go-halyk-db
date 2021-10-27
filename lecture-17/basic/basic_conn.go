package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	username = "postgres"
	password = "postgres"
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

type User struct {
	ID           int64
	Name         string
	Created      time.Time
	FavoriteBook string
}

func dbInsertUser(db *sql.DB, user *User) error {
	var result sql.Result
	var err error
	query := "INSERT INTO users(Name, favorite_book) VALUES ($1, $2)"
	if user.FavoriteBook == "" {
		result, err = db.Exec(query, user.Name, nil)
	} else {
		result, err = db.Exec(query, user.Name, user.FavoriteBook)
	}
	if err != nil {
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("dbInsertUser: %w", err)
	}
	// lastID, err := result.LastInsertId()
	// if err != nil {
	// 	return fmt.Errorf("dbInsertUser: %w", err)
	// }
	fmt.Printf("Rows affected: %d\n", affected)
	return nil
}

func dbGetUsers(db *sql.DB) ([]*User, error) {
	var (
		id                 int64
		name               string
		created            time.Time
		favoriteBook       sql.NullString
		favoriteBookString string
	)

	result := make([]*User, 0, 10)

	rows, err := db.Query("SELECT ID, Name, created_at, favorite_book FROM users")
	if err != nil {
		return nil, fmt.Errorf("dbGetUsers: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &created, &favoriteBook); err != nil {
			return nil, fmt.Errorf("dbGetUsers: %w", err)
		}
		if !favoriteBook.Valid {
			fmt.Printf("Favorite book for user %s is null\n", name)
			favoriteBookString = "Unknown"
		} else {
			favoriteBookString = favoriteBook.String
		}
		result = append(result, &User{
			ID:           id,
			Name:         name,
			Created:      created,
			FavoriteBook: favoriteBookString,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dbGetUsers: %w", err)
	}
	return result, nil
}

func dbGetUsersCount(db *sql.DB) (int, error) {
	var result int
	if err := db.QueryRow("SELECT count(*) FROM users").Scan(&result); err != nil {
		return 0, err
	}
	return result, nil
}

func dbDeleteUsers(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	return err
}

func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := db.Ping(); err != nil {
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

	insertUsers := []*User{
		{
			Name:         "Wowah",
			FavoriteBook: "Самоучитель корейского",
		},
		{
			Name:         "Yourblizzx",
			FavoriteBook: "Some book",
		},
		{
			Name:         "Anuar",
			FavoriteBook: "Some book",
		},
	}

	if err := dbInsertUsers(db, insertUsers); err != nil {
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
		fmt.Printf("ID: %d. Name: %s. Favorite book: %s. Created at: %v\n", user.ID, user.Name, user.FavoriteBook, user.Created)
	}

	// insertTransactionUsers := []*User{
	// 	{
	// 		Name: "Zhannur",
	// 	},
	// 	{
	// 		Name: "Alex",
	// 	},
	// 	{
	// 		Name: "Wowah",
	// 	},
	// }

	// if err := dbTransaction(db, insertTransactionUsers); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Transaction completed")
	// }

	// fmt.Println("------------------------------------------------------------------------")

	// users, err = dbGetUsers(db)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for _, user := range users {
	// 	fmt.Printf("ID: %d. Name: %s. Created at: %v\n", user.ID, user.Name, user.Created)
	// }

	if err := dbDeleteUsers(db); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All commands were successfully completed")
}
