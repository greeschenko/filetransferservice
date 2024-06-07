package main

//TODO
// + change all data to correct in db and reapload all files and repeat tests
// + copy file to new path with os
// + add threads
// - test without copeing on prod
// - change copy to move

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

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

func DoOsExec(com string, args ...string) {
	cmd := exec.Command(com, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}

func (u User) HandleAll() {
	grouplist := map[string]string{
		"CopyPasp":              u.CopyPasp,
		"CopyIpn":               u.CopyIpn,
		"CopyEdr":               u.CopyEdr,
		"CopyPdv":               u.CopyPdv,
		"CopyEp":                u.CopyEp,
		"CopyPdvCancell":        u.CopyPdvCancell,
		"RefBank":               u.RefBank,
		"CopyFound_doc":         u.CopyFound_doc,
		"CopyEdrpoy":            u.CopyEdrpoy,
		"CopyFoundFocRedaction": u.CopyFoundFocRedaction,
		"CopyHasLegal":          u.CopyHasLegal,
		"RefRegLegal":           u.RefRegLegal,
		"CopyRegInukraine":      u.CopyRegInukraine,
		"RefBankInukraine":      u.RefBankInukraine,
		"OtherFile":             u.OtherFile,
	}
	for e := range grouplist {
		if grouplist[e] != "" && grouplist[e][0] != '/' {
			var atts []Attachment
			var files []File
			err := DB.Where("`group` = ?", grouplist[e]).Find(&atts).Error
			if err != nil {
				fmt.Println("Error get list of attachments: ", err)
			}
			for i := range atts {
				err := DB.Where("id = ?", atts[i].FileID).Find(&files).Error
				if err != nil {
					fmt.Println("Error get list of files: ", err)
				}
			}
			if len(files) > 0 {
				fmt.Println("user handle start", u.ID, u.Email)
				fmt.Println(">>>", e, grouplist[e])
				for m := range files {
					newpath := strings.Replace(files[m].Path, "/uploads/", "/uploads2/", -1)
					//fmt.Println(">>>>>>", files[m].Path, files[m].Name, files[m].Ext)
					//fmt.Println("mkdir", "--parrents", "web/uploads2/"+grouplist[e])
					fmt.Println(">>>>>>", "mv", os.Getenv("POLONEXPUBPATH")+files[m].Path+files[m].Name+"."+files[m].Ext, os.Getenv("POLONEXPUBPATH")+newpath)

					//TODO turn on after testing
					//					DoOsExec("mkdir", "-p", os.Getenv("POLONEXPUBPATH")+newpath)
					//
					//					if files[m].Type == 1 {
					//						DoOsExec("mv", os.Getenv("POLONEXPUBPATH")+files[m].Path+files[m].Name+"_big_."+files[m].Ext, os.Getenv("POLONEXPUBPATH")+newpath)
					//						DoOsExec("mv", os.Getenv("POLONEXPUBPATH")+files[m].Path+files[m].Name+"_tumb_."+files[m].Ext, os.Getenv("POLONEXPUBPATH")+newpath)
					//					} else {
					//						DoOsExec("mv", os.Getenv("POLONEXPUBPATH")+files[m].Path+files[m].Name+"."+files[m].Ext, os.Getenv("POLONEXPUBPATH")+newpath)
					//					}
				}
			}
		}
	}

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
	Type   int    `json:"type"`
}

func Start(w http.ResponseWriter, r *http.Request) {
	var (
		data struct {
			Count int64
			Users []User
		}
	)

	fmt.Println(os.Getenv("MYSQLUSER"), os.Getenv("MYSQLPASS"), os.Getenv("MYSQLDOMEN"), os.Getenv("MYSQLDBNAME"))

	// Count the number of records in the user table
	err := DB.Model(&User{}).Count(&data.Count).Error
	if err != nil {
		fmt.Println("Error counting records: ", err)
	}

	err = DB.Where("id > 0").Find(&data.Users).Error
	if err != nil {
		fmt.Println("Error get list of users: ", err)
	}

	limiter := make(chan int, 2)
	for e := range data.Users {
		limiter <- 1
		go func(e int) {
			data.Users[e].HandleAll()
			<-limiter
		}(e)
	}

	res, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		fmt.Println("Error data marshal: ", err)
	}

	// Return the number of records
	w.Write(res)
}

func main() {
	var err error
	fmt.Println("DEMO SERVER STARTED!!!")

	// Connect to MySQL database using GORM
	DB, err = gorm.Open(
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
	DB.DB().SetMaxIdleConns(5)

	// Set max open connection to 10
	DB.DB().SetMaxOpenConns(10)

	http.HandleFunc("/start", Start)
	http.ListenAndServe(":8000", nil)
}
