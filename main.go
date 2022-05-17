package main

import (
	"dota2_fantasy/src/repo"
	"dota2_fantasy/src/router"
	"dota2_fantasy/src/service"
	"dota2_fantasy/src/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	config := util.LoadSecrets()
	dbConnection := util.NewDBConnection(config.Secrets)

	if dbErr := dbConnection.EstablishConnection(); dbErr != nil {
		panic(dbErr)
	}

	pool := dbConnection.GetPool()
	defer pool.Close()

	r := mux.NewRouter()

	mw := router.NewMiddleware(dbConnection)
	repos := repo.SetupRepos()
	services := service.SetupServices(config, repos)
	router.SetupRouters(config, mw, services, r)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println(("working!"))

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	log.Println(("goodbye"))

}
