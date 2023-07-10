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
	EmployeeName       string `json:"employee_name"`
	ManagerName      string `json:"manager_name"`
	TeamName     string `json:"team_name"`
	LeaveType  string `json:"leave_type"`
	LeaveDate   string `json:"leave_date"`
	LeaveDuration string `json:"leave_duration"`

	}
type ApprovedData struct {
	LeaveId string `json:"leave_id"`
	ReportingManager string `json:"reporting_manager"`
	Apporoved string `json:"approved"`
	}
type Kpi3 struct {
	EmployeeName string `json:"employee_name"`
	LeaveCount string `json:"leave_count"`
}
type Kpi4 struct {
	ManagerName string `json:"manager_name"`
	LeaveCount string `json:"leave_count"`
}
type Kpi6 struct {
	LeaveType string `json:"leave_type"`
	LeaveCount string `json:"leave_count"`
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
	r.GET("/kpi3", func(c *gin.Context) {
		users, err := kpi3(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	r.GET("/kpi4", func(c *gin.Context) {
		users, err := kpi4(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	r.GET("/kpi6One", func(c *gin.Context) {
		users, err := kpi6One(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	r.GET("/kpi6Two", func(c *gin.Context) {
		users, err := kpi6Two(db)
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
		id INT NOT NULL,
		employee_name TEXT NOT NULL,
		manager_name TEXT NOT NULL,
		team_name TEXT NOT NULL,
		leave_type TEXT NOT NULL,
		leave_dates TEXT NOT NULL,
		leave_duration TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func insertFormData(db *sql.DB, formData FormData) error {
	_, err := db.Exec("INSERT INTO employee_leave_data ( id,employee_name, manager_name,  team_name,leave_type, leave_dates, leave_duration) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		formData.Id, formData.EmployeeName, formData.ManagerName, formData.TeamName, formData.LeaveType, formData.LeaveDate, formData.LeaveDuration)
	if err != nil {
		log.Println("Failed to insert form data:", err)
		return err
	}
	return nil
}
func getUsers(db *sql.DB) ([]FormData, error) {
	rows, err := db.Query("SELECT id, employee_name, manager_name, team_name, leave_type,  leave_dates, leave_duration FROM employee_leave_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]FormData, 0)

	for rows.Next() {
		var user FormData
		err := rows.Scan(&user.Id, &user.EmployeeName, &user.ManagerName,&user.TeamName, &user.LeaveType, &user.LeaveDate, &user.LeaveDuration)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func getApproveData(db *sql.DB) ([]ApprovedData, error) {
	rows, err := db.Query("SELECT id, manager_name, approved FROM notification")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]ApprovedData, 0)

	for rows.Next() {
		var user ApprovedData
		err := rows.Scan(&user.LeaveId, &user.ReportingManager, &user.Apporoved)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func kpi3(db *sql.DB) ([]Kpi3, error) {
	rows, err := db.Query("SELECT employee_name, leave_count FROM kpi3")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]Kpi3, 0)

	for rows.Next() {
		var user Kpi3
		err := rows.Scan(&user.EmployeeName, &user.LeaveCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


func kpi4(db *sql.DB) ([]Kpi4, error) {
	rows, err := db.Query("SELECT manager_name, leave_count FROM kpi4")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]Kpi4, 0)

	for rows.Next() {
		var user Kpi4
		err := rows.Scan(&user.ManagerName, &user.LeaveCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func kpi6One(db *sql.DB) ([]Kpi6, error) {
	rows, err := db.Query("SELECT leave_type, leave_count FROM kpi6 where team_name = 'AI'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]Kpi6, 0)

	for rows.Next() {
		var user Kpi6
		err := rows.Scan( &user.LeaveType, &user.LeaveCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func kpi6Two(db *sql.DB) ([]Kpi6, error) {
	rows, err := db.Query("SELECT leave_type, leave_count FROM kpi6 where team_name = 'IT'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]Kpi6, 0)

	for rows.Next() {
		var user Kpi6
		err := rows.Scan( &user.LeaveType, &user.LeaveCount)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}