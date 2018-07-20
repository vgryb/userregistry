package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var userTokenRoleMap = make(map[string]string)

var roleMap = map[int]string{
	1: "admin",
	2: "user",
	3: "unknown",
}

// User represents user data for response
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

var db, errDb = sql.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/userregistry")

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/v1", getIndexPage).Methods("GET")
	myRouter.HandleFunc("/api/v1/login", login).Methods("POST")
	myRouter.HandleFunc("/api/v1/users", getUsers).Methods("GET")
	myRouter.HandleFunc("/api/v1/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/api/v1/user/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/api/v1/user", addUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8090", myRouter))
}

func main() {
	checkErr(errDb)
	handleRequests()
	defer db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func roleToString(roleInt int) (role string) {
	switch {
	case roleInt == 1:
		role = roleMap[1]
	case roleInt == 2:
		role = roleMap[2]
	default:
		role = roleMap[3]
	}
	return
}

func roleToInt(roleString string) (role int) {
	switch {
	case roleString == roleMap[1]:
		role = 1
	case roleString == roleMap[2]:
		role = 2
	default:
		role = 0
	}
	return
}
