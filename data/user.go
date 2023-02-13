package data

import (
	"context"
	"fmt"
	"time"

	"github.com/velizarzlatev/bill-splitter/models"
)

type DBUser struct {
	ID int `sql:"id"`
	Name string `sql:"name,varchar(64)"`
	Email string `sql:"email,varchar(102)"`
	Password string `sql:"passowrd,varchar(255)"`
	CreatedAt time.Time `sql:"created_at,timestamp"`
}



func (b DBUser) ToModelUser() models.User {
	return models.User{
		Name: b.Name,
		Email: b.Email,
		Password: b.Password,
	}
}

func fromModelUser(user models.User) DBUser {
	return DBUser{
		Name: user.Name,
		Email: user.Email,
		Password: user.Password,
	}
}



func (db *BillSplitterDB) Create(context.Context, models.User) (u models.User, err error) {
	dbUser := fromModelUser(u)
	_ , err = db.conn.Exec("insert into user (name, email, password, created_at) values (?, ?, ?, ?)", dbUser.Name, dbUser.Email, dbUser.Password, time.Now())
	
	if err != nil {
		return models.User{}, fmt.Errorf("Create user failed: %v", err)
	}
	err = db.conn.QueryRow("SELECT * FROM user WHERE name = ?", dbUser.Name).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password, &dbUser.CreatedAt)
	
	return	dbUser.ToModelUser(), nil

}

func (db *BillSplitterDB) GetPassword(ctx context.Context, email string) (hash string, err error) {
	stmtOut, err := db.conn.Prepare("SELECT password FROM user WHERE email = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(email).Scan(&hash)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}