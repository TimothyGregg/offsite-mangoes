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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Atlas connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://tim:"+url.QueryEscape(getPass())+"@cluster0.kkwum.mongodb.net/the_db?retryWrites=true&w=majority",
	))
	if err != nil { panic(err) }

	defer func() {
		if err = client.Disconnect(ctx); err != nil { panic(err) }
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil { panic(err) }

	fmt.Println("Successfully connected and pinged.")

	c := make(chan []interface{})
	// goroutine for adding things to the channel. Without this, the loop hangs forever
	go loc_db_talker.Table_Reader(c, "tracks")
	var outstring string
	var field_names []string
	var input []string
	var filter bson.D
	var update bson.D
	initial := true
	for interface_array := range c {
		if initial {
			field_names = make([]string, len(interface_array))
			input = make([]string, len(interface_array))
			initial = false
			for i, pointer := range interface_array {
				// Extraction of value from interface requires type assertion. 
				field_names[i] = string(*pointer.(*sql.RawBytes))
			}
			filter = bson.D{field_names}
		} else {
			for i, pointer := range interface_array {
				// Extraction of value from interface requires type assertion. 
				input[i] = pointer.(*sql.RawBytes)
			}
			for _, p := range pointers {
				outstring += string(*p) + " / "
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
