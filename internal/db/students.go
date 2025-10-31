package db

import (
	"database/sql"

	"github.com/bmkersey/go-gui-dean/internal/models"
)

func FetchStudentsBySchool(db *sql.DB, schoolID int) ([]models.Student, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, grade_level, parent_email, status, enrolled_on
		FROM students
		WHERE school_id = ?
		ORDER BY last_name, first_name`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.GradeLevel, &s.ParentEmail, &s.Status, &s.EnrolledOn); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
