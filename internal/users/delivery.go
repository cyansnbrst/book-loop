package users

import "net/http"

type Handlers interface {
	Register() http.HandlerFunc
	Activate() http.HandlerFunc
	Login() http.HandlerFunc
	CreateActivationToken() http.HandlerFunc
}
