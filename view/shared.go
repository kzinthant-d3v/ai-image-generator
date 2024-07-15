package view

import (
	"context"
	"fmt"
	"kzinthant-d3v/ai-image-generator/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	fmt.Println(user)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
