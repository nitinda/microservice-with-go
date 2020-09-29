package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

type userInfoSummary struct {
	ID         int
	UserName   string
	DepartName string
	Created    string
}

type userinfos struct {
	UserInfos []userInfoSummary
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// fmt.Println("# Inserting values")

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "asdasd", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	// fmt.Println("last inserted id =", lastInsertId)

	repos := userinfos{}

	// fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		repo := userInfoSummary{}
		err = rows.Scan(
			&repo.ID,
			&repo.UserName,
			&repo.DepartName,
			&repo.Created,
		)

		repos.UserInfos = append(repos.UserInfos, repo)
	}

	// for rows.Next() {
	// 	var uid int
	// 	var username string
	// 	var department string
	// 	var created time.Time
	// 	err = rows.Scan(&uid, &username, &department, &created)
	// 	checkErr(err)
	// 	fmt.Println("uid | username | department | created ")
	// 	fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
	// }

	out, err := json.Marshal(repos)
	if err != nil {
		fmt.Println("============Error=======")
	}

	fmt.Fprintf(os.Stdout, string(out))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
