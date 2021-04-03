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
}

func getPass() (file_data string) {
	data, err := ioutil.ReadFile("../password")

	if err != nil {
		log.Fatal(err)
	}

	file_data = string(data)
	return
}