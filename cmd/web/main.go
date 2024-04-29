package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/seemsod1/ancy/internal/config"
	"log"
	"net/http"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	if err := setup(&app); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
