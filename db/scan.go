package db

import (
	"context"
	"encoding/json"
	"fmt"
	"guard_rails/model"

	"github.com/jmoiron/sqlx"
)

func (sd *db) GetScanWithStatus(ctx context.Context, repositoryId int64, status []model.ScanStatus) (*model.Scan, error) {
	var scans []model.Scan

	statement := `
        SELECT * FROM  
            scans
        WHERE
            repository_id = ? AND
            status IN (?)
        ;`

	statement, args, err := sqlx.In(statement, repositoryId, status)
	if err != nil {
		return nil, err
	}

	statement = sqlx.Rebind(sqlx.DOLLAR, statement)

	err = sd.db.SelectContext(ctx, &scans, statement, args...)
	if err != nil {
		return nil, err
	}

	if scans == nil {
		return nil, nil
	}

	return &scans[0], err
}

func (sd *db) CreateNewScan(ctx context.Context, repositoryId int64) (*model.Scan, error) {

	scan := &model.Scan{}

	statement := `
        INSERT INTO 
            scans
        (
            repository_id
        )
        VALUES 
        (
            $1
        )
        RETURNING
            *
        ;`

	fmt.Println("HERE")

	err := sd.db.GetContext(ctx, scan, statement, repositoryId)
	if err != nil {
		return nil, err
	}

	return scan, err
}

func (sd *db) GetScans(ctx context.Context, repositoryName string) ([]model.Scan, error) {
	var scans []model.Scan

	statement := `
        SELECT * FROM  
            scans
        WHERE
            repository_id =
            (
                SELECT id FROM
                    repositories
                WHERE
                    name = $1 AND
                    deleted_at IS NULL
            )
        ;`

	err := sd.db.SelectContext(ctx, &scans, statement, repositoryName)
	if err != nil {
		return nil, err
	}

	return scans, err

}

func (sd *db) UpdateScanStatus(scanId int64, status model.ScanStatus) error {

	statement := `
        UPDATE
            scans
        SET
            status = $1
        WHERE
            id = $2
        ;`

	_, err := sd.db.Exec(statement, status, scanId)

	return err
}

func (sd *db) StartScan(scanId int64) error {

	statement := `
        UPDATE
            scans
        SET
            status = $1,
            started_at = now()
        WHERE
            id = $2
        ;`

	_, err := sd.db.Exec(statement, model.InProgress, scanId)

	return err
}

func (sd *db) StopScan(scanId int64, findings model.Findings, status model.ScanStatus) error {
	var (
		findingsJson []byte
		err          error
	)

	if findings != nil {
		findingsJson, err = json.Marshal(findings)
		if err != nil {
			return err
		}
	}

	statement := `
        UPDATE
            scans
        SET
            status = $1,
            findings = $2,
            ended_at = now()
        WHERE
            id = $3
        ;`

	_, err = sd.db.Exec(statement, status, findingsJson, scanId)

	return err
}
