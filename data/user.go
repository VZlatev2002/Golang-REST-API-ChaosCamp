package data

import (
	"fmt"
	"time"
)

// Participant struct has info of a Participant
type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	CreatedAt string
}


func Create(name string, email string, password string) (u User, err error) {
	fmt.Println("inside create")
	_ , err = Db.Exec("insert into user (name, email, passowrd, created_at) values (?, ?, ?, ?)", name, email, password, time.Now())
	
	if err != nil {
		return User{}, fmt.Errorf("Create user failed: %v", err)
	}
	err = Db.QueryRow("SELECT * FROM user WHERE name = ?", name).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	
	return

}

func CheckPassword(name, passoword string) (pass string, err error) {
	stmtOut, err := Db.Prepare("SELECT passowrd FROM user WHERE email = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmtOut.Close()


	err = stmtOut.QueryRow(name).Scan(&pass)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return 
}
