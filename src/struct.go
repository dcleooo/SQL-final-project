package src

type Employee struct {
    EmployeeId     int
    FirstName      string
    LastName       string
    Birthdate      string
    Address        string
    DepartementName string
    PostName       string
    DateOfJoin     string
    Phone          string
    Email          string
    IsHead         bool
    PostId         int
    DepartementId   int
}

type Department struct {
    DepartementId int
    Name         string
}

type Post struct {
    PostId int
    Name   string
}

type Manager struct {
    PostId         int
    PostName       string
    DepartementId   int
    DepartementName string
}

type PageData struct {
    Employees   []Employee
    Departments []Department
    Posts       []Post
}