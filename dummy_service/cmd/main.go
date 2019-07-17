package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	dummy_service "github.com/ivch/dummy-service"
	"github.com/ivch/dummy-service/config"
	"github.com/ivch/dummy-service/repository"
	"github.com/ivch/dummy-service/router"
)

//todo gracefull shytdown
//todo middlewares logging/instrumenting

func main() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.New(cfg.AWSConfig, cfg.DynamoConfig)
	if err != nil {
		log.Fatal(err)
	}

	svc := dummy_service.New(repo)
	router := router.New(svc)

	// Start HTTP server.
	if len(cfg.HTTPPort) > 0 {
		go func() {
			if err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), router.Handler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
		fmt.Println("Started HTTP server on " + cfg.HTTPPort)
	}

	// Start HTTPS server.
	if len(cfg.HTTPSPort) > 0 {
		go func() {
			if err := fasthttp.ListenAndServeTLS(fmt.Sprintf(":%s", cfg.HTTPSPort), "localhost.crt", "localhost.key", router.Handler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
		fmt.Println("Started HTTPS server on " + cfg.HTTPSPort)
	}

	select {}
}
