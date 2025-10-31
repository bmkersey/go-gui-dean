package db

import (
	"database/sql"

	"github.com/bmkersey/go-gui-dean/internal/models"
)

func FetchSchoolsByDistrict(db *sql.DB, districtID int) ([]models.School, error) {
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

	var out []models.School
	for rows.Next() {
		var s models.School
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
