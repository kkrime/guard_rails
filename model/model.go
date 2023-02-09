package model

import (
	"encoding/json"
	"guard_rails/config"
	"time"
)

type Repository struct {
	Id        int64     `db:"id" json:"-"`
	Name      string    `json:"name" binding:"required,alphanum" db:"name"`
	Url       string    `json:"url" binding:"required,url" db:"url"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
	// time.Time cannot be set to nil (for records that have not been deleted)
	// so setting DeletedAt to *time.Time
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

type ScanStatus string

const (
	Queued     ScanStatus = "QUEUED"
	InProgress ScanStatus = "IN PROGRESS"
	Success    ScanStatus = "SUCCESS"
	Failure    ScanStatus = "FAILURE"
)

type Scan struct {
	Id           int64      `db:"id"`
	RepositoryId int64      `db:"repository_id"`
	Status       ScanStatus `db:"status"`
	Findings     Findings   `db:"findings"`
	CreatedAt    time.Time  `db:"created_at" json:"-"`
	StartedAt    *time.Time `db:"started_at" json:"-"`
	EndeddAt     *time.Time `db:"ended_at" json:"-"`
}

type Findings []Finding

type Location struct {
	Path      string    `json:"path"`
	Positions Positions `json:"positions"`
}

type Positions struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line int64 `json:"line"`
}

type Finding struct {
	*config.ScanData
	Location Location         `json:"location"`
	MetaData *config.MetaData `json:"metadata"`
}

func (f *Findings) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	input := src.([]uint8)
	return json.Unmarshal(input, &f)
}
