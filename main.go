package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/url"
	"time"
	"github.com/TimothyGregg/offsite-mangoes/loc_db_talker"
)

func main () {
	// Atlas connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
	"mongodb+srv://tim:" + url.QueryEscape(getPass()) + "@cluster0.kkwum.mongodb.net/the_db?retryWrites=true&w=majority",
	))
	if err != nil { panic(err) }

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	c := make(chan loc_db_talker.PlaylistSong)
	go loc_db_talker.Songs_Table_Reader(c)
	for ps := range c {
		fmt.Println(ps.PlaylistID + " : " + ps.SongID)
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