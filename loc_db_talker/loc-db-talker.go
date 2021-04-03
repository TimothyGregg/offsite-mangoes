package loc_db_talker

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"io/ioutil"
)

type PlaylistSong struct {
	PlaylistID string
	SongID string
}

func Songs_Table_Reader(c chan PlaylistSong) {
	db, err := sql.Open("mysql", "root:" + getPass() + "@(localhost:3306)/spotify?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT PlaylistID, SongID FROM playlists_songs"
	rows, err := db.Query(query)

	var ps PlaylistSong

	more := rows.Next()
	for {
		// break on end of query results
		if !more {
			close(c)
			break
		}
		err = rows.Scan(&ps.PlaylistID, &ps.SongID)
		if err != nil {
			log.Fatal(err)
		}
		more = rows.Next()
		c <- ps
	}
}

func getPass() (file_data string) {
	data, err := ioutil.ReadFile("passwords/local")

	if err != nil {
		log.Fatal(err)
	}

	file_data = string(data)
	return
}