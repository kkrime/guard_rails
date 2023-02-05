package model

import "time"

type Repository struct {
	Id        string    `db:"id" json:"-"`
	Name      string    `json:"name" binding:"required,alphanum" db:"name"`
	Url       string    `json:"url" binding:"required,url" db:"url"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	// time.Time cannot be set to nil (for records that have not been deleted)
	// so setting DeletedAt to *time.Time
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
