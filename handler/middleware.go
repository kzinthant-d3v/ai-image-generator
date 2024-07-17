package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kzinthant-d3v/ai-image-generator/db"
	"kzinthant-d3v/ai-image-generator/pkg/sb"
	"kzinthant-d3v/ai-image-generator/types"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	if rv := r.Context().Value(types.UserContextKey); rv != nil {
		if user, ok := rv.(types.AuthenticatedUser); ok {
			return user
		}
	}
	return types.AuthenticatedUser{}
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)
		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func WithAccountSetup(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getAuthenticatedUser(r)
		account, err := db.GetAccountByID(user.ID)
		// The user has not setup his account yet.
		// Hence, redirect him to /account/setup
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("there is some errors")
				http.Redirect(w, r, "/account/setup", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("in here!!!!!")
		user.Account = account
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithAuthUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
		session, err := store.Get(r, "user")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		accessToken := session.Values["accessToken"]
		if accessToken == nil {
			next.ServeHTTP(w, r)
			return
		}
		// accessToken, err := r.Cookie("at")
		// if err != nil {
		// 	fmt.Println(err)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		res, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:       uuid.MustParse(res.ID),
			Email:    res.Email,
			LoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
