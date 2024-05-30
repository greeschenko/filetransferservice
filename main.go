package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User is a struct that represents a user in the database
type Files struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Democount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Demo count function started")

    fmt.Println(os.Getenv("MYSQLUSER"), os.Getenv("MYSQLPASS"), os.Getenv("MYSQLDOMEN"), os.Getenv("MYSQLDBNAME"))

	// Connect to MySQL database using GORM
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			os.Getenv("MYSQLUSER"),
			os.Getenv("MYSQLPASS"),
			os.Getenv("MYSQLDOMEN"),
			os.Getenv("MYSQLDBNAME")))
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
	} else {
		fmt.Println("Connected to db")
	}

	// Set max idle connection to 5
	db.DB().SetMaxIdleConns(5)

	// Set max open connection to 10
	db.DB().SetMaxOpenConns(10)

	// Count the number of records in the user table
	var count int64
	err = db.Model(&Files{}).Count(&count).Error
	if err != nil {
		fmt.Println("Error counting records: ", err)
	}

	// Return the number of records
	w.Write([]byte(fmt.Sprintf("Number of records in user table: %d\n", count)))
}

func main() {
	fmt.Println("DEMO SERVER STARTED!!!")

	http.HandleFunc("/demo", Democount)
	http.ListenAndServe(":8000", nil)
}
