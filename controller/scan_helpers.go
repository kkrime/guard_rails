package controller

import "guard_rails/model"

type scan struct {
	ScanId   int64            `json:"scan_id"`
	Status   model.ScanStatus `json:"status"`
	Findings model.Findings   `json:"findings,omitempty"`
}

func transsformScans(scans []model.Scan) []scan {

	result := make([]scan, len(scans))

	for i, scan := range scans {
		result[i].ScanId = scan.Id
		result[i].Status = scan.Status
		result[i].Findings = scan.Findings
	}

	return result

}
