package ext

import (
	"net/http"
)

type Server struct {
	Addr    string
	Handler http.Handler
}
