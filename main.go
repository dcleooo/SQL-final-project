package main

import (
    "fmt"
    "log"
	"SQL-final-project/src"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/", src.HomeHandler)
    http.HandleFunc("/add-employee", src.AddEmployeeHandler)
    http.HandleFunc("/delete-employee", src.DeleteEmployeeHandler)
    http.HandleFunc("/departments", src.DepartmentsHandler)
    http.HandleFunc("/posts", src.PostsHandler)
    http.HandleFunc("/managers", src.ManagerHandler)
    http.HandleFunc("/edit-employee", src.EditEmployeeHandler)
    fmt.Println("Starting server on :8080")
    fmt.Println("http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}