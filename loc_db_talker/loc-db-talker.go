package loc_db_talker

import (
	"database/sql"
	"io/ioutil"
	"log"

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

func Table_Reader(c chan []sql.RawBytes, table_name string) {
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
	output := make([]sql.RawBytes, len(values))
	// Fill scanArgs with pointers to values, so we can output the pointers to the sql.RawBytes
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		// Need to make a copy so that data is not changed concurrently before it is used in another process
		copy(output, values)
		c <- output
	}

	close(c)
}
