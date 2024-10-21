<<<<<<< HEAD
SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post, s.Salary
FROM employee e
JOIN post p ON e.PostId = p.PostId
JOIN salary s ON e.PostId = s.PostId
WHERE s.Salary >= 60000;

SELECT e.EmployeeId, e.FirstName, e.LastName, e.Birthdate, p.Name AS Post
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE (strftime('%Y', 'now') - strftime('%Y', e.Birthdate)) BETWEEN 30 AND 40;

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE p.Name = 'Manager';

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post
FROM employee e
JOIN hierarchy h ON e.EmployeeId = h.EmployeeId
JOIN employee m ON h.ManagerId = m.EmployeeId
JOIN post p ON e.PostId = p.PostId
WHERE m.FirstName = 'Alice' AND m.LastName = 'Dupont';

SELECT e.EmployeeId, e.FirstName, e.LastName, e.DateOfJoin, p.Name AS Post,
       (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) AS Anciennete
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) > 5;

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post,
       (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) AS Anciennete
FROM employee e
JOIN post p ON e.PostId = p.PostId
JOIN hierarchy h ON e.EmployeeId = h.EmployeeId
JOIN employee m ON h.ManagerId = m.EmployeeId
WHERE p.Name = 'Manager'
  AND (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) > 5
  AND m.FirstName = 'Alice' AND m.LastName = 'Dupont';
=======
SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post, s.Salary
FROM employee e
JOIN post p ON e.PostId = p.PostId
JOIN salary s ON e.PostId = s.PostId
WHERE s.Salary >= 60000;

SELECT e.EmployeeId, e.FirstName, e.LastName, e.Birthdate, p.Name AS Post
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE (strftime('%Y', 'now') - strftime('%Y', e.Birthdate)) BETWEEN 30 AND 40;

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE p.Name = 'Manager';

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post
FROM employee e
JOIN hierarchy h ON e.EmployeeId = h.EmployeeId
JOIN employee m ON h.ManagerId = m.EmployeeId
JOIN post p ON e.PostId = p.PostId
WHERE m.FirstName = 'Alice' AND m.LastName = 'Dupont';

SELECT e.EmployeeId, e.FirstName, e.LastName, e.DateOfJoin, p.Name AS Post,
       (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) AS Anciennete
FROM employee e
JOIN post p ON e.PostId = p.PostId
WHERE (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) > 5;

SELECT e.EmployeeId, e.FirstName, e.LastName, p.Name AS Post,
       (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) AS Anciennete
FROM employee e
JOIN post p ON e.PostId = p.PostId
JOIN hierarchy h ON e.EmployeeId = h.EmployeeId
JOIN employee m ON h.ManagerId = m.EmployeeId
WHERE p.Name = 'Manager'
  AND (strftime('%Y', 'now') - strftime('%Y', e.DateOfJoin)) > 5
  AND m.FirstName = 'Alice' AND m.LastName = 'Dupont';
  
>>>>>>> 08d1b84af5fbed06ba91fd6501c93e4f53478701
