package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


type BillSplitterDB struct{
	conn *sql.DB
}

func NewSQLBillSplitterDB(db *sql.DB) *BillSplitterDB {
	return &BillSplitterDB{
		conn: db,
	}
}

func InitDb() *BillSplitterDB{

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:"+os.Getenv("DBPORT"),
        DBName: os.Getenv("DBNAME"),
		AllowNativePasswords: true,
		ParseTime: true, // Convert to time.Time as in the structs
    }


	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	
    fmt.Println("Connected!")
	return NewSQLBillSplitterDB(db)
}

// Use below for testing purposes
var (
	// DB_USER  database username
	DB_USER = GetEnv("DBUSER", "root")
	// DB_PASSWORD username password
	DB_PASSWORD = GetEnv("DBPASS", "root")
	// DB_NAME database name
	DB_TESTNAME = GetEnv("DBTESTNAME", "v1_bsplitter_test")
	// DB_PORT database connection port
	DB_PORT = GetEnv("DBPORT", "3306")
)
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		fmt.Println(key, fallback)
		return fallback
	}
	fmt.Println(key, value)
	return value
}

// Initialize test database connection with the same configuration as normel, just different
func InitTestDb() *BillSplitterDB{

	
	cfg := mysql.Config{
        User:   DB_USER,
        Passwd: DB_PASSWORD,
        Net:    "tcp",
        Addr:   "127.0.0.1:"+DB_PORT,
        DBName: DB_TESTNAME,
		AllowNativePasswords: true,
		ParseTime: true, // Convert to time.Time as in the structs
    }


	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	
    fmt.Println("Connected!")
	return NewSQLBillSplitterDB(db)
}

// SetupDB clears the database completely
func (db *BillSplitterDB) SetupDB() {
	err := db.ParticipantExpenseDeleteAll()
	if err != nil {
		log.Fatal(err)
	}
	err = db.ExpenseDeleteAll()
	if err != nil {
		log.Fatal(err)
	}
	err = db.ParticipantDeleteAll()
	if err != nil {
		log.Fatal(err)
	}
	err = db.BillSplitDeleteAll()
	if err != nil {
		log.Fatal(err)
	
	}
}

