package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type District struct {
	ID   int
	Name string
}

type School struct {
	ID   int
	Name string
}

type Student struct {
	ID          int
	FirstName   string
	LastName    string
	GradeLevel  sql.NullInt16
	ParentEmail sql.NullString
	Status      string
	EnrolledOn  sql.NullString
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults if any")
	}

	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Error opening DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Error connecting to DB: %v", err)
	}

	fmt.Println("✅ Connected to MySQL successfully!")

	rows, err := db.Query("SELECT id, name FROM districts ORDER BY name")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	defer rows.Close()

	districts, err := fetchDistricts(db)
	if err != nil {
		log.Fatalf("❌ Failed to fetch districts: %v", err)
	}

	// Prepare the list of names for the dropdown
	names := make([]string, len(districts))
	for i, d := range districts {
		names[i] = d.Name
	}

	// --- GUI ---
	myApp := app.New()
	win := myApp.NewWindow("District Selector")
	win.Resize(fyne.NewSize(520, 260))

	districtDropdown := widget.NewSelect(names, nil)
	districtDropdown.PlaceHolder = "Select a district..."

	schoolDropdown := widget.NewSelect([]string{}, nil)
	schoolDropdown.PlaceHolder = "Select a school..."
	schoolDropdown.Disable()

	studentBox := container.NewVBox(widget.NewLabel("👩‍🎓 Students"))
	studentBox.Hide()

	districtDropdown.OnChanged = func(name string) {
		var distID int
		for _, d := range districts {
			if d.Name == name {
				distID = d.ID
				break
			}
		}
		if distID == 0 {
			dialog.ShowError(fmt.Errorf("could not find district ID for %q", name), win)
			return
		}
		schools, err := fetchSchoolsByDistrict(db, distID)
		if err != nil {
			dialog.ShowError(fmt.Errorf("load schools: %w", err), win)
		}

		schoolNames := make([]string, len(schools))
		for i, s := range schools {
			schoolNames[i] = s.Name
		}
		schoolDropdown.Options = schoolNames
		schoolDropdown.SetSelected("") // clear previous choice
		schoolDropdown.Enable()
		schoolDropdown.Refresh()

		schoolDropdown.OnChanged = func(schoolName string) {
			var schoolID int
			for _, s := range schools {
				if s.Name == schoolName {
					schoolID = s.ID
					break
				}
			}
			students, err := fetchStudentsBySchool(db, schoolID)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			studentBox.Objects = []fyne.CanvasObject{widget.NewLabel("👩‍🎓 Students")}
			if len(students) == 0 {
				studentBox.Add(widget.NewLabel("No students found for this school."))
				studentBox.Show()
				studentBox.Refresh()
				return
			}

			for _, s := range students {
				stu := s // capture for closure
				btn := widget.NewButton(fmt.Sprintf("%s %s", s.FirstName, s.LastName), func() {
					showStudentDetails(win, stu)
				})
				studentBox.Add(btn)
			}
			studentBox.Show()
			studentBox.Refresh()
		}
	}

	content := container.NewVBox(
		widget.NewLabel("📚 District & School Browser"),
		widget.NewLabel("District"),
		districtDropdown,
		widget.NewLabel("School"),
		schoolDropdown,
		studentBox,
	)

	win.SetContent(content)
	win.ShowAndRun()
}

func fetchDistricts(db *sql.DB) ([]District, error) {
	rows, err := db.Query("SELECT id, name FROM districts ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var districts []District
	for rows.Next() {
		var d District
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		districts = append(districts, d)
	}
	return districts, rows.Err()
}

func fetchSchoolsByDistrict(db *sql.DB, districtID int) ([]School, error) {
	rows, err := db.Query(`
		SELECT id, name
		FROM schools
		WHERE district_id = ?
		ORDER BY name
	`, districtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []School
	for rows.Next() {
		var s School
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func fetchStudentsBySchool(db *sql.DB, schoolID int) ([]Student, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, grade_level, parent_email, status, enrolled_on
		FROM students
		WHERE school_id = ?
		ORDER BY last_name, first_name`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.GradeLevel, &s.ParentEmail, &s.Status, &s.EnrolledOn); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func showStudentDetails(win fyne.Window, s Student) {
	info := fmt.Sprintf(
		"Name: %s %s\nGrade: %v\nStatus: %s\nParent: %v\nEnrolled On: %v",
		s.FirstName, s.LastName,
		s.GradeLevel.Int16, s.Status,
		s.ParentEmail.String, s.EnrolledOn.String,
	)
	dialog.ShowInformation("Student Details", info, win)
}
