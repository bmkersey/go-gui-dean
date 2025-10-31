package models

import "database/sql"

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
