package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const indexPage = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/api/v1/login" method="post">
            Username:<input type="text" name="username">
            Password:<input type="password" name="password">
            <input type="submit" value="Login">
        </form>
    </body>
</html>
`

// UserFromRequest represents user data from request
type UserFromRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func getIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexPage)
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userName := r.FormValue("username")
	password := r.FormValue("password")

	passHash := toHashPassString(password)

	rows, err := db.Query("SELECT * FROM user WHERE name =? AND password=?", userName, passHash)
	checkErr(err)

	var users []User

	for rows.Next() {
		var password string
		var createTime string
		var role int
		var tmpUser User

		err = rows.Scan(&tmpUser.ID, &tmpUser.Name, &tmpUser.Email, &password, &createTime, &role)
		checkErr(err)

		tmpUser.Role = roleToString(role)
		users = append(users, tmpUser)
	}

	redirectTarget := "/"
	if len(users) == 0 || users[0].Name != userName {
		fmt.Fprintf(w, "401 unauthorized, wrong credentials")
	} else {
		user := users[0]
		// passHasStr := base64.URLEncoding.EncodeToString(passHash)
		// passHasStr := passHash
		userTokenRoleMap[passHash] = user.Role

		// w.Header().Set("Authorization", passHash)

		redirectTarget = "/api/v1/users?token=" + passHash
		http.Redirect(w, r, redirectTarget, 302)
	}

}

func addUser(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		fmt.Fprintf(w, "401 unauthorized")
	}
	role := userTokenRoleMap[token]

	decoder := json.NewDecoder(r.Body)
	var newUser UserFromRequest
	err := decoder.Decode(&newUser)
	checkErr(err)

	stmt, err := db.Prepare("INSERT user SET name=?,email=?,password=?,role=?")
	checkErr(err)

	password := toHashPassString(newUser.Password)

	var roleID int
	switch {
	case role == roleMap[1]:
		roleID = roleToInt(newUser.Role)
	case role == roleMap[2]:
		roleID = roleToInt(role)
	default:
		fmt.Fprintf(w, "401 unauthorized, wrong credentials")
	}

	res, err := stmt.Exec(newUser.Name, newUser.Email, password, roleID)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	log.Printf("New user added with id = %d", id)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		fmt.Fprintf(w, "401 unauthorized")
	}
	role := userTokenRoleMap[token]

	if role != roleMap[1] {
		fmt.Fprintf(w, "401 unauthorized")
	}

	id := r.URL.Path[len("/api/v1/user/"):]

	stmt, err := db.Prepare("delete from user where id=?")
	checkErr(err)

	res, err := stmt.Exec(id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	log.Printf("User with id = %v is deleted: %v", id, affect)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update User Endpoint Hit")
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		fmt.Fprintf(w, "401 unauthorized, wrong credentials")
	}
	role := userTokenRoleMap[token]

	var users []User

	switch {
	case role == roleMap[1]:
		rows, err := db.Query("SELECT * FROM user")
		checkErr(err)
		users = fillUsers(rows)
	case role == roleMap[2]:
		rows, err := db.Query("SELECT * FROM user WHERE role=?", roleToInt(role))
		checkErr(err)
		users = fillUsers(rows)
	default:
		fmt.Fprintf(w, "401 unauthorized, wrong credentials")
	}

	data, err := json.Marshal(users)
	checkErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func fillUsers(rows *sql.Rows) (users []User) {
	for rows.Next() {
		var password []byte
		var createTime string
		var role int
		var tmpUser User

		err := rows.Scan(&tmpUser.ID, &tmpUser.Name, &tmpUser.Email, &password, &createTime, &role)
		checkErr(err)

		tmpUser.Role = roleToString(role)
		users = append(users, tmpUser)
	}
	return
}

func toHashPassString(passwordStr string) (passwordHashStr string /*[]byte*/) {
	return base64.URLEncoding.EncodeToString([]byte(passwordStr))
	// passwordHash := sha256.New()
	// passwordHash.Write([]byte(passwordStr))
	// return passwordHash.Sum(nil)
}
