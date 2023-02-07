package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


var Db *sql.DB

func InitDb(){

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "v1_bsplitter",
		AllowNativePasswords: true,
    }


	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}
	
    fmt.Println("Connected!")
}



// // createUUID creates a random UUID with from RFC 4122
// // adapted from http://github.com/nu7hatch/gouuid
// func createUUID() (uuid string) {
// 	u := new([16]byte)
// 	_, err := rand.Read(u[:])
// 	if err != nil {
// 		log.Fatalln("Cannot generate UUID", err)
// 	}

// 	// 0x40 is reserved variant from RFC 4122
// 	u[8] = (u[8] | 0x40) & 0x7F
// 	// Set the four most significant bits (bits 12 through 15) of the
// 	// time_hi_and_version field to the 4-bit version number.
// 	u[6] = (u[6] & 0xF) | (0x4 << 4)
// 	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
// 	return
// }

// func CreateBillSplit(name string) (billsplit BillSplit, err error) {
// 	//defer db.Close()
// 	statement := "insert into billsplit (uuid, name, created_at) values ($1, $2, $3) returning id, uuid, name, created_at"
// 	stmt, err := Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	// use QueryRow to return a row and scan the returned id into the Session struct
// 	err = stmt.QueryRow(createUUID(), name, time.Now()).Scan(&billsplit.Id, &billsplit.Uuid, &billsplit.Name, &billsplit.CreatedAt)
// 	if err != nil {
// 		return
// 	}
// 	err = stmt.Close()
// 	if err != nil {
// 		return
// 	}
// 	return
// }
// func BillSplits() (billSplits []BillSplit, err error) {
// 	//defer db.Close()
// 	rows, err := Db.Query("SELECT id, uuid, name, created_at FROM billsplit ORDER BY created_at DESC")
// 	if err != nil {
// 		return
// 	}
// 	for rows.Next() {
// 		conv := BillSplit{}
// 		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Name, &conv.CreatedAt); err != nil {
// 			return
// 		}
// 		billSplits = append(billSplits, conv)
// 	}
// 	rows.Close()
// 	return


// func CreateBillSplit(name string) (billsplit BillSplit, err error) {
// 	//defer db.Close()
// 	statement := "insert into billsplit (uuid, name, created_at) values ($1, $2, $3) returning id, uuid, name, created_at"
// 	stmt, err := Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	// use QueryRow to return a row and scan the returned id into the Session struct
// 	err = stmt.QueryRow(createUUID(), name, time.Now()).Scan(&billsplit.Id, &billsplit.Uuid, &billsplit.Name, &billsplit.CreatedAt)
// 	if err != nil {
// 		return
// 	}
// 	err = stmt.Close()
// 	if err != nil {
// 		return
// 	}
// 	return
// }