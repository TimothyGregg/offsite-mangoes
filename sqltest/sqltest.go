package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"io/ioutil"
)

type PlaylistSong struct {
	PlaylistID string
	SongID string
}

func main() {
	db, err := sql.Open("mysql", "root:" + getPass() + "@(localhost:3306)/spotify?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT PlaylistID, SongID FROM playlists_songs"
	rows, err := db.Query(query)

	var ps PlaylistSong

	rows.Next()
	err = rows.Scan(&ps.PlaylistID, &ps.SongID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ps.PlaylistID + " : " + ps.SongID)
}

func getPass() (file_data string) {
	data, err := ioutil.ReadFile("../local")

	if err != nil {
		log.Fatal(err)
	}

	file_data = string(data)
	return
}