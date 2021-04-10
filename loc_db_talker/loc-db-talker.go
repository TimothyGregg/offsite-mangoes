package loc_db_talker

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getPass() (file_data string) {
	data, err := ioutil.ReadFile("passwords/local")

	if err != nil {
		log.Fatal(err)
	}

	file_data = string(data)
	return
}

func Table_Reader(c chan []interface{}, table_name string) {
	db, err := sql.Open("mysql", "root:"+getPass()+"@(localhost:3306)/spotify?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SELECT * FROM " + table_name)

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// This is a wacky workaround to meet the expected input type of rows.Scan()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	// Fill scanArgs with pointers to values, so we can output the pointers to the sql.RawBytes
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// This is dumb. I should really pass back "values" but the compiler doesn't like that.
		// OK SO this needs to go at the top of the function because concurrency.
		// When the channel is asked for a value, it returns the array of just pointers AND THEN gets the first values
		// If it is at the end, it gets the first values, returns them, then gets the second before the first are used
		// This is because I am returning pointers, not values. This wouldn't be an issue if it were values.
		c <- scanArgs
		err = rows.Scan(scanArgs...)
		fmt.Println(string(values[1]))
		time.Sleep(time.Second)
		if err != nil {
			panic(err.Error())
		}
	}

	close(c)
}
