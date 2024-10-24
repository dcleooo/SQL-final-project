package src 

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"html/template"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    template, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    template.Execute(w, nil)
}

func EmployeeHandler(w http.ResponseWriter, r *http.Request) {
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
        emp.Birthdate = emp.Birthdate[:10]
        employees = append(employees, emp)
    }

    data := PageData{
        Employees: employees,
    }

    tmpl, err := template.ParseFiles("./templates/employee.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
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
        emp.Birthdate = emp.Birthdate[:10]
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

    tmpl, err := template.ParseFiles("templates/add.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func AddEmployeeHandler(w http.ResponseWriter, r *http.Request) {
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

        var count int
        err = db.QueryRow("SELECT COUNT(*) FROM hierarchy").Scan(&count)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if count == 0 {
            if r.FormValue("post_id") == "1" {
                _, err = db.Exec("INSERT INTO hierarchy (EmployeeId, ManagerId) SELECT EmployeeId, ? FROM employee WHERE EmployeeId != ?", employeeId, employeeId)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
            }
        } else {
            var managerId int
            err = db.QueryRow(`
                SELECT employee.EmployeeId
                FROM employee 
                LEFT JOIN hierarchy ON employee.EmployeeId = hierarchy.ManagerId
                WHERE employee.PostId = 1
                GROUP BY employee.EmployeeId
                ORDER BY COUNT(hierarchy.EmployeeId) ASC
                LIMIT 1
            `).Scan(&managerId)
            if err != nil {
                if err == sql.ErrNoRows {
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
            if r.FormValue("post_id") != "1" {
                _, err = db.Exec("INSERT INTO hierarchy (EmployeeId, ManagerId) VALUES (?, ?)", employeeId, managerId)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
            }
        }

        var responsableId sql.NullInt64
        departmentId := r.FormValue("departement_id")
        err = db.QueryRow("SELECT ResponsableId FROM departement WHERE DepartementId = ?", departmentId).Scan(&responsableId)
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if err == sql.ErrNoRows || !responsableId.Valid {
            _, err = db.Exec("UPDATE departement SET ResponsableId = ? WHERE DepartementId = ?", employeeId, departmentId)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    if r.Method == "GET" {
        r.ParseForm()
        employeeId := r.URL.Query().Get("id")
        redirectURL := r.URL.Query().Get("redirect")

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
                LEFT JOIN hierarchy h ON employee.EmployeeId = h.ManagerId
                WHERE employee.PostId = 1 AND employee.EmployeeId != ?
                GROUP BY employee.EmployeeId
                ORDER BY COUNT(h.EmployeeId) ASC
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

        var DepartementId sql.NullInt64
        err = db.QueryRow("SELECT COUNT(*), DepartementId FROM departement WHERE ResponsableId = ?", employeeId).Scan(&count, &DepartementId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        isDepartmentHead := count > 0

        if isDepartmentHead && DepartementId.Valid {
            var newDepartmentHeadId int
            err = db.QueryRow(`
                SELECT employee.EmployeeId
                FROM employee
                WHERE employee.DepartementId = ? AND employee.EmployeeId != ?
                ORDER BY DateOfJoin ASC
                LIMIT 1
            `, DepartementId.Int64, employeeId).Scan(&newDepartmentHeadId)
            if err != nil {
                if err == sql.ErrNoRows {
                    _,err = db.Exec("UPDATE departement SET ResponsableId = NULL WHERE DepartementId = ?", DepartementId.Int64)
                    if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            _, err = db.Exec("UPDATE departement SET ResponsableId = ? WHERE DepartementId = ?", newDepartmentHeadId, DepartementId.Int64)
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

        if redirectURL == "" {
            redirectURL = "/"
        }
        http.Redirect(w, r, redirectURL, http.StatusSeeOther)
    }
}

func DepartmentsHandler(w http.ResponseWriter, r *http.Request) {
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
            emp.Birthdate = emp.Birthdate[:10]
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

func PostsHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    postRows, err := db.Query("SELECT PostId, Name FROM post")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer postRows.Close()

    type Post struct {
        PostId int
        Name   string
        Employees []Employee
    }

    var posts []Post

    for postRows.Next() {
        var post Post
        err := postRows.Scan(&post.PostId, &post.Name)
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
            employee.Email
        FROM 
            employee
            INNER JOIN departement ON employee.DepartementId = departement.DepartementId
            INNER JOIN post ON employee.PostId = post.PostId
        WHERE
            employee.PostId = ?
        `, post.PostId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer empRows.Close()

        var employees []Employee
        for empRows.Next() {
            var emp Employee
            err := empRows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.DepartementName, &emp.PostName, &emp.DateOfJoin, &emp.Phone, &emp.Email)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            emp.Birthdate = emp.Birthdate[:10]
            employees = append(employees, emp)
        }
        post.Employees = employees
        posts = append(posts, post)
    }

    data := struct {
        Posts []Post
    }{
        Posts: posts,
    }

    tmpl, err := template.ParseFiles("templates/posts.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func ManagerHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    managersRows, err := db.Query(`
        SELECT
            employee.EmployeeId,
            employee.FirstName,
            employee.LastName
        FROM
            employee
            INNER JOIN post ON employee.PostId = post.PostId
        WHERE 
            post.PostId = 1
    `)
            
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer managersRows.Close()

    type Manager struct {
        EmployeeId int
        FirstName string
        LastName  string
        Employees []Employee
    }

    var managers []Manager

    for managersRows.Next() {
        var manager Manager
        err := managersRows.Scan(&manager.EmployeeId, &manager.FirstName, &manager.LastName)
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
            departement.Name AS DepartementName,
            post.Name AS PostName, 
            employee.DateOfJoin, 
            employee.Phone, 
            employee.Email
        FROM 
            employee
            INNER JOIN post ON employee.PostId = post.PostId
            INNER JOIN departement ON employee.DepartementId = departement.DepartementId
            INNER JOIN hierarchy ON employee.EmployeeId = hierarchy.EmployeeId
        WHERE
            hierarchy.ManagerId = ?
        `, manager.EmployeeId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer empRows.Close()

        var employees []Employee
        for empRows.Next() {
            var emp Employee
            err := empRows.Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.DepartementName, &emp.PostName, &emp.DateOfJoin, &emp.Phone, &emp.Email)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            emp.Birthdate = emp.Birthdate[:10]
            employees = append(employees, emp)
        }
        manager.Employees = employees
        managers = append(managers, manager)
    }
    data := struct {
        Managers []Manager
    }{
        Managers: managers,
    }

    tmpl, err := template.ParseFiles("templates/managers.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func EditEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    defer db.Close()

    if r.Method == "GET" {
        employeeId := r.URL.Query().Get("id")
        redirectURL := r.URL.Query().Get("redirect")

        // Fetch employee details
        var emp Employee
        err := db.QueryRow(`
            SELECT 
                employee.EmployeeId, 
                employee.FirstName, 
                employee.LastName, 
                employee.Birthdate, 
                employee.Address, 
                post.Name AS PostName, 
                employee.DateOfJoin, 
                employee.Phone, 
                employee.Email,
                departement.Name AS DepartmentName,
                post.Name AS PostName,
                employee.DepartementId,
                employee.PostId,
                CASE WHEN departement.ResponsableId = employee.EmployeeId THEN 1 ELSE 0 END AS IsHead
            FROM 
                employee
                INNER JOIN departement ON employee.DepartementId = departement.DepartementId
                INNER JOIN post ON employee.PostId = post.PostId
            WHERE 
                employee.EmployeeId = ?
        `, employeeId).Scan(&emp.EmployeeId, &emp.FirstName, &emp.LastName, &emp.Birthdate, &emp.Address, &emp.PostName, &emp.DateOfJoin, &emp.Phone, &emp.Email, &emp.DepartementName, &emp.PostName,&emp.DepartementId, &emp.PostId, &emp.IsHead )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        emp.Birthdate = emp.Birthdate[:10]
        emp.DateOfJoin = emp.DateOfJoin[:10]

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

        data := struct {
            Employee    Employee
            Departments []Department
            Posts       []Post
            RedirectURL string
        }{
            Employee:    emp,
            RedirectURL: redirectURL,
            Departments: departments,
            Posts:       posts,
        }

        tmpl, err := template.ParseFiles("templates/edit_employee.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else if r.Method == "POST" {
        r.ParseForm()
        employeeId := r.FormValue("employee_id")
        redirectURL := r.FormValue("redirect")

        var currentDepartmentId int
        err := db.QueryRow("SELECT DepartmentId FROM department WHERE ResponsableId = ?", employeeId).Scan(&currentDepartmentId)
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var currentManagerId int
        err = db.QueryRow("SELECT ManagerId FROM hierarchy WHERE EmployeeId = ?", employeeId).Scan(&currentManagerId)
        if err != nil && err != sql.ErrNoRows {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        newDepartmentId, err := strconv.Atoi(r.FormValue("department_id"))
        if err != nil {
            http.Error(w, "Invalid department ID", http.StatusBadRequest)
            return
        }

        _, err = db.Exec(`
            UPDATE employee 
            SET FirstName = ?, LastName = ?, Birthdate = ?, Address = ?, DepartmentId = ?, PostId = ?, DateOfJoin = ?, Phone = ?, Email = ?
            WHERE EmployeeId = ?
        `, r.FormValue("first_name"), r.FormValue("last_name"), r.FormValue("birthdate"), r.FormValue("address"), r.FormValue("department_id"), r.FormValue("post_id"), r.FormValue("date_of_join"), r.FormValue("phone"), r.FormValue("email"), employeeId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if currentDepartmentId != 0 && currentDepartmentId != newDepartmentId {
            var newHeadId int
            err = db.QueryRow(`
                SELECT EmployeeId 
                FROM employee 
                WHERE DepartmentId = ? 
                ORDER BY DateOfJoin ASC 
                LIMIT 1
            `, currentDepartmentId).Scan(&newHeadId)
            if err != nil {
                if err == sql.ErrNoRows {
                    _,err = db.Exec("UPDATE department SET ResponsableId = NULL WHERE DepartmentId = ?", currentDepartmentId)
                    if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                    return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            _, err = db.Exec("UPDATE department SET ResponsableId = ? WHERE DepartmentId = ?", newHeadId, currentDepartmentId)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        newPostId, err := strconv.Atoi(r.FormValue("post_id"))
        if err != nil {
            http.Error(w, "Invalid post ID", http.StatusBadRequest)
            return
        }

        if currentManagerId != 0 && currentManagerId != newPostId {
            var newManagerId int
            err = db.QueryRow(`
                SELECT EmployeeId 
                FROM employee 
                WHERE PostId = ? 
                ORDER BY DateOfJoin ASC 
                LIMIT 1
            `, currentManagerId).Scan(&newManagerId)
            if err != nil {
                if err == sql.ErrNoRows {
                    _,err = db.Exec("DELETE FROM hierarchy WHERE ManagerId = ?", currentManagerId)
                    if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            // Update the old manager's hierarchy
            _, err = db.Exec("UPDATE hierarchy SET ManagerId = ? WHERE ManagerId = ?", newManagerId, currentManagerId)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        if redirectURL == "" {
            redirectURL = "/"
        }
        http.Redirect(w, r, redirectURL, http.StatusSeeOther)
    }
}