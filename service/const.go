package service

// errors
const (
	// repository
	Repository_Already_Added = "repository already added"
	Repository_Not_Found     = "repository not found"

	// scan
	Scan_Already_Exists = "scan already exists queued or running"
	No_Scans_Found      = "no scans found"

	Internal_Error             = "internal error"
	Unable_To_Reach_Repository = "unable to reach repository"
)
