package users

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ut-sama-art-studio/art-market-backend/database"
)

// this User struct models database object, where as the one in model.user is a data transfer object
type User struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
}

func (user User) Insert() (string, error) {
	// use " " around names with uppercase
	query := `
		INSERT INTO "User" (name, email, password, "profilePicture", bio)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id string
	err := database.Db.QueryRow(query, user.Name, user.Email, user.Password, user.ProfilePicture, user.Bio).Scan(&id)
	if err != nil {
		log.Print("Error executing query: ", err)
		return "", err
	}

	log.Print("Row inserted!")
	return id, nil
}

func GetUserByID(id string) (*User, error) {
	query := `
		SELECT id, name, email, password, "profilePicture", bio
		FROM "User"
		WHERE id = $1
	`

	var user User
	err := database.Db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		log.Print("Error executing query: ", err)
		return nil, err
	}

	return &user, nil
}

// updates the user based on the user object
func (user *User) Update() error {
	// Construct SQL query to update user information
	query := `
		UPDATE "User"
		SET name = $1, email = $2, "profilePicture" = $3, bio = $4
		WHERE id = $5
	`

	// Execute the query using the database connection
	_, err := database.Db.Exec(query, user.Name, user.Email, user.ProfilePicture, user.Bio, user.ID)
	if err != nil {
		log.Print("Error updating user: ", err)
		return err
	}

	return nil
}

// deletes the user
func DeleteById(id string) error {
	// Construct SQL query to delete user
	query := `
		DELETE FROM "User"
		WHERE id = $1
	`

	// Execute the query using the database connection
	_, err := database.Db.Exec(query, id)
	if err != nil {
		log.Print("Error deleting user: ", err)
		return err
	}

	return nil
}
