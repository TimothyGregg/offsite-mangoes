package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"github.com/TimothyGregg/offsite-mangoes/loc_db_talker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type spotify_object struct {
	OwnerID string
	Owns    string
}

func main() {
	// Atlas connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://tim:"+url.QueryEscape(getPass())+"@cluster0.kkwum.mongodb.net/the_db?retryWrites=true&w=majority",
	))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	c := make(chan []sql.RawBytes)
	input := new(spotify_object)

	for _, table_name := range [3]string{"artists", "playlists", "playlists_songs"} {
		// goroutine for adding things to the channel. Without this, the loop hangs forever
		go loc_db_talker.Table_Reader(c, table_name)
		if table_name == "playlists_songs" {
			table_name = "playlist_songs"
		}
		collection := client.Database("the_db").Collection(table_name)
		fmt.Println("Processing table: " + table_name)
		for interface_array := range c {
			input.OwnerID = string(interface_array[0])
			input.Owns = string(interface_array[1])
			_, err = collection.InsertOne(context.TODO(), input)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getPass() (file_data string) {
	data, err := ioutil.ReadFile("passwords/the_db")

	if err != nil {
		log.Fatal(err)
	}

	file_data = string(data)
	return
}
