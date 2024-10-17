package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	EmployeeId   string
    FirstName    string
    LastName     string
    Birthdate    string
    Address      string
    DepartementId string
    PostId       string
    DateOfJoin   string
    Phone        string
    Email        string
}

func dbConn() (db *sql.DB) {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    rows, err := db.Query("SELECT EmployeeId, FirstName, LastName, Birthdate, Address, DepartementId, PostId, DateOfJoin, Phone, Email FROM employee")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var emp Employee
        err := rows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.DepartementId, &emp.PostId, &emp.DateOfJoin, &emp.Phone, &emp.Email)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employees = append(employees, emp)
    }

    tmpl, err := template.ParseFiles("static/html/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, employees)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}