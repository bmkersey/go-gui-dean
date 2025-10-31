package db

import (
	"database/sql"

	"github.com/bmkersey/go-gui-dean/internal/models"
)

func FetchDistricts(db *sql.DB) ([]models.District, error) {
	rows, err := db.Query(`SELECT id, name FROM districts ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.District
	for rows.Next() {
		var d models.District
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
