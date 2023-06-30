package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type FormData struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Team     string `json:"team"`
	LeaveType  string `json:"leave_type"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	ReportTo   string `json:"report_to"`
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
	// Apply CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} 
	r.Use(cors.New(config))

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
	r.PUT("/api/approve/:id", func(c *gin.Context) {
		leaveID := c.Param("id")
	
		// Update the approval status in the database
		if err := updateApprovalStatus(db, leaveID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Approval status updated successfully"})
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS emp_leave_data (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		mobile TEXT NOT NULL,
		team TEXT NOT NULL,
		leave_type TEXT NOT NULL,
		from_date TEXT NOT NULL,
		to_date TEXT NOT NULL,
		report_to TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func insertFormData(db *sql.DB, formData FormData) error {
	_, err := db.Exec("INSERT INTO emp_leave_data (id, name, email,  mobile, team,leave_type, from_date, to_date, report_to) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		formData.Id, formData.Name, formData.Email, formData.Mobile, formData.Team, formData.LeaveType, formData.FromDate,formData.ToDate, formData.ReportTo)
	if err != nil {
		log.Println("Failed to insert form data:", err)
		return err
	}
	return nil
}
func getUsers(db *sql.DB) ([]FormData, error) {
	rows, err := db.Query("SELECT id, name, email, mobile, team, leave_type,  from_date, to_date, report_to FROM emp_leave_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]FormData, 0)

	for rows.Next() {
		var user FormData
		err := rows.Scan(&user.Id, &user.Name, &user.Email,&user.Mobile, &user.Team, &user.LeaveType, &user.FromDate, &user.ToDate, &user.ReportTo)
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
func updateApprovalStatus(db *sql.DB, leaveID string) error {
	stmt, err := db.Prepare("UPDATE notifications SET approved = $1 WHERE leave_id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, leaveID)
	if err != nil {
		return err
	}

	return nil
}

