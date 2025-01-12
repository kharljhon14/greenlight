package main

import (
	"context"
	"net/http"

	"github.com/kharljhon14/greenlight/internal/data"
)

type contextKey string

const useContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), useContextKey, user)

	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(useContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
