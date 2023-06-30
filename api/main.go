package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type FormData struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Manager      string `json:"manager"`
	Team     string `json:"team"`
	LeaveType  string `json:"leave_type"`
	LeaveDate   string `json:"leave_date"`
	}
type ApprovedData struct {
	LeaveId string `json:"leave_id"`
	ReportingManager string `json:"reporting_manager"`
	Apporoved string `json:"approved"`
	EmployeeId string `json:"employee_id"`
}


var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "xenonstack"
)

func main() {
	var err error
	db, err = sql.Open("postgres", getDBConnectionString())
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	// Create the form_data table if it doesn't exist
	createTableIfNotExists(db)

	// Set up Gin
	r := gin.Default()
	// Define the API route to save the form data
	r.POST("/api/save", func(c *gin.Context) {
		var formData FormData

		if err := c.ShouldBindJSON(&formData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the form data into the database
		if err := insertFormData(db, formData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Form data saved successfully"})
	})
	r.GET("/users", func(c *gin.Context) {
		users, err := getUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	r.GET("/approved", func(c *gin.Context) {
		users, err := getApproveData(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	
	// Start the server
	log.Println("Server listening on http://localhost:8000")
	log.Fatal(r.Run(":8000"))
}

func getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func createTableIfNotExists(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS employee_leave_data (
		id SERIAL PRIMARY KEY,
		employee_name TEXT NOT NULL,
		manager_name TEXT NOT NULL,
		team_name TEXT NOT NULL,
		leave_type TEXT NOT NULL,
		leave_dates TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func insertFormData(db *sql.DB, formData FormData) error {
	_, err := db.Exec("INSERT INTO employee_leave_data ( employee_name, manager_name,  team_name,leave_type, leave_dates) VALUES ($1, $2, $3, $4, $5)",
		formData.Name, formData.Manager, formData.Team, formData.LeaveType, formData.LeaveDate)
	if err != nil {
		log.Println("Failed to insert form data:", err)
		return err
	}
	return nil
}
func getUsers(db *sql.DB) ([]FormData, error) {
	rows, err := db.Query("SELECT id, employee_name, manager_name, team_name, leave_type,  leave_dates FROM employee_leave_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]FormData, 0)

	for rows.Next() {
		var user FormData
		err := rows.Scan(&user.Id, &user.Name, &user.Manager,&user.Team, &user.LeaveType, &user.LeaveDate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func getApproveData(db *sql.DB) ([]ApprovedData, error) {
	rows, err := db.Query("SELECT leave_id, reporting_manager, approved,employee_id FROM notifications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]ApprovedData, 0)

	for rows.Next() {
		var user ApprovedData
		err := rows.Scan(&user.LeaveId, &user.ReportingManager, &user.Apporoved, &user.EmployeeId )
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

