<<<<<<< HEAD
=======
<<<<<<< HEAD
>>>>>>> 91793a69d9e7820d45b4db0d0a6da92b808a2da4
PRAGMA foreign_keys = ON;

CREATE TABLE employee (
    EmployeeId INTEGER PRIMARY KEY AUTOINCREMENT,
    FirstName VARCHAR(50),
    LastName VARCHAR(50),
    Birthdate DATE,
    Address VARCHAR(255),
    DepartementId INTEGER,
    PostId INTEGER,
    DateOfJoin DATE,
    Phone VARCHAR(15),
    Email VARCHAR(100),
    FOREIGN KEY (DepartementId) REFERENCES departement(DepartementId),
    FOREIGN KEY (PostId) REFERENCES post(PostId)
);

CREATE TABLE departement (
    DepartementId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name VARCHAR(50),
    ResponsableId INTEGER,
    FOREIGN KEY (ResponsableId) REFERENCES employee(EmployeeId)
);

CREATE TABLE post (
    PostId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name VARCHAR(50)
);

CREATE TABLE salary (
    PostId INTEGER,
    Salary FLOAT,
    FOREIGN KEY (PostId) REFERENCES post(PostId)
);

CREATE TABLE hierarchy (
    EmployeeId INTEGER,
    ManagerId INTEGER,
    FOREIGN KEY (EmployeeId) REFERENCES employee(EmployeeId),
    FOREIGN KEY (ManagerId) REFERENCES employee(EmployeeId)
);

INSERT INTO departement (Name, ResponsableId) VALUES
    ('HR', NULL),
    ('Engineering', NULL),
    ('Marketing', NULL);

INSERT INTO post (Name) VALUES
    ('Manager'),
    ('Engineer'),
    ('Technician'),
    ('HR Specialist'),
    ('Marketing Specialist');

INSERT INTO salary (PostId, Salary) VALUES
    (1, 70000.00),
    (2, 60000.00),
    (3, 50000.00),
    (4, 55000.00),
    (5, 50000.00);

INSERT INTO employee (FirstName, LastName, Birthdate, Address, DepartementId, PostId, DateOfJoin, Phone, Email) VALUES
    ('Alice', 'Dupont', '1985-04-12', '12 Rue des Lilas, Paris', 1, 1, '2020-01-15', '0601020304', 'alice.dupont@example.com'),
    ('Bob', 'Martin', '1990-07-19', '34 Avenue de la République, Lyon', 2, 2, '2018-03-20', '0605060708', 'bob.martin@example.com'),
    ('Claire', 'Leroy', '1982-11-05', '56 Boulevard Victor Hugo, Marseille', 3, 5, '2015-10-10', '0610101112', 'claire.leroy@example.com'),
    ('David', 'Bernard', '1992-03-22', '78 Rue de la Liberté, Bordeaux', 1, 4, '2019-05-12', '0613141516', 'david.bernard@example.com'),
    ('Emma', 'Durand', '1988-09-09', '89 Place de la Bourse, Nantes', 1, 3, '2017-11-03', '0617181920', 'emma.durand@example.com');

UPDATE departement SET ResponsableId = (SELECT EmployeeId FROM employee WHERE FirstName = 'Alice' AND LastName = 'Dupont') WHERE Name = 'HR';
UPDATE departement SET ResponsableId = (SELECT EmployeeId FROM employee WHERE FirstName = 'Bob' AND LastName = 'Martin') WHERE Name = 'Engineering';
UPDATE departement SET ResponsableId = (SELECT EmployeeId FROM employee WHERE FirstName = 'Claire' AND LastName = 'Leroy') WHERE Name = 'Marketing';

INSERT INTO hierarchy (EmployeeId, ManagerId) VALUES
    ((SELECT EmployeeId FROM employee WHERE FirstName = 'Bob' AND LastName = 'Martin'), (SELECT EmployeeId FROM employee WHERE FirstName = 'Alice' AND LastName = 'Dupont')),
    ((SELECT EmployeeId FROM employee WHERE FirstName = 'Claire' AND LastName = 'Leroy'), (SELECT EmployeeId FROM employee WHERE FirstName = 'Alice' AND LastName = 'Dupont')),
    ((SELECT EmployeeId FROM employee WHERE FirstName = 'David' AND LastName = 'Bernard'), (SELECT EmployeeId FROM employee WHERE FirstName = 'Alice' AND LastName = 'Dupont')),
