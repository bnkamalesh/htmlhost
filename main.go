package main

import (
	"log"

	"github.com/bnkamalesh/htmlhost/internal/api"
	"github.com/bnkamalesh/htmlhost/internal/configs"
	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/htmlhost/internal/server/http"
)

func main() {
	cfgSvc := configs.New()

	pgSvc, err := pages.NewService(cfgSvc.Pages())
	if err != nil {
		log.Fatal(err)
		return
	}

	apiSvc := api.New(pgSvc)

	http, err := http.New(
		cfgSvc.HTTP(),
		apiSvc,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(http.Start())
}
