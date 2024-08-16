package users

import (
	"database/sql"
	"log"

	"github.com/ut-sama-art-studio/art-market-backend/database"
)

// this User struct models database object, where as the one in model.user is a data transfer object
type User struct {
	ID      string `json:"id"`
	OauthID string `json:"oauth_id"`
	Role    string `json:"role"`

	Username       string  `json:"username,omitempty"`
	Name           string  `json:"name"` // display name
	Email          *string `json:"email,omitempty"`
	Password       *string `json:"password,omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
}

func (user User) Insert() (string, error) {
	// SQL query to insert a new user into the User table
	query := `
        INSERT INTO "User" (oauth_id, username, name, email, password, profile_picture, bio)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	// log.Printf("Inserting user: oauth_id=%s, name=%s, email=%s, password=%s, profile_picture=%s, bio=%s",
	// 	user.OauthID,
	// 	user.Name,
	// 	user.Email,
	// 	user.Password,
	// 	user.ProfilePicture,
	// 	user.Bio,
	// )
	var id string
	// Execute the query and get the newly inserted user's ID
	err := database.Db.QueryRow(query, user.OauthID, user.Username, user.Name, user.Email, user.Password, user.ProfilePicture, user.Bio).Scan(&id)
	if err != nil {
		log.Print("Error executing query: ", err)
		return "", err
	}

	log.Print("User created: ", id)
	return id, nil
}

func GetUserByOauthID(oauthID string) (*User, error) {
	query := `
		SELECT id, username, name, email, password, profile_picture, bio
		FROM "User"
		WHERE oauth_id = $1
	`

	var user User
	err := database.Db.QueryRow(query, oauthID).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Print("Error executing query: ", err)
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id string) (*User, error) {
	query := `
		SELECT id, username, name, email, password, profile_picture, bio
		FROM "User"
		WHERE id = $1
	`

	var user User
	err := database.Db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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
		SET name = $1, email = $2, profile_picture = $3, bio = $4
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

// retrieves all users from the database
func GetAllUsers() ([]User, error) {
	query := `
		SELECT id, username, name, email, password, profile_picture, bio
		FROM "User"
	`

	rows, err := database.Db.Query(query)
	if err != nil {
		log.Print("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Password, &user.ProfilePicture, &user.Bio)
		if err != nil {
			log.Print("Error scanning row: ", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Print("Error with rows: ", err)
		return nil, err
	}

	return users, nil
}
