package main

//TODO
// - get user list
// - check user list docs records if not empty handle it
// - doc record code handle:
//   - find attachment
//   - find file and gen path
// - try connect volume with files
// - copy file to new path with os
// - change record in file table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User is a struct that represents a user in the database
type User struct {
	ID                    int    `json:"id"`
	Email                 string `json:"email"`
	Status                int    `json:"status"`
	CopyPasp              string `json:"copy_pasp"`
	CopyIpn               string `json:"copy_ipn"`
	CopyEdr               string `json:"copy_edr"`
	CopyPdv               string `json:"copy_pdv"`
	CopyEp                string `json:"copy_ep"`
	CopyPdvCancell        string `json:"copy_pdv_cancell"`
	RefBank               string `json:"ref_bank"`
	CopyFound_doc         string `json:"copy_found_doc"`
	CopyEdrpoy            string `json:"copy_edrpoy"`
	CopyFoundFocRedaction string `json:"copy_found_doc_redaction"`
	CopyHasLegal          string `json:"copy_has_legal"`
	RefRegLegal           string `json:"ref_reg_legal"`
	CopyRegInukraine      string `json:"copy_reg_inukraine"`
	RefBankInukraine      string `json:"ref_bank_inukraine"`
	OtherFile             string `json:"other_file"`
}

func (User) TableName() string {
	return "user"
}

type Attachment struct {
	ID     int    `json:"id"`
	Group  string `json:"group"`
	FileID int    `json:"file_id"`
}

type File struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Ext    string `json:"ext"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

func Democount(w http.ResponseWriter, r *http.Request) {
	var (
		data  struct {
			Count int64
			Users []User
            Attachments []Attachment
            Files []File
		}
	)

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
	err = db.Model(&User{}).Count(&data.Count).Error
	if err != nil {
		fmt.Println("Error counting records: ", err)
	}

	err = db.Where("id > 0").Limit(10).Find(&data.Users).Error
	if err != nil {
		fmt.Println("Error get list of users: ", err)
	}

	err = db.Where("id > 0").Find(&data.Attachments).Error
	if err != nil {
		fmt.Println("Error get list of attachments: ", err)
	}

	err = db.Where("id > 0").Find(&data.Files).Error
	if err != nil {
		fmt.Println("Error get list of files: ", err)
	}

	res, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		fmt.Println("Error data marshal: ", err)
	}

	// Return the number of records
	w.Write(res)
}

func main() {
	fmt.Println("DEMO SERVER STARTED!!!")

	http.HandleFunc("/demo", Democount)
	http.ListenAndServe(":8000", nil)
}
