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
    EmployeeId   int
    FirstName    string
    LastName     string
    Birthdate    string
    Address      string
    DepartementName string
    PostName     string
    DateOfJoin   string
    Phone        string
    Email        string
}

type Department struct {
    DepartementId int
    Name          string
}

type Post struct {
    PostId int
    Name   string
}

type PageData struct {
    Employees   []Employee
    Departments []Department
    Posts       []Post
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

    rows, err := db.Query(`
        SELECT 
            employee.EmployeeId, 
            employee.FirstName, 
            employee.LastName, 
            employee.Birthdate, 
            employee.Address, 
            departement.Name AS DepartmentName, 
            post.Name AS PostName, 
            employee.DateOfJoin, 
            employee.Phone, 
            employee.Email 
        FROM 
            employee
            INNER JOIN departement ON employee.DepartementId = departement.DepartementId
            INNER JOIN post ON employee.PostId = post.PostId
    `)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var emp Employee
        err := rows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.DepartementName, &emp.PostName, &emp.DateOfJoin, &emp.Phone, &emp.Email)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employees = append(employees, emp)
    }

    deptRows, err := db.Query("SELECT DepartementId, Name FROM departement")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer deptRows.Close()

    var departments []Department
    for deptRows.Next() {
        var dept Department
        err := deptRows.Scan(&dept.DepartementId, &dept.Name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        departments = append(departments, dept)
    }

    postRows, err := db.Query("SELECT PostId, Name FROM post")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer postRows.Close()

    var posts []Post
    for postRows.Next() {
        var post Post
        err := postRows.Scan(&post.PostId, &post.Name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    data := PageData{
        Employees:   employees,
        Departments: departments,
        Posts:       posts,
    }

    tmpl, err := template.ParseFiles("static/html/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    if r.Method == "POST" {
        r.ParseForm()
        result, err := db.Exec("INSERT INTO employee (FirstName, LastName, Birthdate, Address, DepartementId, PostId, DateOfJoin, Phone, Email) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
            r.FormValue("first_name"), r.FormValue("last_name"), r.FormValue("birthdate"), r.FormValue("address"), r.FormValue("departement_id"), r.FormValue("post_id"), r.FormValue("date_of_join"), r.FormValue("phone"), r.FormValue("email"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        employeeId, err := result.LastInsertId()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var managerId int
        err = db.QueryRow(`
            SELECT ResponsableId
            FROM departement
            WHERE DepartementId = ?
        `, r.FormValue("departement_id")).Scan(&managerId)
        if err != nil {
            if err == sql.ErrNoRows {
                // No manager found, delete the newly added employee
                _, delErr := db.Exec("DELETE FROM employee WHERE EmployeeId = ?", employeeId)
                if delErr != nil {
                    http.Error(w, delErr.Error(), http.StatusInternalServerError)
                    return
                }
                http.Error(w, "No manager found for the selected department. Employee not added.", http.StatusBadRequest)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        _, err = db.Exec(`
            INSERT INTO hierarchy (EmployeeId, ManagerId)
            VALUES (?, ?)
        `, employeeId, managerId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", indexHandler)
    http.HandleFunc("/add-employee", addEmployeeHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}