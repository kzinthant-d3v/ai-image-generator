package handler

import (
	"fmt"
	"kzinthant-d3v/ai-image-generator/pkg/kit/validate"
	"kzinthant-d3v/ai-image-generator/pkg/sb"
	"kzinthant-d3v/ai-image-generator/pkg/utils"
	"kzinthant-d3v/ai-image-generator/view/auth"
	"net/http"

	"github.com/nedpals/supabase-go"
)

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

	cookie := &http.Cookie{
		Value:    res.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
