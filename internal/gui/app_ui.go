package gui

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bmkersey/go-gui-dean/internal/db"
	"github.com/bmkersey/go-gui-dean/internal/models"
)

func RunApp(sqlDB *sql.DB) {

	myApp := app.New()
	win := myApp.NewWindow("District Selector")
	win.Resize(fyne.NewSize(520, 260))

	districts, err := db.FetchDistricts(sqlDB)
	if err != nil {
		log.Fatalf("❌ Failed to fetch districts: %v", err)
	}

	names := make([]string, len(districts))
	for i, d := range districts {
		names[i] = d.Name
	}
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
		schools, err := db.FetchSchoolsByDistrict(sqlDB, distID)
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
			students, err := db.FetchStudentsBySchool(sqlDB, schoolID)
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

func showStudentDetails(win fyne.Window, s models.Student) {
	info := fmt.Sprintf(
		"Name: %s %s\nGrade: %v\nStatus: %s\nParent: %v\nEnrolled On: %v",
		s.FirstName, s.LastName,
		s.GradeLevel.Int16, s.Status,
		s.ParentEmail.String, s.EnrolledOn.String,
	)
	dialog.ShowInformation("Student Details", info, win)
}
