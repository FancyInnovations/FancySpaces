package app

import (
	"net/http"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux *http.ServeMux
}

func Start(cfg Configuration) {

}
