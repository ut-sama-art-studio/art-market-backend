package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/ut-sama-art-studio/art-market-backend/middlewares"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	// Check if the user ID is present in the context
	userID := middlewares.ContextUserID(ctx)
	if userID == "" {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}
