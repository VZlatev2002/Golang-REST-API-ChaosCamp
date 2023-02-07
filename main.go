package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/velizarzlatev/final_project/data"
	"github.com/velizarzlatev/final_project/helpers"
	"golang.org/x/crypto/bcrypt"
)

type Database interface {
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp(db *sql.DB, r *mux.Router) *App {
	return &App{Router: r, DB: db}
}

func (a *App) Signup(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("inside signup")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	//Create the user

	result, err := data.Create(user.Name, user.Email, string(hash))

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

func (a *App) Login(writer http.ResponseWriter, request *http.Request){
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		log.Print(err)
	}
	
	hash, err := data.CheckPassword(body.Email, body.Password)
	if err != nil {
		helpers.ErrorMessage(writer, request, "Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password))
	if err != nil {
		helpers.ErrorMessage(writer, request, "Invalid email or password")
	}
	
	tokenString, err := CreateJWT()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := &http.Cookie{Name: "token",Value:tokenString,Expires:expiration}
	http.SetCookie(writer, cookie)

}


func (a *App) GetAllUsers(writer http.ResponseWriter, request *http.Request){
	writer.Write([]byte("Inside the Api Authenticated path"))
	fmt.Println("inside the api auth protected path")
}

 



func CreateJWT() (string, error){
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"]= time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		fmt.Println("error from Creatr JWT", err.Error(), tokenStr)
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err:= r.Cookie("token")

		

		if err!= nil{
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized due to no token in the header"))
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized"))
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
			})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
			if float64(time.Now().Unix()) > claims["exp"].(float64){
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized bcause token has expired"))
			}
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized because no time claims in the token"))
		}
			
	})
}


func main() {
	l := log.New(os.Stdout, "bill_splitter", log.LstdFlags)

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	data.InitDb()
	Router := mux.NewRouter()
	a := NewApp(data.Db, Router)

	a.Router.HandleFunc("/signup", a.Signup).Methods("POST")
	a.Router.HandleFunc("/error", a.Error)
	a.Router.Handle("/api/v1/users", ValidateToken(a.GetAllUsers)).Methods("GET")


	a.Router.HandleFunc("/login", a.Login).Methods("POST")
	// a.Router.HandleFunc("/signup", a.Signup).Methods("POST")

	s := &http.Server{
		Addr:         ":3000",
		Handler:      Router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
