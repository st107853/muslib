package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/st107853/muslib/lib"
)

type Data struct {
	Title  string
	Musics []lib.Music
}

var tmpl, _ = template.ParseFiles("template.html")
var index, _ = template.ParseFiles("index.html")

func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	port := os.Getenv("PORT")

	err = lib.Connect()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/muslib", musicGetAll).Methods("GET")
	r.HandleFunc("/muslib/{parametr}/{data}", musicGetBy).Methods("GET")
	r.HandleFunc("/muslib/song/{group}/{song}", musicGetSong).Methods("GET")

	r.HandleFunc("/muslib/{group}/{song}/{parametr}/{data}", musicPut).Methods("PUT")

	r.HandleFunc("/muslib/{group}/{song}", musicPost).Methods("POST")

	r.HandleFunc("/muslib/{group}/{song}", musicDelete).Methods("DELETE")

	//create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      r,
	}

	slog.Any("Server starts at: %v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func musicGetAll(w http.ResponseWriter, r *http.Request) {
	mus, err := lib.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"can't show all songs",
			slog.String("method", "GET"),
			slog.String("path", "/muslib"),
			slog.Any("error", err),
		)
		return
	}

	data := Data{
		Title:  "All what we have",
		Musics: mus,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"template.html problem",
			slog.String("method", "GET"),
			slog.String("path", "/muslib"),
			slog.Any("error", err),
		)
		return
	}
}

func musicGetBy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mus, err := lib.GetBy(vars["parametr"], vars["data"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"lib.GetBy problem",
			slog.String("method", "GET"),
			slog.String("path", "/muslib/"+vars["parametr"]+"/"+vars["data"]),
			slog.Any("error", err),
		)
		return
	}

	data := Data{
		Title:  "Songs of band",
		Musics: mus,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"template.html problem",
			slog.String("method", "GET"),
			slog.String("path", "/muslib/{parametr}/{data}"),
			slog.Any("error", err),
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func musicGetSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group, song := vars["group"], vars["song"]

	mus, err := lib.GetSong(group, song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"lib.GetSong problem",
			slog.String("method", "GET"),
			slog.String("path", "/muslib/song/"+vars["group"]+"/"+vars["song"]),
			slog.Any("error", err),
		)
		return
	}

	title := fmt.Sprintf("%v %v text", group, song)

	err = index.Execute(w, Data{Title: title, Musics: []lib.Music{mus}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Warn(
			"ingex.html problem",
			slog.String("method", "GET"),
			slog.String("path", "/muslib/song/{group}/{song}"),
			slog.Any("error", err),
		)
		return
	}
}

func musicPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := lib.Put(vars["group"], vars["song"], vars["parametr"], vars["data"])

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Warn(
			"lib.Put problem",
			slog.String("method", "PUT"),
			slog.String("path", "/muslib/"+vars["group"]+"/"+vars["song"]+
				"/"+vars["parametr"]+"/"+vars["data"]),
			slog.Any("error", err),
		)
		return
	}

	w.Write([]byte("<h1>Successesfullu updated!</h1>"))
}

func musicPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := time.Now()

	err := lib.Post(vars["group"], vars["song"], t.Format("2006.01.02"))

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Warn(
			"lib.Put problem",
			slog.String("method", "POST"),
			slog.String("path", "/muslib/"+vars["group"]+"/"+vars["song"]),
			slog.Any("error", err),
		)
		log.Print(err)
	}
	w.Write([]byte("<h1>Successesfullu created!</h1>"))
}

func musicDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := lib.Delate(vars["group"], vars["song"])

	if err != nil {
		slog.Warn(
			"lib.Put problem",
			slog.String("method", "DELETE"),
			slog.String("path", "/muslib/"+vars["group"]+"/"+vars["song"]),
			slog.Any("error", err),
		)
	}

	w.WriteHeader(http.StatusOK)
}
