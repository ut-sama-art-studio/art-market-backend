package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/ut-sama-art-studio/art-market-backend/graph/model"
	"github.com/ut-sama-art-studio/art-market-backend/middlewares"
	"github.com/ut-sama-art-studio/art-market-backend/services/fileservice"
	"github.com/ut-sama-art-studio/art-market-backend/services/userservice"
	"github.com/ut-sama-art-studio/art-market-backend/utils/jwt"
)

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*model.User, error) {
	ctxId := middlewares.ContextUserID(ctx)
	if ctxId != id {
		// TODO: allow admin
		err := errors.New("No premission, can't update another user")
		log.Print("Error updating user: ", err)
		return nil, err
	}

	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error finding user: ", err)
		return nil, err
	}
	// update
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.ProfilePicture != nil {
		user.ProfilePicture = input.ProfilePicture
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if err = user.Update(); err != nil {
		log.Print("Error updating user: ", err)
		return nil, err
	}
	return user.ToGraphUser(), nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	ctxId := middlewares.ContextUserID(ctx)
	if ctxId != id {
		// TODO: allow admin
		err := errors.New("no premission, can't delete another user")
		log.Print("Error updating user: ", err)
		return false, err
	}

	if err := userservice.DeleteById(id); err != nil {
		log.Print("Error deleting user: ", err)
		return false, err
	}

	return true, nil
}

// UpdateProfilePicture is the resolver for the updateProfilePicture field.
func (r *mutationResolver) UpdateProfilePicture(ctx context.Context, file graphql.Upload) (*model.User, error) {
	userID := middlewares.ContextUserID(ctx)
	user, err := userservice.GetUserByID(userID)
	if err != nil {
		log.Print("Error finding user: ", err)
		return nil, err
	}

	// It's possible to just name the image with the same name to replace it, but it will cause next.js cached image to not update if it's the same URL
	folderPath := "profile-picture"
	fileservice.DeleteUserFolder(userID, folderPath)
	fileURL, err := fileservice.UploadFileToS3(file, userID, folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload profile picture: %w", err)
	}

	// update user
	user.ProfilePicture = &fileURL
	if err = user.Update(); err != nil {
		log.Print("Error updating user: ", err)
		return nil, err
	}

	return user.ToGraphUser(), nil
}

// Sets the role of the user
func (r *mutationResolver) SetRole(ctx context.Context, id string, role string) (*model.User, error) {
	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error finding user: ", err)
		return nil, err
	}
	// update
	user.Role = role
	if err = user.Update(); err != nil {
		log.Print("Error updating user: ", err)
		return nil, err
	}
	return user.ToGraphUser(), nil
}

// Apply token generated to give user artist role
func (r *mutationResolver) ApplyArtistRoleToken(ctx context.Context, token string) (*model.User, error) {
	id := middlewares.ContextUserID(ctx)

	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user in getUserByID: ", err)
		return nil, err
	}

	if user.Role == "admin" || user.Role == "director" || user.Role == "artist" {
		return user.ToGraphUser(), nil
	}

	if _, err := jwt.VerifyKeyValueToken(token, "role", "artist"); err != nil {
		// Token invalid
		return nil, err
	} else {
		// Token valid, update user role
		user.Role = "artist"
		if err = user.Update(); err != nil {
			log.Print("Error updating user: ", err)
			return nil, err
		}
		return user.ToGraphUser(), nil
	}
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	id := middlewares.ContextUserID(ctx)

	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user in getUserByID: ", err)
		return nil, err
	}
	// log.Printf("Found user: %s\n", user.Name)

	return user.ToGraphUser(), nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user in getUserByID: ", err)
		return nil, err
	}
	// log.Printf("Found user: %s\n", user.Name)

	return user.ToGraphUser(), nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	allUsers, err := userservice.GetAllUsers()
	if err != nil {
		log.Print("Error fetching users: ", err)
		return nil, err
	}

	var result []*model.User
	for _, user := range allUsers {
		result = append(result, user.ToGraphUser())
	}

	return result, nil
}

// Returns all artists
func (r *queryResolver) Artists(ctx context.Context) ([]*model.User, error) {
	allUsers, err := userservice.GetAllArtists()
	if err != nil {
		log.Print("Error fetching artists: ", err)
		return nil, err
	}

	var result []*model.User
	for _, user := range allUsers {
		result = append(result, user.ToGraphUser())
	}

	return result, nil
}

// Generate a token allowing user who apply it to become an artist
func (r *queryResolver) GenerateArtistRoleToken(ctx context.Context) (string, error) {
	id := middlewares.ContextUserID(ctx)

	user, err := userservice.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user in getUserByID: ", err)
		return "", err
	}

	if user.Role == "admin" || user.Role == "director" {
		// Token valid for 7 days
		return jwt.GenerateKeyValueToken("role", "artist", time.Now().AddDate(0, 0, 7))
	}

	return "", fmt.Errorf("no permission to create token")
}
