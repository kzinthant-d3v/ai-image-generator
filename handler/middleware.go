package handler

import (
	"context"
	"kzinthant-d3v/ai-image-generator/types"
	"net/http"
	"strings"
)

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	if rv := r.Context().Value(types.UserContextKey); rv != nil {
		if user, ok := rv.(types.AuthenticatedUser); ok {
			return user
		}
	}
	return types.AuthenticatedUser{}
}

func WithAuthUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
