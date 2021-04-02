package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"io/ioutil"
)

func main() {
	db, err := sql.Open("mysql", "root:" + *(getPass()) + "@(localhost:3306)/spotify?parseTime=true")

	fmt.Println("Here")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Here Too")

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT PlaylistID, SongID FROM playlists_songs"
	rows, err := db.Query(query)

	var p, s string

	rows.Next()
	err = rows.Scan(&p, &s)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p + " : " + s)
}

func getPass() *string {
	data, err := ioutil.ReadFile("../password")

	if err != nil {
		log.Fatal(err)
	}

	out:= string(data)

	return &out
}