package handler

import (
	"net/http"

	"github.com/a-h/templ"
)

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func MakeHandler(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error()+r.URL.Path, http.StatusInternalServerError)
		}
	}

}
