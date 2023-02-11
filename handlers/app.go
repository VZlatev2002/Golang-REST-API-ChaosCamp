// Package classification of BilSplittter API
//
// Documentation for BillSplitter API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- applications/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/velizarzlatev/final_project/data"
	"github.com/velizarzlatev/final_project/helpers"
	"github.com/velizarzlatev/final_project/models"

	"golang.org/x/crypto/bcrypt"
)

type Expense struct {
	Id          int
	Uuid        string
	Name        string
	Amount      float64
	BillSplitID int
	PayerName   string
	CreatedAt   time.Time
}

type Database interface {
	Create(context.Context, models.User) (models.User, error)

	GetPassword(context.Context, string) (string, error)

	CreateBillSplit(bs models.BillSplit) (billsplit models.BillSplit, err error)

	CreateParticipants(ctx context.Context, billSplit models.BillSplit) (err error)

	ParticipantByName(name string, billSplit_id int) (participant data.ParticipantDB, err error)

	CreateExpense(name string, amount float64, payer string, id int) (expense data.ExpenseDB, err error)

	BillSplitByID(id int) (billsplit data.DBBillSplit, err error)

	AddParticipants(names []string, billsplit_id int, expense_id int) (err error)

	ExpenseByID(expense_id int) (expense data.ExpenseDB, err error)

	ExpenseParticipants(expenseId int) (items []string, err error)

	BillSplits() (billSplits []data.DBBillSplit, err error)

	BillSplitByName(name string) (billsplit data.DBBillSplit, err error)

	Expenses(billSplitID int) (items []data.ExpenseDB, err error)

	GetFullBalance(billSplitID int) (fullBalance map[string]float64, err error)
	
}

type App struct {
	Router *mux.Router
	db    Database
}

func NewApp(db Database, r *mux.Router) *App {
	return &App{Router: r, db: db}
}

func Initialize() (a *App) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := data.InitDb()
	Router := mux.NewRouter()
	return NewApp(db, Router)
}

// swagger:route POST /signup signup SignupUser
// Signups a user, checks if the user exists and return a cookie with jwt token

func (a *App) Signup(writer http.ResponseWriter, request *http.Request) {
	log.Println("Signup")

	ctx := context.Background()
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	u := models.User{}
	err := json.NewDecoder(request.Body).Decode(&u)
	if err != nil {
		log.Print(err)
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(hash)
	//Create the user

	result, err := a.db.Create(ctx, u)

	if err != nil {
		log.Print(err)
	}
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot create new User")
	} else {
		helpers.RespondWithJSON(writer, http.StatusCreated, result)
	}

}
func (a *App) Error(writer http.ResponseWriter, request *http.Request){
	helpers.RespondWithJSON(writer, http.StatusNotExtended, "error: USE POST")
}




func (a *App) GetAllUsers(writer http.ResponseWriter, request *http.Request){
	writer.Write([]byte("Inside the Api Authenticated path"))
	fmt.Println("inside the api auth protected path")
}

// NewBillSplit creates a new billSplit in the database
func (a *App) NewBillSplit(writer http.ResponseWriter, request *http.Request){
	log.Println("NewBillSplit")
	ctx := context.Background()
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	bs := models.BillSplit{} //omit billsplitId here

	err := json.NewDecoder(request.Body).Decode(&bs)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot create new BillSplit")
	}
	billSplitToParticipants, err := a.db.CreateBillSplit(bs)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	}
	log.Println("Created new BillSplit", billSplitToParticipants)
	err = a.db.CreateParticipants(ctx, billSplitToParticipants)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot create new BillSplit")
	} else {
		helpers.RespondWithJSON(writer, http.StatusCreated, bs)
	}
}

// GetExpense get the expense in the database given its id
func (a *App) GetExpense(writer http.ResponseWriter, request *http.Request) {
	log.Println("GetExpense")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	expenseID := mux.Vars(request)["ExpenseId"]

	intValue, err := strconv.ParseInt(expenseID, 10, 0)
	log.Println(int(intValue))
	if err != nil {
		log.Println("mistake in parsing billsplitID to intValue")
		return
	}
	expense, err := a.db.ExpenseByID(int(intValue))
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	}
	participants, err := a.db.ExpenseParticipants(expense.Id)

	expenseInfo := models.ExpenseInfo{
		BillSplitID:  expense.BillSplitID,
		Name:         expense.Name,
		PayerName:    expense.PayerName,
		Amount:       expense.Amount,
		CreatedAt:    expense.CreatedAt,
		Participants: participants,
	}
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	} else {
		//generateHTML(writer, surveys, "layout","index")
		helpers.RespondWithJSON(writer, 200, expenseInfo)
	}
}

// NewExpense creates a new expense record in the database
func (a *App) NewExpense(writer http.ResponseWriter, request *http.Request) {
	log.Println("NewExpense")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var body struct {
		Expense      string
		Amount       float64
		Payer        string
		Participants []string
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	}
	billSplitName := mux.Vars(request)["BillSplitName"]
	billSplit, err := a.db.BillSplitByName(billSplitName)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	}
	expense, err := a.db.CreateExpense(body.Expense, body.Amount, body.Payer, billSplit.ID)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Can't create expense")
	}
	err = a.db.AddParticipants(body.Participants, expense.BillSplitID, expense.Id)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get threads")
	} else {
		helpers.RespondWithJSON(writer, http.StatusCreated, body)
	}
}

	// GetBillSplits gets a billsplits in the database
func (a *App) GetBillSplits(writer http.ResponseWriter, request *http.Request) {
	log.Println("GetBillSplits")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	surveys, err := a.db.BillSplits()
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get BillSplits")
	} else {
		//generateHTML(writer, surveys, "layout", "public.navbar", "index")
		helpers.RespondWithJSON(writer, 200, surveys)
	}
}

func (a *App) GetBillSplitById(writer http.ResponseWriter, request *http.Request) {
	log.Println("GetBillSplitByID")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	billSplitId := mux.Vars(request)["BillSplitId"]
	intValue, err := strconv.ParseInt(billSplitId, 10, 0)
	log.Println(int(intValue))
	if err != nil {
		log.Println("mistake in parsing billsplitID to intValue")
		return
	}
	billSplit, err := a.db.BillSplitByID(int(intValue))
	if err != nil {
		helpers.ErrorMessage(writer, request, "Cannot get billSplits")
	} else{
		helpers.RespondWithJSON(writer, http.StatusOK, billSplit)
	}
}

func (a *App) Login(writer http.ResponseWriter, request *http.Request){
	log.Println("Login")

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		log.Println(err)
	}
	
	hash, err := a.db.GetPassword(request.Context(), body.Email)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password))
	if err != nil {
		writer.Write([]byte("IInvalid email or password"))
		log.Println(err)
	}
	
	tokenString, err := CreateJWT(body.Email)
	cookie := &http.Cookie{Name: "token",Value:tokenString}
	http.SetCookie(writer, cookie)

}

type CustomClaims struct {
	Username string 
	jwt.StandardClaims
}

func CreateJWT(username string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer: "https://localhost",
		},
	})

	tkn, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println("error from Creatr JWT", err.Error(), tkn)
		return "", err
	}
	return tkn, err
}

func ValidateToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err:= r.Cookie("token")

		
		if err!= nil{
			w.Write([]byte("not authorized due to no token in the header"))
			w.WriteHeader(http.StatusUnauthorized)

		}

		token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				w.Write([]byte("not authorized"))
				w.WriteHeader(http.StatusUnauthorized)
			}
			// checks for the method signature, maybe different from HMAC, ex: RSA(public/private key)
			return []byte(os.Getenv("SECRET_KEY")), nil
			})
		tkn, ok := token.Claims.(*CustomClaims)
		if !ok {
			return
		}
		if err := tkn.Valid(); err != nil {
			return
		}
		next(w,r)


		// if claims, ok := token.Claims.(jwt.MapClaims); ok{
		// 	if float64(time.Now().Unix()) > claims["exp"].(float64){
		// 		w.Write([]byte("not authorized bcause token has expired"))
		// 		w.WriteHeader(http.StatusUnauthorized)
			
		// 	}else {
		// 		next(w, r)
		// 	}
		// } else {
		// 	w.Write([]byte("not authorized because no time claims in the token"))
		// 	w.WriteHeader(http.StatusUnauthorized)
		// }
			
	})
}

func (a *App) GetParticipantsBalance(writer http.ResponseWriter, request *http.Request) {
	log.Println("GetParticipantsBalance")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

	billSplitId := mux.Vars(request)["BillSplitId"]
	intValue, err := strconv.ParseInt(billSplitId, 10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	billSplit, err := a.db.BillSplitByID(int(intValue)) // Query the database to check if there is such a BillSplitId
	if err != nil {
		helpers.ErrorMessage(writer, request, "No such BillSplitId")
	}

	balance, err := a.db.GetFullBalance(billSplit.ID)
	log.Println(err)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Error getting full balance")
	} else {
		log.Println("Responding to client")
		helpers.RespondWithJSON(writer, http.StatusCreated, balance)
	}
}



func (a *App) SetRoutes() {


	a.Router.HandleFunc("/signup", a.Signup).Methods("POST") // Validated but TODO: claims are removed after short period of time for some reason
	a.Router.HandleFunc("/error", a.Error)
	a.Router.HandleFunc("/login", a.Login).Methods("POST") //Validated

	a.Router.Handle("/api/v1/", ValidateToken(a.GetBillSplits)).Methods("GET") // Validate but returns null Participants field, TODO: refactor this
	a.Router.Handle("/api/v1/users", ValidateToken(a.GetAllUsers)).Methods("GET")
	a.Router.Handle("/api/v1/billsplit/new", ValidateToken(a.NewBillSplit)).Methods("POST") // Validated
	a.Router.Handle("/api/v1/billsplit/{BillSplitId}", ValidateToken(a.GetBillSplitById)).Methods("GET") //Validated
	// a.Router.HandleFunc("api/v1/billsplit/{BillSplitId}/expenses", a.GetBillSplitExpenses).Methods("GET")
	
	a.Router.Handle("/api/v1/billsplit/{BillSplitName}/expenses/new", ValidateToken(a.NewExpense)).Methods("POST") // Validated
	a.Router.Handle("/api/v1/expense/{ExpenseId}", ValidateToken(a.GetExpense)).Methods("GET") //Validated 

	// a.Router.HandleFunc("/billsplit/{BillSplitId}/participants", a.GetBillSplitParticipants).Methods("GET")
	// a.Router.HandleFunc("/billsplit/{BillSplitId}/participants/new", a.NewParticipants).Methods("POST")
	a.Router.Handle("/billsplit/{BillSplitId}/balance", ValidateToken(a.GetParticipantsBalance)).Methods("GET")

}
	

