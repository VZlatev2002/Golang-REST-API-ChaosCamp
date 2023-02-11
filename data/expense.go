package data

import (
	"log"
	"strings"
	"time"
)

type ExpenseDB struct {
	Id          int
	Name        string
	Amount      float64
	BillSplitID int
	PayerName   string
	CreatedAt   time.Time
}

func (db *BillSplitterDB) AddParticipants(names []string, billsplit_id int, expense_id int) (err error) {
	log.Println("AddParticipants")
	participants, err := db.ParticipantsByName(names, billsplit_id)
	if err != nil {
		log.Fatal(err)
	}
	sqlStr := "insert into participant_expense(participant_id, expense_id) VALUES "
	vals := []interface{}{}

	for _, row := range participants {
		sqlStr += "(?, ?),"
		vals = append(vals, row.ID, expense_id)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	
	//prepare the statement
	stmt, _ := db.conn.Prepare(sqlStr)

	
	//format all vals at once
	_, err = stmt.Exec(vals...)
	return
}

func (db *BillSplitterDB) CreateExpense(name string, amount float64, payer string, billsplitid int) (expense ExpenseDB, err error){
	//defer db.Close()
	participant, err := db.ParticipantByName(payer, billsplitid)
	if err != nil {
		log.Println("Payer is not registered participant in a billsplit", err)
		return
	}
	_, err = db.conn.Exec("insert into expense (name, amount, billsplit_id, participant_id, created_at) values (?,?,?,?,?)", name, amount, billsplitid, participant.ID, time.Now())
	statement := "SELECT e.id, e.name, e.amount, e.billsplit_id, p.email_name, e.created_at FROM expense AS e INNER JOIN participant AS p ON e.participant_id = p.id where e.name = ? and e.billsplit_id = ?"
	if err != nil {
		return
	}
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = db.conn.QueryRow(statement, name, billsplitid).Scan(&expense.Id, &expense.Name, &expense.Amount, &expense.BillSplitID, &expense.PayerName, &expense.CreatedAt)
	if err != nil {
		log.Println("Problem with the Join Query in CreateExpense", err)

		return
	}
	log.Println(expense)
	return
}

func (db *BillSplitterDB) ExpenseByID(expense_id int) (expense ExpenseDB, err error) {
	log.Println("ExpenseByID", expense_id)
	err = db.conn.QueryRow("SELECT e.id, e.name, e.amount, e.billsplit_id, p.email_name, e.created_at FROM expense as e INNER JOIN participant as p ON e.participant_id = p.id where e.id = ?", expense_id).
		Scan(&expense.Id, &expense.Name, &expense.Amount, &expense.BillSplitID, &expense.PayerName, &expense.CreatedAt)
	return
}

func (db *BillSplitterDB) Expenses(billSplitID int) (items []ExpenseDB, err error) {
	//defer db.Close()s
	log.Println("Expenses of BillSplit:", billSplitID)
	rows, err := db.conn.Query("SELECT e.id, e.name, e.amount, e.billsplit_id, p.email_name, e.created_at FROM expense as e INNER JOIN participant as p ON e.participant_id = p.id where e.billSplit_id = ? ORDER BY created_at DESC", billSplitID)
	if err != nil {
		return
	}
	for rows.Next() {
		post := ExpenseDB{}
		if err = rows.Scan(&post.Id, &post.Name, &post.Amount, &post.BillSplitID, &post.PayerName, &post.CreatedAt); err != nil {
			return
		}
		items = append(items, post)
	}
	rows.Close()
	return

}
func (db *BillSplitterDB) Balance(expense ExpenseDB) map[string]float64 {
	log.Println("Balance of", expense.Name)
	participants, err := db.ExpenseParticipants(expense.Id)
	if err != nil {
		log.Println("Error getting participants in the expense", err)
		log.Fatal(err)
	}
	payer, err := db.ParticipantByName(expense.PayerName, expense.BillSplitID)
	if err != nil {
		log.Println("Error getting the Participant by name", err)
		log.Fatal(err)
	}
	balance := make(map[string]float64)
	balance[payer.Name] = expense.Amount
	for _, participant := range participants {
		balance[participant] += -expense.Amount / float64(len(participants))
	}
	return balance
}

func (db *BillSplitterDB) ExpenseParticipants(expenseId int) (items []string, err error) {
	//defer db.Close()
	log.Println("Participants of expense", expenseId)
	rows, err := db.conn.Query("SELECT p.email_name FROM participant_expense as pe INNER JOIN participant p ON p.id = pe.participant_id WHERE pe.expense_id = ? ORDER BY p.created_at DESC", expenseId)
	if err != nil {
		log.Println("Error getting participants of expense", err)
		return
	}
	for rows.Next() {
		var participant string
		if err = rows.Scan(&participant); err != nil {
			return
		}
		items = append(items, participant)
	}
	rows.Close()
	return
}
func (db *BillSplitterDB) ExpenseDeleteAll() (err error) {
	statement := "delete from expense"
	_, err = db.conn.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
	return
}