package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/st107853/muslib/lib"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	var config = lib.PostgresDBParams{
		Host:     os.Getenv("HOST"),
		DBName:   os.Getenv("DBNAME"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
	}

	lib.Connect(config)

	r := mux.NewRouter()

	r.HandleFunc("/muslib", musicGetAll).Methods("GET")
	r.HandleFunc("/muslib/{parametr}/{name}", musicGetBy).Methods("GET")

	r.HandleFunc("/muslib/{group}/{song}/{parametr}/{data}", musicPut).Methods("Put")

	r.HandleFunc("/muslib/{group}/{song}", musicPost).Methods("POST")

	r.HandleFunc("/muslib/{group}/{song}", musicDelete).Methods("DELETE")

	//r.HandleFunc("muslib/chang")

	//create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      r,
	}

	log.Printf("Server starts at: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func musicGetAll(w http.ResponseWriter, r *http.Request) {
	mus, err := lib.Get()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(mus)

	resJson, _ := json.Marshal(mus)

	w.Write(resJson)
}

func musicGetBy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mus, err := lib.GetBy(vars["parametr"], vars["name"])

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(mus)

	resJson, _ := json.Marshal(mus)

	w.Write(resJson)
}

func musicPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := lib.Put(vars["group"], vars["song"], vars["parametr"], vars["data"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
}

func musicPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := lib.Post(vars["group"], vars["song"])

	if err != nil {
		w.WriteHeader(503)
		log.Print(err)
	}
	w.WriteHeader(200)
}

func musicDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := lib.Delate(vars["group"], vars["song"])

	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(200)
}
