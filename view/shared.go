package view

import (
	"context"
	"kzinthant-d3v/ai-image-generator/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
