package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"
	"log"

	"github.com/ut-sama-art-studio/art-market-backend/graph/model"
	"github.com/ut-sama-art-studio/art-market-backend/middlewares"
	"github.com/ut-sama-art-studio/art-market-backend/services/users"
)

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUser) (*model.User, error) {
	user, err := users.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user: ", err)
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
	return &model.User{ID: user.ID, Name: user.Name, Email: user.Email, ProfilePicture: user.ProfilePicture, Bio: user.Bio}, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	if err := users.DeleteById(id); err != nil {
		log.Print("Error deleting user: ", err)
		return false, err
	}

	return true, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented: RefreshToken - refreshToken"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userId := middlewares.ContextUserID(ctx)
	user, err := users.GetUserByID(userId)
	if err != nil {
		log.Print("Error fetching user: ", err)
		return nil, err
	}

	return &model.User{
		ID:             user.ID,
		Name:           user.Name,
		Username:       &user.Username,
		Email:          user.Email,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, err := users.GetUserByID(id)
	if err != nil {
		log.Print("Error fetching user: ", err)
		return nil, err
	}
	return &model.User{ID: user.ID, Name: user.Name, Username: &user.Username, Email: user.Email, ProfilePicture: user.ProfilePicture, Bio: user.Bio}, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	allUsers, err := users.GetAllUsers()
	if err != nil {
		log.Print("Error fetching users: ", err)
		return nil, err
	}

	var result []*model.User
	for _, user := range allUsers {
		result = append(result, &model.User{
			ID:             user.ID,
			Name:           user.Name,
			Username:       &user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
			Bio:            user.Bio,
		})
	}

	return result, nil
}
