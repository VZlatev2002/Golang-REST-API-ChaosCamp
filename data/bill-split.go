package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/velizarzlatev/final_project/models"
)

type DBBillSplit struct {
	ID     int `sql:"id"`
	Name      string `sql:"name,varchar(64)"`
	CreatedAt time.Time `sql:"created_at,timestamp"`
	Participants []models.Participant 

}



func (b DBBillSplit) ToModelBill() models.BillSplit {
	return models.BillSplit{
		ID: b.ID,
		Name: b.Name,
		Participants: b.Participants,
	}
}

func fromModelBill(bs models.BillSplit) DBBillSplit {
	return DBBillSplit{
		Name: bs.Name,
		Participants: bs.Participants,
	}
}

func (db *BillSplitterDB) CreateBillSplit(bs models.BillSplit) (billsplit models.BillSplit, err error) {
	dbBill := fromModelBill(bs)	
	_, err = db.conn.Exec("insert into billsplit (name, created_at) values (?, ?)", bs.Name, time.Now())
	if err != nil {
		return bs, fmt.Errorf("Create billSplit failed: %v", err)
	}
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = db.conn.QueryRow("SELECT * FROM billsplit WHERE name = ?", bs.Name).Scan(&dbBill.ID, &dbBill.Name, &dbBill.CreatedAt)
	

	return dbBill.ToModelBill(), err
}

func (db *BillSplitterDB) CreateParticipants(ctx context.Context, billSplit models.BillSplit) (err error) {	
	// stms, err := db.conn.PrepareContext(ctx, "select * from users where user_id = ?", billSplit.)


	for _, participant := range billSplit.Participants {
		
		row := db.conn.QueryRowContext(ctx, "select id from user where email = ?", participant.Email_name)
		if err != nil {
			log.Println("error duing selection from user", err)
			return err
		}
		err = row.Scan(&participant.User_id)
		
		log.Println(participant)
		if participant.User_id != 0 {
			_, err := db.conn.ExecContext(ctx, "insert into participant(email_name, billsplit_id, user_id, created_at) values (?, ?, ?, ?)", &participant.Email_name, billSplit.ID, &participant.User_id, time.Now())
			if err != nil {
				log.Println("error duing insertion to participant with know user_id", err)
				return err
			}
		} else {
			_, err = db.conn.ExecContext(ctx, "insert into participant(email_name, billsplit_id, created_at) values (?, ?, ?)", &participant.Email_name, billSplit.ID, time.Now())
			if err != nil {
				log.Println("error duing insertion to participant with not known user_id", err)
				return err
			}
		}
		

	}
	return  err

}

func (db *BillSplitterDB) BillSplitByName(name string) (billsplit DBBillSplit, err error) {
	row := db.conn.QueryRow("SELECT * FROM billsplit WHERE name = ?", name)
	row.Scan(&billsplit.ID, &billsplit.Name, &billsplit.CreatedAt)
	log.Println(billsplit)
	return
}

func (db *BillSplitterDB) BillSplitByID(id int) (billsplit DBBillSplit, err error) {
	log.Println("BillSplitByID", id)
	row := db.conn.QueryRow("SELECT id,name, created_at FROM billsplit WHERE id = ?", id)
	row.Scan(&billsplit.ID, &billsplit.Name, &billsplit.CreatedAt)
	return
}

func (db *BillSplitterDB) BillSplits() (billSplits []DBBillSplit, err error) {
	//defer db.Close()
	log.Println("BillSplits form data")
	rows, err := db.conn.Query("SELECT * FROM billsplit ORDER BY created_at DESC")
	if err != nil {
		log.Println("Can't query rows from database", err)
		return
	}
	for rows.Next() {
		conv := DBBillSplit{}
		if err = rows.Scan(&conv.ID, &conv.Name, &conv.CreatedAt); err != nil {
			log.Println("Error scanning row to struct", err)
			return
		}
		billSplits = append(billSplits, conv)
	}
	rows.Close()
	return
}

func (db *BillSplitterDB) GetFullBalance(billSplitID int) (fullBalance map[string]float64, err error) {
	log.Println("GetFullBalance", billSplitID)
	expenses, err := db.Expenses(billSplitID)
	if err != nil {
		log.Fatal(err)
	}
	fullBalance = make(map[string]float64)
	participants, err := db.Participants(billSplitID)
	for _, participant := range participants {
		fullBalance[participant.Name] = 0
	}
	for _, expense := range expenses {
		balanceExpense := db.Balance(expense)
		for k, v := range balanceExpense {
			fullBalance[k] += v
		}
	}
	if err != nil {
		return
	}
	return
}
func (db *BillSplitterDB) Participants(billSplitID int) (items []ParticipantDB, err error) {
	log.Println("Participants of BillSplit:", billSplitID)
	rows, err := db.conn.Query("SELECT id, email_name, billsplit_id, user_id FROM participant where billsplit_id = ? ORDER BY created_at DESC", billSplitID)
	if err != nil {
		log.Println("Error getting participants", err)
		return
	}
	for rows.Next() {
		post := ParticipantDB{}
		if err = rows.Scan(&post.ID, &post.Name, &post.BillSplitID, &post.UserID); err != nil {
			log.Println("error getting participants", err)
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return
}

func (db *BillSplitterDB) BillSplitDeleteAll() (err error) {
	statement := "delete from billsplit"
	_, err = db.conn.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}