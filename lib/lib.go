package lib

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var ErrorNoSuchSong = errors.New("No such song")
var ErrorAlreadyExists = errors.New("Already exists")

type Music struct {
	Group       string `json:"group name"`
	Song        string `json:"song name"`
	ReleaseDate string `json:"date of release"`
	Text        string `json:"text"`
	Link        string `json:"youtube link"`
}

type PostgresDBParams struct {
	DBName   string
	Host     string
	User     string
	Password string
}

type Logger struct {
	//	events chan<- Event // Write-only channel for sending events
	//	errors <-chan error // Read-only channels for receving errors
	*gorm.DB // The database access interface
}

var db Logger

func Connect(config PostgresDBParams) error {

	//Capture connection propeties.
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		config.Host, config.DBName, config.User, config.Password)

	//Get a database handle.
	l, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("faild to open db: %w", err)
	}

	db = Logger{l}

	err = db.AutoMigrate(&Music{})
	if err != nil {
		fmt.Print("faild to open db: ", err.Error())
	}

	fmt.Println("Connected!")
	return nil
}

func Get() ([]Music, error) {
	var songs []Music

	db.Find(&songs)
	if db.Error != nil {
		return nil, fmt.Errorf("get songs err: %w", db.Error)
	}
	return songs, nil
}

func GetBy(par, data string) ([]Music, error) {
	var songs []Music

	query := fmt.Sprintf(`"%v" = ?`, par)

	db.Find(&songs, query, data)
	if len(songs) == 0 {
		return nil, ErrorNoSuchSong
	}
	return songs, nil
}

func Put(group, song, parametr, data string) error {
	var temp Music

	db.Where(`"group" = ? AND "song" = ?`, group, song).First(&temp)

	if temp.Group == "" {
		return ErrorNoSuchSong
	}

	switch parametr {
	case "group":
		temp.Group = data
	case "song":
		temp.Song = data
	case "releasedate":
		temp.ReleaseDate = data
	}

	db.Where(`"group" = ? AND "song" = ?`, group, song).Save(&temp)

	return nil
}

func Post(group, song string) error {
	var temp Music

	err := db.Where(Music{Group: group, Song: song}).FirstOrCreate(&temp).Error

	return err
}

func Delate(group, song string) error {
	var songs []Music

	db.Delete(&songs, Music{Group: group, Song: song})
	if db.Error != nil {
		return fmt.Errorf("get by date err: %w", db.Error)
	}

	return nil
}
