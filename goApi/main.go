package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
	const maxRetries = 100 // se que es absurdamente alto pero me fallo en conectarse en 5, 10, 15 y 25 intentos porlomenos una vez durante mis pruebas entonces ¯\_(ツ)_/¯
	const retryDelay = 2 * time.Second

	dsn := "root:1234@tcp(db:3306)/classicmodels"
	var db *sql.DB
	var err error

	// aqui uso las variables de arriba para intentar conectarme a la base de datos, si funciona se sale del loop
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		fmt.Printf("Database not ready, retrying in %v... (attempt %d)\n", retryDelay, i+1)
		time.Sleep(retryDelay) // Espera antes de volver a intentar la conectarse
	}

	if err != nil {
		log.Fatalf("Error connecting to the database after %d attempts: %v", maxRetries, err)
		return
	}
	defer db.Close()

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
