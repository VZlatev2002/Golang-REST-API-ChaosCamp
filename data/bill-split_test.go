package data

import (
	"log"
	"reflect"
	"testing"

	"github.com/velizarzlatev/bill-splitter/models"
)

// ID           int
// Name         string        `json:"name"`
// Participants []Participant `json:"participants"`

// type BillSplit struct {
// 	ID int
// 	Name         string        `json:"name"`
// 	Participants []Participant `json:"participants"`
// }

// type Participant struct {
// 	Email_name string `json:"email_name"`
// 	User_id    int
// }

func TestCreateBillSplit(t *testing.T) {
	db := InitTestDb()
	db.SetupDB()
	tests := []struct {
		name     string
		wantBill models.BillSplit
		wantErr  bool
	}{
		{"test0", 
			models.BillSplit{
				Name: "Holidays in Spain",
				Participants: []models.Participant{
					{Email_name: " Harry Potter"},
					{Email_name: "ParticipantWithouUsername"},
				},
			}, false,
		},
		{"test1", models.BillSplit{
			Name: "Holidays in Spain",
			Participants: []models.Participant{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBillSplit, err := db.CreateBillSplit(tt.wantBill)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBillSplit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBillSplit.Name != tt.wantBill.Name { // Use name, beacuse slices a can't be compared and the id is not inputed by the user
				t.Errorf("CreateBillSplit() gotSurvey = %v, want %v", gotBillSplit.Name, tt.wantBill.Name)
			}
		})
	}
	db.conn.Close()
}

func TestCreateParticipants(t *testing.T) {
	db := InitTestDb()
	db.SetupDB()
	tests := []struct {
		name     string
		wantBill models.BillSplit
		wantErr  bool
	}{
		{"test0", 
			models.BillSplit{
				Name: "Holidays in Spain",
				Participants: []models.Participant{
					{Email_name: " Harry Potter"},
					{Email_name: "ParticipantWithouUsername"},
				},
			}, false,
		},
		{"test1", models.BillSplit{
			Name: "Holidays in Spain",
			Participants: []models.Participant{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBillSplit, err := db.CreateBillSplit(tt.wantBill)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBillSplit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBillSplit.Name != tt.wantBill.Name { // Use name, beacuse slices a can't be compared and the id is not inputed by the user
				t.Errorf("CreateBillSplit() gotSurvey = %v, want %v", gotBillSplit.Name, tt.wantBill.Name)
			}
		})
	}
	db.conn.Close()
}

func TestBillSplit_GetFullBalance(t *testing.T) {
	db := InitTestDb()
	db.SetupDB()
	testbillSplit := models.BillSplit{Name: "Trip to Azkaban",
				Participants: []models.Participant{
					{Email_name: "Harry Potter"},
					{Email_name: "Prof Albus"},
					{Email_name: "Ginny Weasley"},
					{Email_name: "Lord Voldemort"},

				}}
	tests := []struct {
		name        string
		wantBalance map[string]float64
		wantErr     bool
	}{
		{
			"testGetFull",
			map[string]float64{
				"Harry Potter": 22.5,
				"Prof Albus": -12.5,
				"Ginny Weasley": -12.5,
				"Lord Voldemort": 2.5,
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billSplit, err := db.CreateBillSplit(testbillSplit)
			if err != nil {
				log.Fatal(err)
			}
			billSplit, err = db.CreateParticipants(billSplit)
			if err != nil {
				log.Fatal(err)
			}
			expense1, err := db.CreateExpense("expense1", 50.0, "Harry Potter", billSplit.ID)
			if err != nil {
				log.Fatal(err)
			}
			err = db.AddParticipants([]string{"Harry Potter", "Lord Voldemort", "Prof Albus", "Ginny Weasley"}, billSplit.ID, expense1.Id)
			if err != nil {
				log.Fatal(err)
			}
			expense2, err := db.CreateExpense("expense2", 30.0, "Lord Voldemort", billSplit.ID)
			if err != nil {
				log.Fatal(err)
			}
			err = db.AddParticipants([]string{"Harry Potter", "Lord Voldemort"}, billSplit.ID, expense2.Id)
			if err != nil {
				log.Fatal(err)
			}
			gotBalance, err := db.GetFullBalance(billSplit.ID)
			if err != nil {
				log.Fatal(err)
			}

			if !reflect.DeepEqual(gotBalance, tt.wantBalance) {
				t.Errorf("Balance() gotBalance = %v, want %v", gotBalance, tt.wantBalance)
			}

		})
	}
	err := db.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}