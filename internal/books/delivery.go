package books

import "net/http"

type Handlers interface {
	List() http.HandlerFunc
	Create() http.HandlerFunc
	Get() http.HandlerFunc
	Update() http.HandlerFunc
}
