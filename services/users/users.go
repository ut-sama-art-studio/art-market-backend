package users

import (
	"log"

	"github.com/ut-sama-art-studio/art-market-backend/database"
	"github.com/ut-sama-art-studio/art-market-backend/services/auth"
)

type User struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
	Bio            *string `json:"bio,omitempty"`
}

func (user User) Save() (id string) {
	query := `
		INSERT INTO "USER" (Name, Email, Password, ProfilePicture, Bio)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	passwordHash, err := auth.Hash(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Db.QueryRow(query, user.Name, user.Email, passwordHash, user.ProfilePicture, user.Bio).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Row inserted!")
	return id
}
