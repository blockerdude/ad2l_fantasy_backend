package main

import (
	"context"
	"dota2_fantasy/repo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

var dbPool *pgxpool.Pool

func main() {

	r := mux.NewRouter()
	secrets := readSecrets()
	pool, err := pgxpool.Connect(context.Background(), secrets.DBConnectionString)

	if err != nil {
		panic(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	dbPool = pool

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

	r.HandleFunc("/hello", helloHandler)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	log.Println(("goodbye"))

}

type Secrets struct {
	DBConnectionString string `json:"db_conn_string"`
}

func readSecrets() Secrets {
	jsonFile, err := os.Open("secrets.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	secrets := Secrets{}

	err = json.Unmarshal(byteValue, &secrets)
	if err != nil {
		panic(err)
	}

	return secrets

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	confRepo := repo.NewConferenceRepo()

	// conn := r.Context().Value("dbptr")

	confs2, err := confRepo.GetAllConferences(dbPool)
	if err != nil {
		fmt.Fprintf(w, "we failed!")
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := json.Marshal(confs2)
	if err != nil {
		fmt.Fprintf(w, "we failed! %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "confs: %v", string(body))
}
