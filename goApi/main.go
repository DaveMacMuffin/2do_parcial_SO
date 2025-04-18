package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Employee represents an employee record with the desired fields.
type Employee struct {
	EmployeeNumber int
	LastName       string
	FirstName      string
	Email          string
	JobTitle       string
}

func getEmployees(db *sql.DB) ([]Employee, error) {
	rows, err := db.Query(
		"SELECT employeeNumber, lastName, firstName, email, jobTitle FROM employees",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(
			&e.EmployeeNumber,
			&e.LastName,
			&e.FirstName,
			&e.Email,
			&e.JobTitle,
		); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func employeesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	employees, err := getEmployees(db)
	if err != nil {
		http.Error(w, "Failed to fetch employees: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Plain‑text response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, e := range employees {
		fmt.Fprintf(
			w,
			"Employee %d: %s %s — %s — %s\n",
			e.EmployeeNumber,
			e.FirstName,
			e.LastName,
			e.Email,
			e.JobTitle,
		)
	}
}

func main() {
	const maxRetries = 100
	const retryDelay = 2 * time.Second // Duration to wait between retries

	// DSN for connecting to the MySQL classicmodels database
	dsn := "root:1234@tcp(db:3306)/classicmodels"
	var db *sql.DB
	var err error

	// Retry logic for establishing a database connection
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			// Try to ping the database to check if it's available
			err = db.Ping()
			if err == nil {
				break // Successful connection
			}
		}

		fmt.Printf("Database not ready, retrying in %v... (attempt %d)\n", retryDelay, i+1)
		time.Sleep(retryDelay) // Wait before retrying
	}

	if err != nil {
		log.Fatalf("Error connecting to the database after %d attempts: %v", maxRetries, err)
		return
	}
	defer db.Close() // Defer closing the database connection

	fmt.Println("Successfully connected to the database!")

	http.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		employeesHandler(db, w, r)
	})

	fmt.Println("Server started at http://localhost:8080/employees")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
