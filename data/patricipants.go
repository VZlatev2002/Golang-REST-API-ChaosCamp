package data

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type ParticipantDB struct {
	ID int `sql:"id"`
	Name string `sql:"email_name"`
	BillSplitID int `sql:"billsplit_id"`
	UserID int `sql:"user_id"`
	CreatedAt time.Time `sql:"created_at"`
}

func (db *BillSplitterDB) ParticipantsByName(names []string, billsplit_id int) (items []ParticipantDB, err error) {
	//defer db.Close()
	log.Println("ParticipantsByName")
	sqlStr := "SELECT * FROM participant where billsplit_id = ? and email_name in (?" + strings.Repeat(",?", len(names)-1) + ") ORDER BY created_at DESC"

	args := make([]interface{}, len(names)+1)
	args[0] = billsplit_id
	for i, id := range names {
		args[i+1] = id
	}
	rows, err := db.conn.Query(sqlStr, args...)
	// (?` + strings.Repeat(",?", len(args)-1) + `)`

	if err != nil {
		return
	}
	for rows.Next() {
		post := ParticipantDB{}
		if err = rows.Scan(&post.ID, &post.Name, &post.BillSplitID,  &post.UserID, &post.CreatedAt); err != nil {
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return
}

func (db *BillSplitterDB) ParticipantByName(name string, billSplit_id int) (participant ParticipantDB, err error){
	
	log.Println("ParticipantByName", name, billSplit_id)

	row := db.conn.QueryRow("SELECT * FROM participant WHERE email_name = ? and billSplit_id = ?", name, billSplit_id)
	err = row.Scan(&participant.ID, &participant.Name, &participant.BillSplitID, &participant.UserID, &participant.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no rows foundd, no participants for that billsplit_id")
			return
		} else {
			log.Println("Error querying database:", err)
			return
		}
	}
	return 
	
}

func (db *BillSplitterDB) ParticipantExpenseDeleteAll() (err error) {
	statement := "delete from participant_expense"
	_, err = db.conn.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (db *BillSplitterDB) ParticipantDeleteAll() (err error) {
	//defer db.Close()
	statement := "delete  from participant"
	_, err = db.conn.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}