package models

import "time"

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type BillSplit struct {
	ID int
	Name         string        `json:"name"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	Email_name string `json:"email_name"`
	User_id    int
}

type Expense struct {
	Expense      string
	Amount       float64
	Payer        string
	Participants []string
}

type ExpenseInfo struct {
	BillSplitID  int
	Name         string
	PayerName    string
	Amount       float64
	CreatedAt    time.Time
	Participants []string
}