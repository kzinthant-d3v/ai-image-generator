package handler

import (
	"fmt"
	"kzinthant-d3v/ai-image-generator/db"
	"kzinthant-d3v/ai-image-generator/pkg/kit/validate"
	"kzinthant-d3v/ai-image-generator/pkg/sb"
	"kzinthant-d3v/ai-image-generator/pkg/utils"
	"kzinthant-d3v/ai-image-generator/types"
	"kzinthant-d3v/ai-image-generator/view/auth"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.AccountSetup())
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	// cookie := &http.Cookie{
	// 	Value:    "",
	// 	Name:     "at",
	// 	MaxAge:   -1,
	// 	HttpOnly: true,
	// 	Path:     "/",
	// 	Secure:   true,
	// }
	// http.SetCookie(w, cookie)

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, "user")
	session.Values["accessToken"] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func HandleLoginInIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.LogIn())
}

func HandleSignUpIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignUp())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	errors := auth.SignupErrors{}
	if ok := validate.New(&params, validate.Fields{
		"Email":           validate.Rules(validate.Required, validate.Email),
		"Password":        validate.Rules(validate.Required, validate.Password),
		"ConfirmPassword": validate.Rules(validate.Equal(params.Password), validate.Message("The password and confirm password must be the same")),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignupForm(params, errors))
	}
	_, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	fmt.Println(err)
	if err != nil {
		return nil
	}
	return render(r, w, auth.SignUpSuccess(params.Email))
}

func HandleLoginInWithGoogle(w http.ResponseWriter, r *http.Request) error {
	res, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "https://kaskar.xyz/auth/callback",
	})
	if err != nil {
		return nil
	}
	http.Redirect(w, r, res.URL, http.StatusSeeOther)
	return nil
}

func HandleLoginInCreate(w http.ResponseWriter, r *http.Request) error {
	crendentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if !utils.IsValidEmail(crendentials.Email) {
		return render(r, w, auth.LoginForm(crendentials, auth.LoginErrors{
			Email: "The email is not valid",
		}))
	}

	res, err := sb.Client.Auth.SignIn(r.Context(), crendentials)

	if err != nil {
		return render(r, w, auth.LoginForm(crendentials, auth.LoginErrors{
			InvalidCredentials: "The credentials are not correct",
		}))
	}

	if err := setAuthSession(w, r, res.AccessToken); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}

	setAuthSession(w, r, accessToken)
	// setAuthCookie(w, accessToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, "user")
	session.Values["accessToken"] = accessToken
	return session.Save(r, w)
}

func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupParams{
		Username: r.FormValue("username"),
	}

	var errors auth.AccountSetupErrors

	if ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Required, validate.Min(2), validate.Max(50)),
	}).Validate(&errors); !ok {
		return render(r, w, auth.AccountSetupForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		Username: params.Username,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")
}
