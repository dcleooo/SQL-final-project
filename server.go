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
    IsHead       bool
}

type Department struct {
    DepartementId int
    Name          string
}

type Post struct {
    PostId int
    Name   string
}

type Manager struct {
    PostId       int
    PostName     string
    DepartmentId int
    DepartmentName string
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

    tmpl, err := template.ParseFiles("templates/home.html")
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
            SELECT employee.EmployeeId
            FROM employee 
            LEFT JOIN hierarchy  ON employee.EmployeeId = hierarchy.ManagerId
            WHERE employee.PostId = 1
            GROUP BY employee.EmployeeId
            ORDER BY COUNT(hierarchy.EmployeeId) ASC
            LIMIT 1
        `).Scan(&managerId)
        if err != nil {
            if err == sql.ErrNoRows {
                // No manager found, delete the newly added employee
                _, delErr := db.Exec("DELETE FROM employee WHERE EmployeeId = ?", employeeId)
                if delErr != nil {
                    http.Error(w, delErr.Error(), http.StatusInternalServerError)
                    return
                }
                http.Error(w, "No manager found. Employee not added.", http.StatusBadRequest)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Insert into hierarchy
        _, err = db.Exec("INSERT INTO hierarchy (EmployeeId, ManagerId) VALUES (?, ?)", employeeId, managerId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    if r.Method == "POST" {
        r.ParseForm()
        employeeId := r.FormValue("employee_id")

        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM hierarchy WHERE ManagerId = ?", employeeId).Scan(&count)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        isManager := count > 0

        if isManager {
            var newManagerId int
            err = db.QueryRow(`
                SELECT employee.EmployeeId
                FROM employee 
                LEFT JOIN hierarchy ON employee.EmployeeId = hierarchy.ManagerId
                WHERE employee.PostId = 1 AND employee.EmployeeId != ?
                GROUP BY employee.EmployeeId
                ORDER BY COUNT(hierarchy.EmployeeId) ASC
                LIMIT 1
            `, employeeId).Scan(&newManagerId)
            if err != nil {
                if err == sql.ErrNoRows {
                    http.Error(w, "No other manager found to reassign direct reports.", http.StatusBadRequest)
                    return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            _, err = db.Exec("UPDATE hierarchy SET ManagerId = ? WHERE ManagerId = ?", newManagerId, employeeId)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
        var DepartementId int
        err = db.QueryRow("SELECT COUNT(*), DepartementId FROM departement WHERE ResponsableId = ?", employeeId).Scan(&count, &DepartementId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        isDepartmentHead := count > 0

        if isDepartmentHead {
            var newDepartmentHeadId int
            err = db.QueryRow(`
                SELECT employee.EmployeeId
                FROM employee
                WHERE employee.DepartementId = ? AND employee.EmployeeId != ?
                ORDER BY DateOfJoin ASC
                LIMIT 1
            `, DepartementId, employeeId).Scan(&newDepartmentHeadId)
            if err != nil {
                if err == sql.ErrNoRows {
                    http.Error(w, "No other employee found to reassign department head.", http.StatusBadRequest)
                    return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            _, err = db.Exec("UPDATE departement SET ResponsableId = ? WHERE DepartementId = ?", newDepartmentHeadId, DepartementId)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        _, err = db.Exec("DELETE FROM hierarchy WHERE EmployeeId = ?", employeeId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        _, err = db.Exec("DELETE FROM employee WHERE EmployeeId = ?", employeeId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func departmentsHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    deptRows, err := db.Query("SELECT DepartementId, Name FROM departement")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer deptRows.Close()

    type Department struct {
        DepartementId int
        Name          string
        Employees     []Employee
    }

    var departments []Department

    for deptRows.Next() {
        var dept Department
        err := deptRows.Scan(&dept.DepartementId, &dept.Name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        empRows, err := db.Query(`
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
            employee.Email,
            CASE WHEN departement.ResponsableId = employee.EmployeeId THEN 1 ELSE 0 END AS IsHead
        FROM 
            employee
            INNER JOIN departement ON employee.DepartementId = departement.DepartementId
            INNER JOIN post ON employee.PostId = post.PostId
        WHERE
            employee.DepartementId = ?
        `, dept.DepartementId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer empRows.Close()

        var employees []Employee
        for empRows.Next() {
            var emp Employee
            err := empRows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.DepartementName, &emp.PostName, &emp.DateOfJoin, &emp.Phone, &emp.Email, &emp.IsHead)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            employees = append(employees, emp)
        }
        dept.Employees = employees
        departments = append(departments, dept)
    }

    data := struct {
        Departments []Department
    }{
        Departments: departments,
    }

    tmpl, err := template.ParseFiles("templates/departments.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)
    http.HandleFunc("/add-employee", addEmployeeHandler)
    http.HandleFunc("/delete-employee", deleteEmployeeHandler)
    http.HandleFunc("/departments",departmentsHandler)
    //http.HandleFunc("/posts",postsHandler)
    //http.HandleFunc("/manager",managerHandler)
    //http.HandleFunc("/employees",employeesHandler)
	fmt.Println("Starting server on :8080")
    fmt.Println("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}