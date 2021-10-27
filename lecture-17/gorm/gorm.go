package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	username = "postgres"
	password = "postgres"
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

type User struct {
	gorm.Model
	ID           int64
	Name         string
	FavoriteBook string
}

func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, hostname, port, db)
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	db.Create(&User{
		ID:           1,
		Name:         "NewUser",
		FavoriteBook: "User book",
	})

	var user User
	db.First(&user, 1)
	fmt.Printf("User with ID 1: %+v\n", user)

	db.Model(&User{}).Where("ID=?", 1).Update("favorite_book", "Some another book")

	db.First(&user, "name=?", "NewUser")
	fmt.Printf("User with name NewUser: %+v\n", user)

	// Soft delete
	db.Delete(&user)

	db.Find(&user, 1)
	fmt.Printf("User with ID 1: %+v\n", user)

	// Permanent delete
	db.Unscoped().Delete(&user)
}
