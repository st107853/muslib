package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

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
	Link        string `json:"some link"`
}

type Logger struct {
	*gorm.DB // The database access interface
}

var db Logger

func Connect() error {

	//Capture connection propeties.
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		os.Getenv("HOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASS"))

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

func GetSong(group, song string) (Music, error) {
	var music Music

	err := db.Where(`"group" = ? AND "song" = ?`, group, song).First(&music).Error

	return music, err
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
	temp, err := GetSong(group, song)

	if temp.Group == "" || err != nil {
		return ErrorNoSuchSong
	}

	switch parametr {
	case "group":
		temp.Group = data
	case "song":
		temp.Song = data
	case "link":
		{
			url, err := base64.RawURLEncoding.DecodeString(data)

			link := string(url)
			temp.Link = link

			if err != nil {
				return err
			}
		}
	case "date":
		temp.ReleaseDate = data
	case "text":
		temp.Text = formatText(data)
	}

	db.Where(`"group" = ? AND "song" = ?`, group, song).Save(&temp)

	return nil
}

func Post(group, song, time string) error {
	var temp Music

	err := db.Where(Music{Group: group, Song: song, ReleaseDate: time}).FirstOrCreate(&temp).Error

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

func formatText(text string) string {
	result := ""

	for _, v := range text {
		if v >= 'A' && v <= 'Z' {
			result += "\n"
		}
		result += string(v)
	}

	return result
}
